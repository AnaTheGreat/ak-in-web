package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "Polymath Here"
    
    // Generate new hash
    newHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
    fmt.Printf("New hash for 'Polymath Here':\n%s\n\n", newHash)
    
    // Get hash from database (you'll need to copy this from the SELECT output)
    dbHash := "$2a$10$Yj5MerGU2TWplVfwWvoXCua65lbRXU7lJ6luTwCKn3jqocVQlJaPe"
    
    // Test if password matches
    err := bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(password))
    if err != nil {
        fmt.Printf("Password verification FAILED: %v\n", err)
    } else {
        fmt.Println("Password verification SUCCESS!")
    }
}
