package service

var UserManagerInstance UserManager

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

type simpleUser struct {
	id    string
	email string
}

func (r *simpleUser) ID() string {
	return r.id
}

func (r *simpleUser) Email() string {
	return r.email
}

// !!!
func (r *simpleUser) IsAdmin() bool {
	return true
}

func (r *simpleUserManager) NewUser(id string) User {
	user := simpleUser{id: id}
	return &user
}

func init() {
	UserManagerInstance = &simpleUserManager{
		db: make(map[string]User),
	}
}
