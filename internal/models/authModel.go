package models

type SignInDto struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
