package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"toggle-corp/coding-challenges/internal/models"
)

func RootHandler(c *gin.Context, db DB, user models.User, templateCtx gin.H) {
	c.HTML(http.StatusOK, "index.html", templateCtx)
}

func ProfileGetHandler(c *gin.Context, db *gorm.DB) {
}
