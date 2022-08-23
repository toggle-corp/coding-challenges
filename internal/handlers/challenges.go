package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"toggle-corp/coding-challenges/internal/models"
)

func ChallengesGetHandler(c *gin.Context, db *gorm.DB) {
}

func NewChallengeGetHandler(c *gin.Context, db *gorm.DB, admin models.User) {
	c.HTML(http.StatusOK, "new-challenge.html", gin.H{
		"edit": false,
	})
}

func validateChallengeInputs(c *gin.Context) (models.Challenge, map[string]string) {
	title := c.PostForm("title")
	problemStatement := c.PostForm("problem_statement")
	testInputs := c.PostForm("test_inputs")
	testOutputs := c.PostForm("test_outputs")
	score := c.PostForm("score")
	isPublished := c.PostForm("is_published")

	valid := true
	errors := make(map[string]string)
	if title == "" {
		valid = false
		errors["title"] = "Cannot be empty"
	}
	if problemStatement == "" {
		valid = false
		errors["problem_statement"] = "Cannot be empty"
	}
	if testInputs == "" {
		valid = false
		errors["test_inputs"] = "Cannot be empty"
	}
	if testOutputs == "" {
		valid = false
		errors["test_outputs"] = "Cannot be empty"
	}
	// Check if numeric
	scoreVal, err := strconv.Atoi(score)
	if err != nil || scoreVal < 0 {
		valid = false
		errors["score"] = "Invalid value"
	}
	var challenge models.Challenge
	if !valid {
		return challenge, errors
	}
	var published bool
	if isPublished == "true" {
		published = true
	} else {
		published = false
	}
	challenge = models.Challenge{
		Title:            title,
		ProblemStatement: problemStatement,
		TestInputs:       testInputs,
		TestOutputs:      testOutputs,
		IsPublished:      published,
		Score:            scoreVal,
	}
	return challenge, nil
}

func NewChallengePostHandler(c *gin.Context, db *gorm.DB, admin models.User) {
	challenge, errors := validateChallengeInputs(c)
	if errors != nil {
		c.HTML(http.StatusBadRequest, "new-challenge.html", gin.H{
			"errors": errors,
		})
	}
	// Create challenge
	challenge.CreatedBy = int(admin.ID)
	result := db.Create(&challenge)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "new-challenge.html", gin.H{
			"error": "Could not create challenge",
			// TODO: serialize to map[string]string
			"form": gin.H{
				"title":             challenge.Title,
				"problem_statement": challenge.ProblemStatement,
				"test_inputs":       challenge.TestInputs,
				"test_outputs":      challenge.TestOutputs,
				"score":             challenge.Score,
				"publish":           challenge.IsPublished,
			},
		})
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/home")
	c.Abort()
}

func EditChallengeGetHandler(c *gin.Context, db *gorm.DB, admin models.User) {
	c.HTML(http.StatusOK, "new-challenge.html", gin.H{
		"edit": false,
	})
}

func EditChallengePutHandler(c *gin.Context, db *gorm.DB, admin models.User) {
}
