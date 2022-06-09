package nessielight

import (
	"fmt"

	"github.com/Project-Nessie/nessielight/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type simpleUserManager struct {
}

func (r *simpleUserManager) AddUser(user User) error {
	if userdata, ok := user.(*simpleUser); ok {
		DataBase.Create(userdata)
		logger.Print("AddUser id=", userdata.ID, " tid=", userdata.Registerid)
	} else {
		return fmt.Errorf("invalid user type")
	}
	return nil
}

func (r *simpleUserManager) SetUser(user User) error {
	if userdata, ok := user.(*simpleUser); ok {
		logger.Print("SaveUser id=", userdata.ID, " tid=", userdata.Registerid)
		DataBase.Save(userdata)
	} else {
		return fmt.Errorf("invalid user type")
	}
	return nil
}

func (r *simpleUserManager) DeleteUser(user User) error {
	if userdata, ok := user.(*simpleUser); ok {
		logger.Print("DeleteUser id=", userdata.ID, " tid=", userdata.Registerid)
		if userdata.ID == 0 {
			return fmt.Errorf("delete user without ID")
		}
		DataBase.Delete(userdata)
	} else {
		return fmt.Errorf("invalid user type")
	}
	return nil
}

func (r *simpleUserManager) FindUserByTelegramID(tid int) (User, error) {
	var user simpleUser
	DataBase.Where(&simpleUser{Registerid: tid}).First(&user)
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}
func (r *simpleUserManager) FindUserByProxy(proxyid uint) (User, error) {
	var users []simpleUser
	DataBase.Find(&users)
	for _, v := range users {
		for _, p := range v.V2rayProxyID {
			if p == int32(proxyid) {
				return &v, nil
			}
		}
	}
	return nil, nil
}

func (r *simpleUserManager) All() ([]User, error) {
	var users []simpleUser
	DataBase.Find(&users)
	return utils.Map(users, func(user simpleUser) User {
		return &user
	}), nil
}

var _ UserManager = (*simpleUserManager)(nil)

type simpleUser struct {
	gorm.Model
	Registerid int
	Nam        string
	// proxy id
	V2rayProxyID pq.Int32Array `gorm:"type:text[]"`
	Traff        TrafficValue  `gorm:"embedded"`
}

func (r *simpleUser) TelegramID() int {
	return r.Registerid
}

func (r *simpleUser) Name() string {
	return r.Nam
}

func (r *simpleUser) Proxy() []Proxy {
	return utils.Map(r.V2rayProxyID, func(id int32) Proxy {
		var v2rayproxy v2rayProxy
		DataBase.First(&v2rayproxy, id)
		return &v2rayproxy
	})
}

func (r *simpleUser) SetProxy(proxy []Proxy) error {
	r.V2rayProxyID = make([]int32, 0, len(r.V2rayProxyID))
	for _, v := range proxy {
		if v2rayproxy, ok := v.(*v2rayProxy); ok {
			r.V2rayProxyID = append(r.V2rayProxyID, int32(v2rayproxy.ID))
		} else {
			panic(fmt.Errorf("unknown proxy type"))
		}
	}
	return nil
}

func (r *simpleUser) SetName(name string) error {
	r.Nam = name
	return nil
}

func (r *simpleUser) Traffic() TrafficValue {
	return r.Traff
}
func (r *simpleUser) SetTraffic(val TrafficValue) error {
	r.Traff = val
	return nil
}

var _ User = (*simpleUser)(nil)

func (r *simpleUserManager) NewUser(tid int) User {
	user := simpleUser{Registerid: tid}
	return &user
}
