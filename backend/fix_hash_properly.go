package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "Polymath Here"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
    
    // Print SQL with proper escaping for bash
    fmt.Println("-- Copy and run this in psql:")
    fmt.Printf("UPDATE users SET password_hash = '%s' WHERE username = 'Admin';\n", string(hash))
    
    // Or run directly
    fmt.Println("\n-- Or run this bash command:")
    fmt.Printf("sudo -u postgres psql -d akinweb -c \"UPDATE users SET password_hash = '%s' WHERE username = 'Admin';\"\n", string(hash))
}
