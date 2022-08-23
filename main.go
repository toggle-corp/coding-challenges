package main

import (
	"fmt"
	"html/template"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"toggle-corp/coding-challenges/internal/globals"
	"toggle-corp/coding-challenges/internal/middleware"
	"toggle-corp/coding-challenges/internal/routers"
	"toggle-corp/coding-challenges/internal/utils"
)

func main() {
	db, err := utils.ConnectDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"formatAsDate": utils.FormatAsDate,
	})

	r.LoadHTMLGlob("templates/*.html")
	r.Delims("{{", "}}")

	r.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))
	public := r.Group("/")
	routers.PublicRoutes(public, db)

	private := r.Group("/")
	private.Use(middleware.AuthRequired(db))
	routers.PrivateRoutes(private, db)

	r.Run()
}
