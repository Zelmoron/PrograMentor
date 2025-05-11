package domain

type Lessons struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	LessonTitle string `gorm:"type:varchar(255);not null" json:"lesson_title"`
	Content     string `gorm:"type:text;not null" json:"content"`
	Examples    string `gorm:"type:text" json:"examples"`
}
