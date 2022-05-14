// package service include v2ray API control service
package service

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
