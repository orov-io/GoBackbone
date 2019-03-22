package service

import "github.com/gin-gonic/gin"

func (a *App) pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
