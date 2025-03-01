package domain

type Users struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	//Email     string `json:"email"` - потом сделать
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
