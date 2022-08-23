package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"toggle-corp/coding-challenges/internal/models"
)

func RootHandler(c *gin.Context, db DB) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func HomeHandler(c *gin.Context, db *gorm.DB) {
	user_raw, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusMovedPermanently, "/login")
		c.Abort()
		return
	}
	ctxValues := make(map[string]interface{})
	user := user_raw.(models.User)
	ctxValues["user"] = user
	if user.IsAdmin {
		ctxValues["Challenges"] = models.GetChallenges(db)
		c.HTML(http.StatusOK, "admin_dashboard.html", ctxValues)
	} else {
		c.HTML(http.StatusOK, "home.html", ctxValues)
	}
}

func ProfileGetHandler(c *gin.Context, db *gorm.DB) {
}
