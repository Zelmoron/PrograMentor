package handlers

import (
	"main/repository"
	"main/services"
)

type (
	OutHandlers struct {
		repos *repository.Repo
		users *services.Users
	}
	InHandlers struct {
		repos *repository.Repo
		users *services.Users
	}
)

func (o *OutHandlers) GetUsers() *services.Users {
	return o.users
}

func NewOutHandlers(
	repos *repository.Repo,
	users *services.Users,
) *OutHandlers {
	return &OutHandlers{
		repos: repos,
		users: users,
	}
}

func NewInHandlers(
	repos *repository.Repo,
	users *services.Users) *InHandlers {
	return &InHandlers{
		repos: repos,
		users: users,
	}
}
