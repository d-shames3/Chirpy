package auth

import (
	"log"
	"testing"
)

func TestAuth(t *testing.T) {
	password := "test"
	hash, err := HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("hashed password: %v", hash)

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		log.Fatal(err)
	}

	if !match {
		t.Errorf(`HashPassword("test") = %v, want true for match`, hash)
	}
}
