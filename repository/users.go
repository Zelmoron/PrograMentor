package repository

import (
	"gorm.io/gorm"

	"main/domain"
)

type UsersRepo struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (ur *UsersRepo) GetUserByUsername(username string) (*domain.Users, error) {
	var user domain.Users
	result := ur.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &user, nil
}
