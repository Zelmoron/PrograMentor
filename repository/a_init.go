package repository

import (
	"fmt"
	"main/domain"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repo struct {
	db        *gorm.DB
	UsersRepo *UsersRepo
}

func InitRepo(connectionString string) *Repo {

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	fmt.Printf("Repository stats: %+v", sqlDB.Stats())

	return &Repo{
		db:        db,
		UsersRepo: NewUsersRepo(db),
	}
}

func (repo *Repo) Migrate() {

	if err := repo.db.AutoMigrate(
		&domain.Users{},
		&domain.RefreshToken{},
	); err != nil {
		panic(err)
	}
}

func (repo *Repo) SaveRefreshToken(userID int64, refreshToken string) error {
	refreshTokenRecord := domain.RefreshToken{
		UserID: userID,
		Token:  refreshToken,
	}
	result := repo.db.Create(&refreshTokenRecord)
	return result.Error
}

func (repo *Repo) GetUserIDByRefreshToken(refreshToken string) (uint, error) {
	var refreshTokenRecord domain.RefreshToken
	result := repo.db.Where("Token = ?", refreshToken).First(&refreshTokenRecord)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("refresh token not found")
		}
		return 0, result.Error
	}
	return uint(refreshTokenRecord.UserID), nil
}

func (repo *Repo) DeleteRefreshToken(refreshToken string) error {
	result := repo.db.Where("Token = ?", refreshToken).Delete(&domain.RefreshToken{})
	return result.Error
}
