// Note that REST API SHOULD NOT depend on ANY OTHER service except database connection
package service

// interface for User
// need implementation
type User interface {
	ID() string
	Email() string
	IsAdmin() bool
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
