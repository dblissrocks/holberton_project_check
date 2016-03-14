package main

import (
	"fmt"
	"os"
	"path"
)

// Checks that a given task is as expected
func checkTaskFiles(t task) (checkResult, error) {
	if _, err := os.Stat(t.GithubDir); os.IsNotExist(err) { // Directory doesn't exist
		message := fmt.Sprintf("No \"%s\" directory was found in this directory.", t.GithubDir)
		return checkResult{Passed: false, Message: message}, nil
	}

	if _, err := os.Stat(path.Join(t.GithubDir, t.GithubFile)); os.IsNotExist(err) { // Directory doesn't exist
		message := fmt.Sprintf("No \"%s\" file was found in the \"%s\" directory.", t.GithubFile, t.GithubDir)
		return checkResult{Passed: false, Message: message}, nil
	}

	return checkResult{Passed: true}, nil
}
