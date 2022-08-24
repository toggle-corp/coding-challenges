package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"toggle-corp/coding-challenges/internal/models"
	"toggle-corp/coding-challenges/internal/utils"
)

type Error = map[string]string

func ChallengesGetHandler(c *gin.Context, db *gorm.DB, user models.User, templateCtx gin.H) {
	cid := c.Param("id")
	var ch models.Challenge
	result := db.First(&ch, cid)

	if result.Error != nil || result.RowsAffected == 0 {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}
	templateCtx["title"] = "Challenge: " + ch.Title
	templateCtx["challenge"] = ch
	c.HTML(http.StatusOK, "user_challenge.html", templateCtx)
}

func validateSubmission(c *gin.Context) (models.Submission, Error) {
	submission := models.Submission{
		Error:  "",
		Status: models.InQueue,
	}
	errors := make(Error)
	valid := true
	lang := models.SubmissionLanguage(c.PostForm("Language"))
	if lang != models.Python && lang != models.Javascript && lang != models.Go {
		valid = false
		errors["Language"] = "Invalid language"
	}
	submission.Language = models.SubmissionLanguage(lang)
	submittedCode := c.PostForm("SubmittedCode")
	if strings.Trim(submittedCode, " \n\r") == "" {
		valid = false
		errors["SubmittedCode"] = "Code cannot be empty"
	}
	// Replace tab with spaces
	submission.SubmittedCode = strings.Replace(submittedCode, "\t", "    ", -1)
	if !valid {
		return submission, errors
	}
	return submission, nil
}

func ChallengesPostHandler(c *gin.Context, db *gorm.DB, user models.User, templateCtx gin.H) {
	cid := c.Param("id")
	var ch models.Challenge
	result := db.First(&ch, cid)

	if result.Error != nil || result.RowsAffected == 0 {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}
	templateCtx["title"] = "Challenge: " + ch.Title
	templateCtx["challenge"] = ch

	submission, errors := validateSubmission(c)
	if errors != nil {
		templateCtx["errors"] = errors
		templateCtx["submisssion"] = submission
		c.HTML(http.StatusBadRequest, "user_challenge.html", templateCtx)
		return
	}
	// Add challenge id and user id
	submission.ChallengeID = int(ch.ID)
	submission.SubmittedBy = int(user.ID)

	// Save submission
	subResult := db.Create(&submission)
	if subResult.Error != nil {
		templateCtx["error"] = "Could not submit"
		templateCtx["submission"] = submission
		c.HTML(http.StatusBadRequest, "user_challenge.html", templateCtx)
		return
	}
	go utils.Execute(submission, db)
	c.Redirect(http.StatusMovedPermanently, "/my-submissions?action=submitted")
	c.Abort()
}

func NewChallengeGetHandler(c *gin.Context, db *gorm.DB, admin models.User, templateCtx gin.H) {
	var challenge models.Challenge
	c.HTML(http.StatusOK, "new-challenge.html", gin.H{
		"edit":      false,
		"action":    "/new-challenge",
		"method":    "post",
		"title":     "New Challenge",
		"challenge": challenge,
		"user":      templateCtx["user"],
	})
}

