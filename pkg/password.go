package pkg

import (
	"golang.org/x/crypto/bcrypt"
)

// GenSaltPassword 将指定字符串，进行bcrypt编码，用于密码加密
func GenSaltPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword 检查 用户输入的密码 与 db的hashPassword，是否一致
func CheckPassword(password, hashPassword string) bool {
	// 密码是bcrypt加密的，所以要先对password加密，然后再与db比对
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}
