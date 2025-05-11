package repository

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"main/domain"
)

type Repo struct {
	db          *gorm.DB
	UsersRepo   *UsersRepo
	LessonsRepo *LessonsRepo
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
		db:          db,
		UsersRepo:   NewUsersRepo(db),
		LessonsRepo: NewLessonsRepo(db),
	}
}

func (repo *Repo) Migrate() {

	if err := repo.db.AutoMigrate(
		&domain.Users{},
		&domain.Lessons{},
	); err != nil {
		panic(err)
	}
}
