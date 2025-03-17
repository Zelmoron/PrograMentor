package services

import (
	"log"
	"main/domain"
	"main/repository"
	"os"
)

type Config struct {
	JWTSecret string
}

type Users struct {
	repos *repository.Repo
	Cfg   *Config
}

func NewUsers(repo *repository.Repo) *Users {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in .env file")
	}

	cfg := &Config{
		JWTSecret: jwtSecret,
	}

	return &Users{
		repos: repo,
		Cfg:   cfg,
	}
}

func (s *Users) GetUserByUsername(username string) (*domain.Users, error) {
	return s.repos.UsersRepo.GetUserByUsername(username)
}

func (s *Users) SaveRefreshToken(userID int64, refreshToken string) error {
	return s.repos.SaveRefreshToken(userID, refreshToken)
}

func (s *Users) GetUserIDByRefreshToken(refreshToken string) (uint, error) {
	return s.repos.GetUserIDByRefreshToken(refreshToken)
}

func (s *Users) DeleteRefreshToken(refreshToken string) error {
	return s.repos.DeleteRefreshToken(refreshToken)
}
