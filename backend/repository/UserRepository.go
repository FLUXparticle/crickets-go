package repository

import "crickets-go/data"

type UserRepository struct {
	// User-Tabelle (in einer echten Anwendung sollten diese aus einer Datenbank oder z.B. LDAP kommen)
	users map[string]*data.User
}

func NewUserRepository() *UserRepository {
	repository := &UserRepository{
		users: make(map[string]*data.User),
	}

	repository.Save(&data.User{Username: "admin", Password: "Secret123"})
	repository.Save(&data.User{Username: "helpdesk", Password: "Secret123"})
	repository.Save(&data.User{Username: "employee", Password: "Secret123"})
	repository.Save(&data.User{Username: "manager", Password: "Secret123"})

	return repository
}

func (r *UserRepository) FindByUsername(username string) *data.User {
	if user, found := r.users[username]; found {
		return user
	} else {
		return nil
	}
}

func (r *UserRepository) Save(user *data.User) {
	r.users[user.Username] = user
	user.ID = int32(len(r.users))
}
