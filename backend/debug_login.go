package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    _ "github.com/lib/pq"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    // Connect to DB
    connStr := "host=localhost port=5432 user=postgres password=postgres dbname=akinweb sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Test 1: Query for user
    username := "Admin"
    var dbUsername, dbHash string
    var isAdmin bool
    
    err = db.QueryRow("SELECT username, password_hash, is_admin FROM users WHERE username = $1", username).
        Scan(&dbUsername, &dbHash, &isAdmin)
    
    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Printf("❌ User '%s' not found in database\n", username)
        } else {
            fmt.Printf("❌ Query error: %v\n", err)
        }
        return
    }
    
    fmt.Printf("✅ Found user: %s, is_admin: %v\n", dbUsername, isAdmin)
    fmt.Printf("Hash starts with: %.30s...\n", dbHash)
    
    // Test 2: Verify password
    password := "Polymath Here"
    err = bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(password))
    if err != nil {
        fmt.Printf("❌ Password verification failed: %v\n", err)
    } else {
        fmt.Println("✅ Password verification SUCCESS!")
    }
}
