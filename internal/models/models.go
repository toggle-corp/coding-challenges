package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string    `json:"username";gorm:"unique"`
	Password  string    `json:"password"`
	LastLogin time.Time `json:"lastLogin"`
	IsAdmin   bool      `gorm:"default:false"`
	IsGuest   bool      `gorm:"default:false"`
}

type Challenge struct {
	gorm.Model
	Title            string
	ProblemStatement string
	TestInputs       string
	TestOutputs      string
	IsPublished      bool `gorm:"default:true"`
	Score            int  `gorm:"default:10"`
	CreatedBy        int
	User             User `gorm:"foreignKey:CreatedBy"`
}

type SubmissionLanguage string

const (
	Python     SubmissionLanguage = "python"
	Javascript SubmissionLanguage = "javascript"
	Go         SubmissionLanguage = "go"
)

type SubmissionStatus string

const (
	InQueue SubmissionStatus = "in_queue"
	Running SubmissionStatus = "running"
	Failed  SubmissionStatus = "failed"
	Passed  SubmissionStatus = "passed"
)

type Submission struct {
	gorm.Model
	Language      SubmissionLanguage
	SubmittedCode string
	Error         string `gorm:"default:null"`
	Score         int    `gorm:"default:0"`
	Status        SubmissionStatus

	ChallengeID int
	Challenge   Challenge

	SubmittedBy int
	User        User `gorm:"foreignKey:SubmittedBy"`
}

type ChallengesResult struct {
	ID          int
	Submissions int
	Title       string
	CreatedBy   string
	CreatedAt   time.Time
}

func GetChallenges(db *gorm.DB) []ChallengesResult {
	var results []ChallengesResult
	db.Raw(`select challenges.id, count(submissions.id) as submissions, title, users.username as created_by,
		challenges.created_at from challenges left join submissions on
		submissions.challenge_id=challenges.id left join users on users.id =
		challenges.created_by group by (challenges.id, users.id)`).Scan(&results)
	return results
}
