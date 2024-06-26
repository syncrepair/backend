package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

type PasswordHasher interface {
	Hash(string) string
	Check(string, string) bool
}

type passwordHasher struct {
	salt string
}

func NewPasswordHasher(salt string) PasswordHasher {
	return &passwordHasher{
		salt: salt,
	}
}

func (h *passwordHasher) Hash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password + h.salt))

	return hex.EncodeToString(hash.Sum(nil))
}

func (h *passwordHasher) Check(hashedPassword string, password string) bool {
	return h.Hash(password+h.salt) == hashedPassword
}
