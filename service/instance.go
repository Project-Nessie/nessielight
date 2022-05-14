package service

var AuthServiceInstance TelegramAuthService
var UserManagerInstance UserManager
var V2rayServiceInstance V2rayService

func init() {

	AuthServiceInstance = &simpleTelegramAuthService{
		userManager: &UserManagerInstance,
		tokenDB:     make(map[string]bool),
	}
	UserManagerInstance = &simpleUserManager{
		db: make(map[string]User),
	}

}
