package types

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Message string `json:"message"`
}
