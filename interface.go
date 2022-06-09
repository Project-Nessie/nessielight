package nessielight

// Interface for User. Typically implemented by UserManager.NewUser
type User interface {
	TelegramID() int
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
	FindUserByTelegramID(tid int) (User, error)
	// find user by proxy id, nil for not found
	FindUserByProxy(proxyid uint) (User, error)
	// generate new user by id
	NewUser(tid int) User
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
	AddVmessWsInbound(tag string, port uint16, wspath string) error
	RemoveInbound(tag string) error
}

// implemented by simpleTelegramAuthService
type TelegramAuthService interface {
	// 生成一个注册用的 token
	GenToken() (token string)
	// 使用 token 注册用户，注册失败（token不匹配）返回错误
	Register(token string, tid int) (User, error)
}

// need implementation
type SystemCtlService interface {
	StartV2rayServer() error
	StopV2rayServer() error
	RestartV2rayServer() error
}

// describe a proxy config. Proxy can be store in sqldb
type Proxy interface {
	// identify this proxy
	ProxyID() uint
	// apply this proxy
	Activate() error
	// remove this proxy
	Deactivate() error
	// introduce this proxy in telegram message
	Message() string
}
