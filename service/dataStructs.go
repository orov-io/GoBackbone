package service

type pongData struct {
	Name string `json:"name" form:"name" binding:"required"`
}
