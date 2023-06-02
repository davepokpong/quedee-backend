package utils

import "golang.org/x/crypto/bcrypt"

func GenerateHashPassword(password string) (string, error) {
	p := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
	return string(hashedPassword), err
}
