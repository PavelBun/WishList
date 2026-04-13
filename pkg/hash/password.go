// Package hash provides password hashing utilities.
package hash

import "golang.org/x/crypto/bcrypt"

// Password returns the bcrypt hash of the password.
func Password(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Check compares a password with its hash.
func Check(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
