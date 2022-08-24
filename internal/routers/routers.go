package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	h "toggle-corp/coding-challenges/internal/handlers"
	"toggle-corp/coding-challenges/internal/utils"
)

var withDB = utils.WithDB
var withDBAdmin = utils.WithDBAdmin
var withDBUser = utils.WithDBUser

func PublicRoutes(g *gin.RouterGroup, db *gorm.DB) {
	g.GET("/login", h.LoginGetHandler)
	g.POST("/login", withDB(h.LoginHandler, db))
	g.GET("/register", h.RegisterGetHandler)
	g.POST("/register", withDB(h.RegisterHandler, db))
	g.POST("/forbidden", func(c *gin.Context) {
		c.HTML(http.StatusForbidden, "forbidden.html", gin.H{})
	})
}

func PrivateRoutes(g *gin.RouterGroup, db *gorm.DB) {
	g.GET("/", withDBUser(h.RootHandler, db))
	g.GET("/challenges", withDBUser(h.GetChallengesHandler, db))
	g.GET("/challenge/:id", withDBUser(h.ChallengesGetHandler, db))
	g.POST("/challenge/:id", withDBUser(h.ChallengesPostHandler, db))
	g.GET("/my-submissions", withDBUser(h.MySubmissionsGetHandler, db))
	g.GET("/submission/:id", withDBUser(h.SubmissionGetHandler, db))
	g.POST("/submissions/:submissionId", withDBUser(h.SubmissionsGetHandler, db))
	g.GET("/new-challenge", withDBAdmin(h.NewChallengeGetHandler, db))
	g.POST("/new-challenge", withDBAdmin(h.NewChallengePostHandler, db))
	g.GET("/edit-challenge/:id", withDBAdmin(h.EditChallengeGetHandler, db))
	g.POST("/edit-challenge/:id", withDBAdmin(h.EditChallengePostHandler, db))
	g.GET("/profile", withDB(h.ProfileGetHandler, db))
	g.GET("/logout", withDB(h.LogoutHandler, db))
}
