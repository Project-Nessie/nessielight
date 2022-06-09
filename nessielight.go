package nessielight

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/Project-Nessie/nessielight/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormDB = gorm.DB

var logger *log.Logger
var AuthServiceInstance TelegramAuthService
var UserManagerInstance UserManager
var V2rayServiceInstance V2rayService
var DataBase *GormDB

func InitDBwithFile(path string) error {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}
	DataBase = db
	if err := DataBase.AutoMigrate(&v2rayProxy{}); err != nil {
		return err
	}
	if err := DataBase.AutoMigrate(&simpleUser{}); err != nil {
		return err
	}
	var proxies []v2rayProxy
	DataBase.Find(&proxies)
	for _, v := range proxies {
		logger.Print(&v)
	}
	var users []simpleUser
	DataBase.Find(&users)
	for _, v := range users {
		logger.Print("user id=", v.ID, " tid=", v.Registerid, " proxy:", v.V2rayProxyID)
	}
	return nil
}

func InitV2rayService(inboundTag string, vmessPort, vmessClientPort int, vmessAddress, wsPath, v2rayApi string) error {
	client := v2rayClient{
		inboundTag: inboundTag,
		port:       vmessPort,
		clientport: vmessClientPort,
		domain:     vmessAddress,
		path:       wsPath,
	}
	V2rayServiceInstance = &client

	if err := V2rayServiceInstance.Start(v2rayApi); err != nil {
		return err
	}
	V2rayServiceInstance.RemoveInbound(inboundTag)
	if err := V2rayServiceInstance.AddVmessWsInbound(inboundTag, uint16(vmessPort), wsPath); err != nil {
		return err
	}

	return nil
}

func init() {
	AuthServiceInstance = &simpleTelegramAuthService{
		userManager: &UserManagerInstance,
		tokenDB:     make(map[string]bool),
	}
	UserManagerInstance = &simpleUserManager{}
}

func GetUserProxyMessage(user User) string {
	msg := fmt.Sprint("Proxy of ", user.TelegramID(), "\n")
	for _, proxy := range user.Proxy() {
		msg += proxy.Message() + "\n"
	}
	return msg
}

func ApplyUserProxy(user User) error {
	for _, proxy := range user.Proxy() {
		if err := proxy.Activate(); err != nil {
			return fmt.Errorf("ApplyUserProxy(id=%d): %s", user.TelegramID(), err.Error())
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
		id := name[len(V2rayServiceInstance.(*v2rayClient).inboundTag):]

		uid, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			logger.Print("V2rayUpdateUserTraffic invalid email: ", name)
			continue
		}

		logger.Print("V2rayUpdateUserTraffic ", name, " ", linktype, " ", utils.ByteValue(v.Value))
		logger.Print("V2rayUpdateUserTraffic id=", id)
		if user, err := UserManagerInstance.FindUserByProxy(uint(uid)); err == nil && user != nil {
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

func Restore() error {
	users, err := UserManagerInstance.All()
	if err != nil {
		return err
	}
	for _, v := range users {
		for _, p := range v.Proxy() {
			p.Activate()
		}
	}
	return nil
}
