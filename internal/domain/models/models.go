package models

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}

type LoginUser struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Book struct {
	BookId int    `json:"b_id"`
	Lable  string `json:"lable" validate:"required"`
	Author string `json:"author" validate:"required"`
	UserId int    `json:"u_id"`
}

type Response struct {
	Message string `json:"message"`
}
