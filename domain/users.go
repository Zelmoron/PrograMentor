package domain

type Users struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type RefreshToken struct {
	Token  string `gorm:"primaryKey" json:"token"`
	UserID int64  `json:"user_id"`
	User   Users  `gorm:"foreignKey:UserID"`
}
