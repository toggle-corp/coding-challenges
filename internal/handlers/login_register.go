package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"toggle-corp/coding-challenges/internal/globals"
	"toggle-corp/coding-challenges/internal/models"
	"toggle-corp/coding-challenges/internal/utils"
)

type DB = *gorm.DB

func LoginGetHandler(c *gin.Context) {
	action := c.Query("action")
	var message string
	if action == "success" {
		message = "Successful Registration"
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":   "Login",
		"message": message,
	})
}

func RegisterGetHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func RegisterHandler(req globals.HandlerArgs) globals.HandlerResult {
	// Get registration data
	c := req.GinContext
	db := req.DB
	username := c.PostForm("username")
	password := c.PostForm("password")
	confirm_password := c.PostForm("confirm_password")
	errors, valid := utils.ValidateRegistration(password, confirm_password)
	// Query user
	userQuery := models.User{Username: username}
	var user models.User
	result := db.Where(userQuery).First(&user)
	if result.RowsAffected > 0 {
		valid = false
		errors["username"] = "User already exists"
	}
	// Hash password
	hashed, err := utils.HashAndSalt(password)
	if err != nil {
		errors["error"] = "Something went wrong"
	}
	if !valid {
		return globals.HandlerResult{
			ResponseType: globals.HTML,
			GinMap: gin.H{
				"error":                  errors["error"],
				"username_error":         errors["username"],
				"password_error":         errors["password"],
				"confirm_password_error": errors["confirm_password"],
				"form": gin.H{
					"username":         username,
					"password":         password,
					"confirm_password": confirm_password,
				},
			},
			Status:       http.StatusBadRequest,
			TemplatePath: "register.html",
		}
	}
	newuser := models.User{Username: username, Password: hashed}
	db.Create(&newuser)
	c.Redirect(http.StatusMovedPermanently, "/login?action=success")
	c.Abort()
	return globals.HandlerResult{}
}

func LoginHandler(req globals.HandlerArgs) globals.HandlerResult {
	c := req.GinContext
	db := req.DB
	username := c.PostForm("username")
	password := c.PostForm("password")

	errors := make(map[string]string)
	valid := true

	if username == "" {
		valid = false
		errors["username_error"] = "Cannot be empty"
	}
	if password == "" {
		valid = false
		errors["password_error"] = "Cannot be empty"
	}
	// Get user
	userQuery := models.User{Username: username}
	var user models.User
	result := db.Where(userQuery).First(&user)
	if result.RowsAffected == 0 {
		valid = false
		errors["error"] = "Invalid credentials"
	}
	// Check password
	if !utils.PasswordsMatch(password, user.Password) {
		valid = false
		errors["error"] = "Invalid credentials"
	}
	if !valid {
		return globals.HandlerResult{
			ResponseType: globals.HTML,
			GinMap: gin.H{
				"error":          errors["error"],
				"username_error": errors["username_error"],
				"password_error": errors["password_error"],
			},
			Status:       http.StatusBadRequest,
			TemplatePath: "login.html",
		}
	}
	// Set session
	session := sessions.Default(c)
	session.Set("userid", user.ID)
	session.Save()
	c.Redirect(http.StatusMovedPermanently, "/challenges")
	return globals.HandlerResult{}
}

func LogoutHandler(req globals.HandlerArgs) globals.HandlerResult {
	c := req.GinContext
	session := sessions.Default(c)
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
	c.Redirect(http.StatusMovedPermanently, "/login")
	return globals.HandlerResult{}
}
