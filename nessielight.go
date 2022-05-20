package nessielight

import (
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

// interface for User
// need implementation
type User interface {
	ID() string
	Email() string
}

// need implementation
type UserManager interface {
	AddUser(user User) error
	SetUser(user User) error
	DeleteUser(user User) error
	// find user by id, nil for not found
	FindUser(id string) (User, error)
	// generate new user by id
	NewUser(id string) User
}

// implemented by V2rayClient
type V2rayService interface {
	SetUser(email string, uuid string) error
	AddUser(email string) (uuid string, err error)
	RemoveUser(email string) error
	QueryUserTraffic(pattern string, reset bool) (stat []UserTrafficStat, err error)
	Start(listen string) error
}

// need implementation
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