func validateChallengeInputs(c *gin.Context) (models.Challenge, Error) {
	title := c.PostForm("Title")
	problemStatement := c.PostForm("ProblemStatement")
	testInputs := c.PostForm("TestInputs")
	testOutputs := c.PostForm("TestOutputs")
	score := c.PostForm("Score")
	isPublished := c.PostForm("IsPublished")

	valid := true
	errors := make(Error)
	if title == "" {
		valid = false
		errors["Title"] = "Cannot be empty"
	}
	if problemStatement == "" {
		valid = false
		errors["ProblemStatement"] = "Cannot be empty"
	}
	if testInputs == "" {
		valid = false
		errors["TestInputs"] = "Cannot be empty"
	}
	if testOutputs == "" {
		valid = false
		errors["TestOutputs"] = "Cannot be empty"
	}
	// Check if numeric
	scoreVal, err := strconv.Atoi(score)
	if err != nil || scoreVal < 0 {
		valid = false
		errors["Score"] = "Invalid value"
	}
	var challenge models.Challenge
	if !valid {
		return challenge, errors
	}
	var published bool
	if isPublished != "" {
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

func NewChallengePostHandler(c *gin.Context, db *gorm.DB, admin models.User, templateCtx gin.H) {
	challenge, errors := validateChallengeInputs(c)
	templateCtx["method"] = "post"
	templateCtx["title"] = "New Challenge"
	templateCtx["action"] = "/new-challenge"
	if errors != nil {
		templateCtx["errors"] = errors
		templateCtx["challenge"] = challenge
		c.HTML(http.StatusBadRequest, "new-challenge.html", templateCtx)
		return
	}
	// Create challenge
	challenge.CreatedBy = int(admin.ID)
	result := db.Create(&challenge)
	if result.Error != nil {
		templateCtx["error"] = "Could not create challenge"
		c.HTML(http.StatusBadRequest, "new-challenge.html", templateCtx)
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/challenges?action=create")
	c.Abort()
}

func EditChallengeGetHandler(c *gin.Context, db *gorm.DB, admin models.User, templateCtx gin.H) {
	cid := c.Param("id")
	var ch models.Challenge
	result := db.First(&ch, cid)
	if result.Error != nil || result.RowsAffected == 0 {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}
	tCtx := gin.H{
		"edit":      true,
		"action":    "/edit-challenge/" + cid,
		"method":    "post",
		"title":     "Edit Challenge",
		"challenge": ch,
		"user":      templateCtx["user"],
	}
	c.HTML(http.StatusOK, "new-challenge.html", tCtx)
}

func EditChallengePostHandler(c *gin.Context, db *gorm.DB, admin models.User, templateCtx gin.H) {
	cid := c.Param("id")
	var ch models.Challenge
	result := db.First(&ch, cid)
	if result.Error != nil || result.RowsAffected == 0 {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}
	templateCtx = gin.H{
		"edit":   true,
		"action": "/edit-challenge/" + cid,
		"method": "post",
		"title":  "Edit Challenge",
		"user":   templateCtx["user"],
	}
	challenge, errors := validateChallengeInputs(c)
	if errors != nil {
		templateCtx["errors"] = errors
		templateCtx["challenge"] = ch
		c.HTML(http.StatusBadRequest, "new-challenge.html", templateCtx)
	}
	// do update
	res := db.Model(&ch).Select("IsPublished", "Title", "ProblemStatement", "TestInputs", "TestOutputs", "Score").Updates(challenge)
	if res.Error != nil {
		templateCtx["error"] = "Cound not update"
		templateCtx["challenge"] = challenge
		c.HTML(http.StatusBadRequest, "new-challenge.html", templateCtx)
		return
	}
	// Redirect to challenges with a message
	c.Redirect(http.StatusMovedPermanently, "/challenges?action=update")
	c.Abort()
}

func MySubmissionsGetHandler(c *gin.Context, db *gorm.DB, user models.User, templateCtx gin.H) {
	submissions := models.GetSubmissions(db, user)
	action := c.Query("action")
	var message string
	if action == "submitted" {
		message = "Successfully submitted solution"
	}
	templateCtx["submissions"] = submissions
	templateCtx["message"] = message
	c.HTML(http.StatusOK, "submissions.html", templateCtx)
}

func SubmissionGetHandler(c *gin.Context, db *gorm.DB, user models.User, templateCtx gin.H) {
	// show details for a particular submission
	sid := c.Param("id")
	var sb models.Submission
	result := db.Preload("Challenge").First(&sb, sid)
	if result.Error != nil || result.RowsAffected == 0 {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}
	templateCtx["submission"] = sb
	fmt.Println(sid, sb.Status)
	c.HTML(http.StatusOK, "submission.html", templateCtx)
}

func SubmissionsGetHandler(c *gin.Context, db *gorm.DB, user models.User, templateCtx gin.H) {
}
