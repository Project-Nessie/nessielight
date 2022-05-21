package nessielight

import (
	"database/sql"
	sqldriver "database/sql/driver"
	"fmt"
	"regexp"

	"github.com/Project-Nessie/nessielight/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
	QueryUserTraffic(pattern string, reset bool) (stat []UserTrafficStat, err error)
	Start(listen string) error
	VmessText(vmessid string) string
	VmessLink(vmessid string) string
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
	sql.Scanner
	sqldriver.Valuer
	// identify this proxy
	ID() string
	// apply this proxy
	Activate() error
	// remove this proxy
	Deactivate() error
	// introduce this proxy
	Message() string
}

type v2rayProxy struct {
	id string
}

func (r *v2rayProxy) ID() string {
	return r.id
}
func (r *v2rayProxy) Activate() error {
	V2rayServiceInstance.RemoveUser(r.id)
	return V2rayServiceInstance.SetUser(r.id, r.id)
}
func (r *v2rayProxy) Deactivate() error {
	return V2rayServiceInstance.RemoveUser(r.id)
}
func (r *v2rayProxy) Message() string {
	return "<code>" + V2rayServiceInstance.VmessLink(r.id) + "</code>"
}
func (r *v2rayProxy) Value() (sqldriver.Value, error) {
	return r.id, nil
}
func (r *v2rayProxy) Scan(src interface{}) error {
	if id, ok := src.(string); ok {
		r.id = id
		return nil
	}
	return fmt.Errorf("invalid src type when scanning v2rayProxy")
}

func CreateV2rayProxy() Proxy {
	proxy := v2rayProxy{
		id: NewUUID(),
	}
	return &proxy
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

type UserType int

const (
	Unregistered UserType = 0
	Registered   UserType = 1
	Inbound      UserType = 2
)

type UserTraffic struct {
	ID               string
	Name             string
	Type             UserType
	Uplink, Downlink int64
}

var reTraffic = regexp.MustCompile(`(.*?)>>>(.*?)>>>traffic>>>((downlink|uplink))`)

func GetV2rayTraffic() (UserTraffics, error) {
	stats, err := V2rayServiceInstance.QueryUserTraffic("", false)
	if err != nil {
		return nil, err
	}
	var traffics map[string]*UserTraffic = make(map[string]*UserTraffic)

	for _, v := range stats {
		if matches := reTraffic.FindSubmatch([]byte(v.Name)); len(matches) > 3 {
			category := string(matches[1])
			name, id := string(matches[2]), ""
			if category == "user" {
				typ := Unregistered
				id = name
				if user, err := UserManagerInstance.FindUserByProxy(id); err == nil && user != nil {
					// registered
					name = user.Name()
					typ = Registered
				}
				if traffics[id] == nil {
					traffics[id] = &UserTraffic{
						ID:   id,
						Name: name,
						Type: typ,
					}
				}
			} else if category == "inbound" {
				id = "inbound-" + name
				if traffics[id] == nil {
					traffics[id] = &UserTraffic{
						ID:   id,
						Name: name,
						Type: Inbound,
					}
				}
			} else {
				continue
			}

			if string(matches[3]) == "downlink" {
				traffics[id].Downlink = v.Value
			} else if string(matches[3]) == "uplink" {
				traffics[id].Uplink = v.Value
			}
		}
	}

	var usertraffic []UserTraffic = make([]UserTraffic, 0, len(traffics))
	for _, v := range traffics {
		logger.Print(v)
		usertraffic = append(usertraffic, *v)
	}
	return usertraffic, nil
}

type UserTraffics []UserTraffic

func (r UserTraffics) Users() UserTraffics {
	return utils.Filter(r, func(t UserTraffic) bool {
		return t.Type != Inbound
	})
}

func (r UserTraffics) Inbounds() UserTraffics {
	return utils.Filter(r, func(t UserTraffic) bool {
		return t.Type == Inbound
	})
}

// default sort key
func (a UserTraffics) Len() int           { return len(a) }
func (a UserTraffics) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UserTraffics) Less(i, j int) bool { return a[i].Downlink > a[j].Downlink }
