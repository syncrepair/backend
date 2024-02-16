package password

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	Hash(password string) string
	Compare(password, hashedPassword string) bool
}

type hasher struct {
	salt string
	cost int
}

func NewHasher(salt string, hashingCost int) Hasher {
	return &hasher{
		salt: salt,
		cost: hashingCost,
	}
}

func (h *hasher) Hash(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		panic("error hashing password: " + err.Error())
	}

	return string(hashedPassword)
}

func (h *hasher) Compare(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
