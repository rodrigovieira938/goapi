package auth

type Response struct {
	Token string `json:"token"`
}
type Login struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}
