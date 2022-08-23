package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
    "gorm.io/gorm"
	"log"
	"net/http"

    "toggle-corp/coding-challenges/internal/models"
)

func AuthRequired (db *gorm.DB) func(_ *gin.Context) {
    return func (c *gin.Context) {
        session := sessions.Default(c)
        userid := session.Get("userid")
        if userid == nil {
            c.Redirect(http.StatusMovedPermanently, "/login")
            c.Abort()
            return
        }
        // Get user from db
        var user models.User
        result := db.First(&user, int(userid.(uint)))

        if result.RowsAffected == 0 {
            log.Println("Session present but not user", userid)
            c.Redirect(http.StatusMovedPermanently, "/login")
            c.Abort()
            return
        }
        c.Set("user", user)
        c.Next()
    }
}
