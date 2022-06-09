package nessielight

import (
	"fmt"

	"github.com/Project-Nessie/nessielight/utils"
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

func (r *simpleUserManager) NewUser(tid int) User {
	user := simpleUser{Registerid: tid}
	return &user
}
