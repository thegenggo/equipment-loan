package httperr

import "github.com/gin-gonic/gin"

type Response struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func Write(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, Response{Error: code, Message: message})
}
