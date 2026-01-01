package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "Polymath Here"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
    fmt.Printf("UPDATE users SET password_hash = '%s' WHERE username = 'Admin';\n", hash)
}
