package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    hash := "$2a$10$fy9fAhtzFlIBaF8oagGoqO1cqDTnckYCO/tGSCfxgAApNgfj5q2b6"
    password := "Polymath Here"
    
    fmt.Printf("Testing hash: %.30s...\n", hash)
    fmt.Printf("With password: '%s'\n", password)
    
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if err != nil {
        fmt.Printf("❌ FAILED: %v\n", err)
    } else {
        fmt.Println("✅ SUCCESS - This hash is correct!")
    }
}
