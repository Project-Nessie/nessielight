package nessielight

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/Project-Nessie/nessielight/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var logger *log.Logger

type GormDB = gorm.DB

var DB *GormDB

func InitDBwithFile(path string) error {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}

var AuthServiceInstance TelegramAuthService
var UserManagerInstance UserManager
var V2rayServiceInstance V2rayService

func InitV2rayService(inboundTag string, vmessPort int, vmessAddress, wsPath, v2rayApi string) error {
	client := v2rayClient{
		inboundTag: inboundTag,
		port:       vmessPort,
		domain:     vmessAddress,
		path:       wsPath,
	}
	V2rayServiceInstance = &client

	if err := V2rayServiceInstance.Start(v2rayApi); err != nil {
		return err
	}
	return nil
}

func init() {
	AuthServiceInstance = &simpleTelegramAuthService{
		userManager: &UserManagerInstance,
		tokenDB:     make(map[string]bool),
	}
	UserManagerInstance = &simpleUserManager{
		db: make(map[string]User),
	}
}

// Interface for User. Typically implemented by UserManager.NewUser
type User interface {
	ID() string
	Name() string
	Proxy() []Proxy
	SetProxy(proxy []Proxy) error
	SetName(name string) error
	// total traffic stored
	Traffic() TrafficValue
	SetTraffic(val TrafficValue) error
}

// implemented by simpleUserManager
type UserManager interface {
	AddUser(user User) error
	SetUser(user User) error
	DeleteUser(user User) error
	// find user by id, nil for not found
	FindUser(id string) (User, error)
	// find user by proxy id, nil for not found
	FindUserByProxy(proxyid string) (User, error)
	// generate new user by id
	NewUser(id string) User
	All() ([]User, error)
}

// implemented by v2rayClient
type V2rayService interface {
	SetUser(email string, uuid string) error
	AddUser(email string) (uuid string, err error)
	// remove a user identified by email
	RemoveUser(email string) error
	QueryTraffic(pattern string, reset bool) (stat []V2rayTrafficStat, err error)
	// query traffic for user under control only
	QueryUserTraffic(reset bool) (stat []V2rayTrafficStat, err error)
	Start(listen string) error
	VmessText(vmessid string) string
	VmessLink(vmessid string) string
	NewProxy() Proxy
}

// implemented by simpleTelegramAuthService
type TelegramAuthService interface {
	// 生成一个注册用的 token
	GenToken() (token string)
	// 使用 token 注册用户，注册失败（token不匹配）返回错误
	Register(token string, id string) (User, error)
}

// need implementation
type SystemCtlService interface {
	StartV2rayServer() error
	StopV2rayServer() error
	RestartV2rayServer() error
}

// describe a proxy config. Proxy can be store in sqldb
type Proxy interface {
	// sql.Scanner
	// sqldriver.Valuer
	// identify this proxy
	ID() string
	// apply this proxy
	Activate() error
	// remove this proxy
	Deactivate() error
	// introduce this proxy in telegram message
	Message() string
}

func GetUserProxyMessage(user User) string {
	msg := "Proxy of " + user.ID() + "\n"
	for _, proxy := range user.Proxy() {
		msg += proxy.Message() + "\n"
	}
	return msg
}

func ApplyUserProxy(user User) error {
	for _, proxy := range user.Proxy() {
		if err := proxy.Activate(); err != nil {
			return fmt.Errorf("ApplyUserProxy(id=%s): %s", user.ID(), err.Error())
		}
	}
	return nil
}

type TrafficValue struct {
	Uplink, Downlink utils.ByteValue
}

var reTraffic = regexp.MustCompile(`(.*?)>>>(.*?)>>>traffic>>>((downlink|uplink))`)

func trafficNameMatch(trafficname string) (category, name, linktype string) {
	if matches := reTraffic.FindSubmatch([]byte(trafficname)); len(matches) > 3 {
		return string(matches[1]), string(matches[2]), string(matches[3])
	}
	return "", "", ""
}

const UUIDLen = 36

// Get current traffic stats of all inbounds.
// Note that user stats is not correct due to V2rayUpdateUserTraffic
func GetV2rayTraffic() ([]NamedTraffic, error) {
	stats, err := V2rayServiceInstance.QueryTraffic("inbound>>>", false)
	if err != nil {
		return nil, err
	}
	var traffics map[string]*NamedTraffic = make(map[string]*NamedTraffic)

	for _, v := range stats {
		_, name, linktype := trafficNameMatch(v.Name)
		id := name
		if traffics[id] == nil {
			traffics[id] = &NamedTraffic{
				Name: name,
			}
		}

		if linktype == "downlink" {
			traffics[id].Downlink = utils.ByteValue(v.Value)
		} else if linktype == "uplink" {
			traffics[id].Uplink = utils.ByteValue(v.Value)
		}
	}

	return utils.Flatten(traffics, func(val *NamedTraffic) NamedTraffic {
		return *val
	}), nil
	// var usertraffic []NamedTraffic = make([]NamedTraffic, 0, len(traffics))
	// for _, v := range traffics {
	// 	logger.Print(v)
	// 	usertraffic = append(usertraffic, *v)
	// }
	// return usertraffic, nil
}

type NamedTraffic struct {
	TrafficValue
	Name string
}

func V2rayUpdateUserTraffic() error {
	stats, err := V2rayServiceInstance.QueryUserTraffic(true) // !!! false just for debug
	if err != nil {
		return err
	}
	for _, v := range stats {
		_, name, linktype := trafficNameMatch(v.Name)
		id := name[len(name)-UUIDLen:]
		logger.Print("V2rayUpdateUserTraffic ", name, " ", linktype, " ", utils.ByteValue(v.Value))
		logger.Print("V2rayUpdateUserTraffic id=", id)
		if user, err := UserManagerInstance.FindUserByProxy(id); err == nil && user != nil {
			data := user.Traffic()
			if linktype == "downlink" {
				data.Downlink += utils.ByteValue(v.Value)
			} else if linktype == "uplink" {
				data.Uplink += utils.ByteValue(v.Value)
			}
			if err := user.SetTraffic(data); err != nil {
				return err
			}
			if err := UserManagerInstance.SetUser(user); err != nil {
				return err
			}
			logger.Print("V2rayUpdateUserTraffic user:", user.Traffic())
		}
	}
	return nil
}

func init() {
	logger = log.New(os.Stderr, "[nessielight] ", log.LstdFlags|log.Lmsgprefix)
}
