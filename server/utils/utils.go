package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func VerifyPassword(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
