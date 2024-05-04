package internal

import "golang.org/x/crypto/bcrypt"

type user struct {
	username string
	password []byte
}

func (u user) Username() string {
	return u.username
}

func (u user) Password() []byte {
	return u.password
}

func (u user) Compare(password []byte) bool {
	return bcrypt.CompareHashAndPassword(u.password, password) == nil
}
