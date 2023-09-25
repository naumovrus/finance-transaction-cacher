package passwordhash

import (
	"crypto/sha1"
	"fmt"
	"time"
)

const (
	salt     = "123zxdqwe647qwerproetidfasdkjhgkd"
	tokenTTL = 12 * time.Hour
)

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
