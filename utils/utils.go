package utils

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func EncryptStr(str string) (string, error) {

	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		return "", err
	}
	fmt.Println("Hash Generated:", string(hash))
	// Store this "hash" somewhere, e.g. in your database
	return string(hash), nil
}

func CheckStrHash(str, hash string) bool {
	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str)); err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
		return false
	}
	log.Println("Passwords match: ", hash)
	return true
}
