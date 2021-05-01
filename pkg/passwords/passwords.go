package passwords

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrTooShort      = errors.New("password too short")
	ErrCannotGetHash = errors.New("cannot get hash")
)

type Passworder struct {
	KeySecret []byte
	MinLen    int
}

func (pw *Passworder) Hash(password string) ([]byte, error) {
	if len(password) <= pw.MinLen {
		return nil, ErrTooShort
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(append([]byte(password), pw.KeySecret...), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrCannotGetHash
	}
	return hashedPassword, nil
}

func (pw *Passworder) IsCorrect(knownPasswordHash []byte, userPasswordInput string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(knownPasswordHash, append([]byte(userPasswordInput), pw.KeySecret...))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, fmt.Errorf("compare password failture: %v", err)
	}
	return true, nil
}
