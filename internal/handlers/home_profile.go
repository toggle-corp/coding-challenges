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

func GetChallengesHandler(c *gin.Context, db *gorm.DB, user models.User, templateCtx gin.H) {
	templateCtx["action"] = c.Query("action")
	templateCtx["Challenges"] = models.GetChallenges(db, user.IsAdmin)
	if user.IsAdmin {
		c.HTML(http.StatusOK, "admin_dashboard.html", templateCtx)
	} else {
		c.HTML(http.StatusOK, "user_challenges.html", templateCtx)
	}
}

func ProfileGetHandler(c *gin.Context, db *gorm.DB) {
}
