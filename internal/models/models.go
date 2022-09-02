package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique"`
	Password  string
	LastLogin time.Time
	IsAdmin   bool `gorm:"default:false"`
	IsGuest   bool `gorm:"default:false"`
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
	Errored SubmissionStatus = "errored"
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
	ID                 int
	Submissions        int
	Title              string
	CreatedBy          string
	IsPublished        bool
	CreatedAt          time.Time
	CorrectSubmissions int
}

func GetChallenges(db *gorm.DB, user User) []ChallengesResult {
	var results []ChallengesResult
	query := `select challenges.id, count(submissions.id) as submissions, title, users.username as created_by,
		challenges.is_published, challenges.created_at, 0
		from challenges
		left join submissions on
		submissions.challenge_id=challenges.id left join users on users.id =
		challenges.created_by `
	if !user.IsAdmin {
		query += " where challenges.is_published = true "
	}
	query += " group by (challenges.id, users.id)"
	query += " order by challenges.id desc"
	db.Raw(query).Scan(&results)
	return results
}

func GetChallengesUser(db *gorm.DB, user User) []ChallengesResult {
	var results []ChallengesResult
	query := fmt.Sprintf(`select challenges.id, count(s.id) as submissions, title, users.username as created_by,
		challenges.is_published, challenges.created_at, array_length(
			array_agg(s.id) filter (WHERE s.status = 'passed' and s.submitted_by=%v),
            1
        ) as correct_submissions
		from challenges
		left join submissions s on s.challenge_id=challenges.id
		left join users on users.id = challenges.created_by
		`, user.ID)
	if !user.IsAdmin {
		query += " where challenges.is_published = true "
	}
	query += " group by (challenges.id, users.id)"
	query += " order by challenges.id desc"
	db.Raw(query).Scan(&results)
	return results
}

type SubmissionsResult struct {
	ID             int
	ChallengeTitle string
	CreatedAt      time.Time
	Language       string
	Status         string
	Score          int
	SubmittedBy    string
	SubmitterId    int
}

func GetSubmissions(db *gorm.DB, user User) []SubmissionsResult {
	var results []SubmissionsResult
	db.Raw(fmt.Sprintf(`select submissions.id,
		challenges.title as challenge_title, submissions.created_at,
		submissions.language, submissions.status, submissions.score,
		users.username as submitted_by, submissions.submitted_by as submitter_id from submissions
		left join challenges on submissions.challenge_id = challenges.id
		left join users on submissions.submitted_by = users.id 
		where submissions.submitted_by=%v
		order by submissions.id desc
	`, user.ID)).Scan(&results)
	return results
}

func GetSubmissionsForChallenge(db *gorm.DB, chId string) []SubmissionsResult {
	var results []SubmissionsResult
	db.Raw(fmt.Sprintf(`select submissions.id,
		challenges.title as challenge_title, submissions.created_at,
		submissions.language, submissions.status, submissions.score,
		users.username as submitted_by, submissions.submitted_by as submitter_id from submissions
		left join challenges on submissions.challenge_id = challenges.id
		left join users on submissions.submitted_by = users.id 
		where challenges.id=%v
		order by submissions.id desc
	`, chId)).Scan(&results)
	return results
}
