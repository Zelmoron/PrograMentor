package repository

import (
	"errors"

	"main/domain"
	"main/utils"
)

func (repo *Repo) SaveRefreshToken(userID int64, refreshToken string) error {
	refreshTokenRecord := domain.RefreshToken{
		UserID: userID,
		Token:  refreshToken,
	}
	result := repo.db.Create(&refreshTokenRecord)
	return result.Error
}

func (repo *Repo) GetUserIDByRefreshToken(refreshToken string) (int64, error) {
	claims, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return 0, err
	}

	userID, ok := claims["sub"].(int64)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	return userID, nil
}

func (repo *Repo) DeleteRefreshToken(refreshToken string) error {
	result := repo.db.Where("Token = ?", refreshToken).Delete(&domain.RefreshToken{})
	return result.Error
}
