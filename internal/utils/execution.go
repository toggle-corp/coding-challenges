package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"toggle-corp/coding-challenges/internal/models"

	"os/exec"

	"gorm.io/gorm"
)

func Execute(subs models.Submission, db *gorm.DB) {
	var runnerScript = "./scripts/runner.sh"
	// Create a tmp file
	code := []byte(subs.SubmittedCode)
	tmpFname := fmt.Sprintf("/tmp/code_%v", subs.ID)
	err := os.WriteFile(tmpFname, code, 0444)
	if err != nil {
		log.Fatal(err)
		subs.Status = models.Errored
		// TODO: check if db error occurs
		db.Model(&subs).Updates(subs)
	}
	cmd := exec.Command(runnerScript, string(subs.Language), tmpFname)
	cmdCode, cmderr := runAndGetError(cmd)
	if cmdCode == 0 {
		log.Println("Successful execution")
		// Remove tmp file
		exec.Command("rm", tmpFname).Run()
		subs.Status = models.Passed
		// TODO: check if db error occurs
		db.Model(&subs).Updates(subs)
	} else {
		subs.Error = string(cmderr)
		subs.Status = models.Errored
		// TODO: check if db error occurs
		db.Model(&subs).Updates(subs)
	}
}

func runAndGetError(cmd *exec.Cmd) (int, string) {
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error when running command.  Error log:")
		errstr := stderr.String()
		fmt.Printf("Got command status: %s\n", err.Error())
		return 1, errstr
	}
	return 0, ""
}
