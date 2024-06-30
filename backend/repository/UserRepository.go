package repository

type UserRepository struct {
	// User-Tabelle (in einer echten Anwendung sollten diese aus einer Datenbank oder z.B. LDAP kommen)
	users map[string]*User
}

func NewUserRepository() *UserRepository {
	repository := &UserRepository{
		users: make(map[string]*User),
	}

	repository.Save(&User{Username: "admin", Password: "Secret123"})
	repository.Save(&User{Username: "helpdesk", Password: "Secret123"})
	repository.Save(&User{Username: "employee", Password: "Secret123"})
	repository.Save(&User{Username: "manager", Password: "Secret123"})

	return repository
}

func (r *UserRepository) FindByUsername(username string) *User {
	if user, found := r.users[username]; found {
		return user
	} else {
		return nil
	}
}

func (r *UserRepository) Save(user *User) {
	r.users[user.Username] = user
	user.ID = len(r.users)
}
