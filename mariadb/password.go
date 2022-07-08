package mariadb

import (
	"crypto/sha256"
	"fmt"
)

func Password(password string) string {
	s := fmt.Sprintf("阿嬤說少林功夫好ㄟ~@.@%s#.#真是好ㄟ~D", password)
	hash := sha256.New()
	hash.Write([]byte(s))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
