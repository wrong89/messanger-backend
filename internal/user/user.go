package user

import "time"

type User struct {
	ID           uint
	Name         string
	Login        string
	PasswordHash []byte
	CreatedAt    time.Time
}

func NewUser(name, login string, passHash []byte) User {
	return User{
		Name:         name,
		Login:        login,
		PasswordHash: passHash,
	}
}
