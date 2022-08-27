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

func writeTmpcode(subs *models.Submission, name string, content string) string {
	tmpFname := fmt.Sprintf("/tmp/%v_%v", name, subs.ID)
	err := os.WriteFile(tmpFname, []byte(content), 0777)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return tmpFname
}

func errorSubscription(subs *models.Submission, errmsg string, db *gorm.DB) {
	subs.Status = models.Errored
	subs.Error = errmsg
	// TODO: check if db error occurs
	db.Model(&subs).Updates(subs)
}

func failSubscription(subs *models.Submission, errmsg string, db *gorm.DB) {
	subs.Status = models.Failed
	subs.Error = errmsg
	// TODO: check if db error occurs
	db.Model(&subs).Updates(subs)
}

func Execute(subs models.Submission, db *gorm.DB) {
	var runnerScript = "./scripts/runner.sh"
	var challenge models.Challenge
	result := db.First(&challenge, subs.ChallengeID)

	if result.Error != nil || result.RowsAffected == 0 {
		errMsg := fmt.Sprintf("Challenge %v not found for submission %v", challenge.ID, subs.ID)
		log.Fatal(errMsg)
		errorSubscription(&subs, errMsg, db)
		return
	}
	tmpCode := writeTmpcode(&subs, "code", subs.SubmittedCode)
	// TODO: tmpTest/Outs need not be created for each submission. They are
	// common for a particular challange
	tmpTest := writeTmpcode(&subs, "test", challenge.TestInputs)
	tmpOuts := writeTmpcode(&subs, "out", challenge.TestOutputs)

	if tmpCode == "" || tmpTest == "" || tmpOuts == "" {
		errMsg := "Could not create tmp files"
		log.Fatal(errMsg, " Subscription: ", subs.ID)
		errorSubscription(&subs, errMsg, db)
		return
	}

	defer func() {
		rmcmd := exec.Command("rm", tmpCode, tmpTest, tmpOuts)
		runAndGetError(rmcmd)
	}()

	var cmd *exec.Cmd
	if GetOSEnv("GIN_MODE", "local") == "release" {
		cmd = exec.Command("ssh", "runner@runnerenv", runnerScript, string(subs.Language), fmt.Sprint(subs.ID))
	} else {
		// Do not ssh, directly call
		cmd = exec.Command(runnerScript, string(subs.Language), fmt.Sprint(subs.ID))
	}
	cmdCode, cmderr := runAndGetError(cmd)
	if cmdCode == 0 {
		log.Println("Successful execution")
		subs.Status = models.Passed
		subs.Score = challenge.Score
		// TODO: check if db error occurs
		db.Model(&subs).Updates(subs)
	} else if cmdCode == 43 { // 43 is custom
		// means test failed
		failSubscription(&subs, string(cmderr), db)
	} else {
		errorSubscription(&subs, string(cmderr), db)
	}
}

func runAndGetError(cmd *exec.Cmd) (int, string) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		errstr := stderr.String()
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), errstr
		}
		return 1, errstr
	}
	return 0, ""
}
