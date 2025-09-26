package users

type User struct {
	ID       int    `json:"id" validate:"gte=1"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
