package dto

type LoginRequest struct {
	UserName     string `json:"username"`
	UserPassword string `json:"password"`
}
