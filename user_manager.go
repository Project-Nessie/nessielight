package nessielight

type simpleUserManager struct {
	db map[string]User
}

func (r *simpleUserManager) AddUser(user User) error {
	r.db[user.ID()] = user
	return nil
}

func (r *simpleUserManager) SetUser(user User) error {
	r.db[user.ID()] = user
	return nil
}

func (r *simpleUserManager) DeleteUser(user User) error {
	delete(r.db, user.ID())
	return nil
}

func (r *simpleUserManager) FindUser(id string) (User, error) {
	if r.db[id] != nil {
		return r.db[id], nil
	}
	return nil, nil
}

func (r *simpleUserManager) All() ([]User, error) {
	var users []User = make([]User, 0, len(r.db))
	for _, v := range r.db {
		users = append(users, v)
	}
	return users, nil
}

var _ UserManager = (*simpleUserManager)(nil)

type simpleUser struct {
	id    string
	name  string
	proxy []Proxy
}

func (r *simpleUser) ID() string {
	return r.id
}

func (r *simpleUser) Name() string {
	return r.name
}

func (r *simpleUser) Proxy() []Proxy {
	return r.proxy
}

func (r *simpleUser) SetProxy(proxy []Proxy) error {
	r.proxy = proxy
	return nil
}

func (r *simpleUser) SetName(name string) error {
	r.name = name
	return nil
}

var _ User = (*simpleUser)(nil)

func (r *simpleUserManager) NewUser(id string) User {
	user := simpleUser{id: id}
	return &user
}
