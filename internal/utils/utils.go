package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"toggle-corp/coding-challenges/internal/globals"
	"toggle-corp/coding-challenges/internal/models"
)

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

func FormatAsDate(t time.Time) string {
	loc, err := time.LoadLocation("Asia/Kathmandu")
	if err == nil {
		t = t.In(loc)
	}
	return t.Format("02-Jan-2006 03:04:05 PM")
}

func SnakeToTitle(snake string) string {
	splitted := strings.Split(snake, "_")
	titled := Map(splitted, strings.Title)
	return strings.Join(titled, " ")
}

func GetOSEnv(k string, defaultVal string) string {
	if value, ok := os.LookupEnv(k); ok {
		return value
	}
	return defaultVal
}

func ConnectDB() (*gorm.DB, error) {
	user := GetOSEnv("POSTGRES_USER", "postgres")
	password := GetOSEnv("POSTGRES_PASSWORD", "postgres")
	dbname := GetOSEnv("POSTGRES_DB", "coding-challenge")
	port := GetOSEnv("POSTGRES_PORT", "5434")
	hostname := GetOSEnv("POSTGRES_HOSTNAME", "localhost")

	dbConnString := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		hostname, user, password, dbname, port,
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbConnString,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err == nil {
		db.AutoMigrate(&models.User{})
		db.AutoMigrate(&models.Challenge{})
		db.AutoMigrate(&models.Submission{})
	}
	return db, err
}

type DBHandler = func(_ globals.HandlerArgs) globals.HandlerResult
type DBUserHandler = func(_ *gin.Context, db_ *gorm.DB, _ models.User, _ gin.H)

func WithDB(handler DBHandler, db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		args := globals.HandlerArgs{
			GinContext: c,
			DB:         db,
			GinMap:     gin.H{},
		}
		handlerResp := handler(args)
		if handlerResp.ResponseType == globals.HTML {
			c.HTML(
				http.StatusOK,
				handlerResp.TemplatePath,
				handlerResp.GinMap,
			)
		}
	}
}

func WithDBAdmin(handler DBUserHandler, db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		user_raw, exists := c.Get("user")
		if !exists || !user_raw.(models.User).IsAdmin {
			c.Redirect(http.StatusMovedPermanently, "/forbidden")
			c.Abort()
			return
		}
		tctx := gin.H{"user": user_raw.(models.User)}
		handler(c, db, user_raw.(models.User), tctx)
	}
}

func WithDBUser(handler DBUserHandler, db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		user_raw, exists := c.Get("user")
		if !exists {
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.Abort()
			return
		}
		tctx := gin.H{"user": user_raw.(models.User)}
		handler(c, db, user_raw.(models.User), tctx)
	}
}

func HashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

func PasswordsMatch(pwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	if err != nil {
		return false
	}
	return true
}
