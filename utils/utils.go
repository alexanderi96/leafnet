package utils

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func EncryptStr(str string) (string, error) {

	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	fmt.Println("Hash Generated")
	return string(hash), nil
}

func CheckStrHash(str, hash string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str)); err != nil {
		return false, err
	}
	log.Println("String matches")
	return true, nil
}

func CheckPathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
