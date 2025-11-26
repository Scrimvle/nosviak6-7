package database

import (
	"bytes"
	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

func NewSalt(i int) *[]byte {
	alloc := make([]byte, i)
	if _, err := rand.Read(alloc); err != nil {
		return NewSalt(i)
	}

	return &alloc
}

func NewHash(password []byte, key *[]byte) []byte {
	if key == nil {
		key = NewSalt(16)
	}

	return argon2.IDKey(password, *key, 3, (64 * 1024), 2, 64)
}

// IsPassword checks if the password provided equals the users password
func (u *User) IsPassword(password []byte) bool {
	return bytes.Equal(NewHash(password, &u.Key), u.Password)
}
