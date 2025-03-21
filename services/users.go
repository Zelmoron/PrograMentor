package services

import (
	"crypto/sha256"
	"encoding/hex"

	"main/repository"
)

type Users struct {
	repos *repository.Repo
}

func NewUsers(repos *repository.Repo) *Users {
	return &Users{
		repos: repos,
	}
}

func VerifyPassword(password string, hashedPassword string) bool {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:]) == hashedPassword
}
