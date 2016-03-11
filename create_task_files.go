package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Checks that a given task is as expected
func createTaskFiles(t task) (checkResult, error) {
	// Checks if in a git repo
	rootPath, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		message := fmt.Sprintf("Current directory does not seem to be a git repository.\n")
		return checkResult{Passed: false, Message: message}, nil
	}

	dirPath := path.Join(strings.TrimRight(string(rootPath), "\n"), strings.TrimRight(t.GithubDir, "\n"))
	filePath := path.Join(dirPath, t.GithubFile)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) { // Directory doesn't exist

		if err = os.MkdirAll(dirPath, 0755); err != nil {
			message := fmt.Sprintf("%s", err)
			return checkResult{Passed: false, Message: message}, nil
		}
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) { // File doesn't exist
		// Creates the file
		if _, err = os.Create(filePath); err != nil {
			message := fmt.Sprintf("%s", err)
			return checkResult{Passed: false, Message: message}, nil
		}
	}

	return checkResult{Passed: true}, nil
}
