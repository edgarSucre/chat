package usecase

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
}

func (Hasher) IsPasswordValid(pass, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
	return err == nil
}

func (Hasher) SecurePassword(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hashed), err
}
