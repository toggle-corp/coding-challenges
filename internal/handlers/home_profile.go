package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"toggle-corp/coding-challenges/internal/globals"
	"toggle-corp/coding-challenges/internal/models"
)

func RootHandler(c *gin.Context, db DB, user models.User, templateCtx gin.H) {
	c.HTML(http.StatusOK, "index.html", templateCtx)
}

func ProfileGetHandler(req globals.HandlerArgs) globals.HandlerResult {
	return globals.HandlerResult{}
}
