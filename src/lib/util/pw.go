package util

import (
    "golang.org/x/crypto/bcrypt"
)

// bcrypt salts the password for us
func HashPassword(password string) string {
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
    hashedPw := string(hash) 
    return hashedPw
}

func VerifyPassword(hashedPassword string, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        return false
    }
    return true
}