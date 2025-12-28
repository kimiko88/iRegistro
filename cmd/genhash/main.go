package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "password123"
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	fmt.Println(string(bytes))
}
