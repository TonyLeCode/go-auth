package auth

import "golang.org/x/crypto/argon2"

func HashPassword(password string) []byte {
	return argon2.IDKey([]byte(password), nil, 1, 64*1024, 4, 32)
}
