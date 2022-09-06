package globals

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Secret = []byte("Mysecret")

type HandlerArgs struct {
	GinContext *gin.Context
	DB         *gorm.DB
	GinMap     gin.H
}

type ResponseType string

const (
	HTML ResponseType = "html"
	JSON ResponseType = "json"
)

type HandlerResult struct {
	ResponseType ResponseType
	GinMap       gin.H
	Status       int
	TemplatePath string
}
