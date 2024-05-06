package internal

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username string
	Password []byte
}

func (u User) Compare(password []byte) bool {
	return bcrypt.CompareHashAndPassword(u.Password, password) == nil
}
