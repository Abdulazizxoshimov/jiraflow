package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

// Hasher is the interface for password hashing and verification.
type Hasher interface {
	Hash(password string) (string, error)
	Check(plain, hashed string) bool
}

type bcryptHasher struct {
	cost int
}

// New returns a Hasher backed by bcrypt with the given cost.
// Pass bcrypt.DefaultCost (10) for most cases.
func New(cost int) Hasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &bcryptHasher{cost: cost}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Check returns true if plain matches the stored hash.
func (h *bcryptHasher) Check(plain, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
