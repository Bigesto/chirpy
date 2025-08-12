package auth

import "golang.org/x/crypto/bcrypt"

func CheckPasswordHash(password, hash string) error {
	bytePass := []byte(password)
	byteHash := []byte(hash)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePass)
	if err != nil {
		return err
	}

	return nil
}
