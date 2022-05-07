// Note that REST API SHOULD NOT depend on ANY OTHER service except database connection
package service

// interface for User
// need implementation
type User interface {
	GetID() string
	GetEmail() string
	isAdmin() bool
}

// need implementation
type UserManager interface {
	AddUser(user User) error
	SetUser(user User) error
	DeleteUser(user User) error
	// find user by id, nil for not found
	FindUser(id string) (User, error)
}
