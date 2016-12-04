package hash

import (
	"bytes"
	"crypto/rand"
	"golang.org/x/crypto/scrypt"
)

const (
	salt_length     = 32
	hash_length     = 64
	work_factor     = 16384
	mem_blocksize   = 8
	parallel_factor = 1
)

func Hash(password string) (h, s []byte, e error) {
	salt := make([]byte, salt_length)
	_, err := rand.Read(salt)
	if err != nil {
		return make([]byte, hash_length), salt, err
	}
	hash, err := scrypt.Key([]byte(password), salt, work_factor, mem_blocksize, parallel_factor, hash_length)
	return hash, salt, err
}

func Verify(password string, salt, passhash []byte) bool {
	dk, err := scrypt.Key([]byte(password), salt, work_factor, mem_blocksize, parallel_factor, hash_length)
	if err != nil {
		return false
	}
	return bytes.Equal(dk, passhash)
}
