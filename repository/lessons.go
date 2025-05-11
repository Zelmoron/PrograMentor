package repository

import (
	"main/domain"

	"gorm.io/gorm"
)

type LessonsRepo struct {
	db *gorm.DB
}

func NewLessonsRepo(db *gorm.DB) *LessonsRepo {
	return &LessonsRepo{db: db}
}

func (lr *LessonsRepo) GetLessonByID(id string) (*domain.Lessons, error) {
	var lesson domain.Lessons
	result := lr.db.First(&lesson, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &lesson, nil
}

func (lr *LessonsRepo) GetTotalLessonsCount() (int64, error) {
	var count int64
	result := lr.db.Model(&domain.Lessons{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
