package services

import "main/repository"

type Users struct {
	repos *repository.Repo
}

func NewUsers(repos *repository.Repo) *Users {
	return &Users{
		repos: repos,
	}
}
