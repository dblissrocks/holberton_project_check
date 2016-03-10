package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func createDirectory(dir string) error {
	// Create new directory
	if err := os.Mkdir(dir, 0755); err != nil {
		return err
	}
	return nil
}

func createFile(file string) error {
	if _, err := os.Create(file); err != nil {
		return err
	}
	return nil
}

// Checks that a given task is as expected
func createTaskFiles(t task) (checkResult, error) {
	// Checks if in a git repo
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		message := fmt.Sprintf("Current directory does not seem to be a git repository.\n")
		return checkResult{Passed: false, Message: message}, nil
	}

	// Checks for the specified repo
	if t.GithubRepo != "" {
		// Obtains the name of current repo
		repoPath := strings.Split(string(out), "/")
		depth := len(repoPath) - 1
		repo := strings.TrimRight(repoPath[depth], "\n")

		// Determines if run from correct repository
		if repo != t.GithubRepo {
			message := fmt.Sprintf("No remote pointing to a \"%s\" Github repository; make sure to run this application from your project's directory.", t.GithubRepo)
			return checkResult{Passed: false, Message: message}, nil
		}

		// Obtains the current directory name
		currPath, err := os.Getwd()
		if err != nil {
			message := fmt.Sprintf("Error: %s\n", err)
			return checkResult{Passed: false, Message: message}, nil
		}

		currDir := strings.Split(currPath, "/")
		depth = len(currDir) - 1
		dir := strings.TrimRight(currDir[depth], "\n")

		// Checks if in the root of the repo
		for dir != t.GithubRepo {
			// Move to the parent directory
			if err := os.Chdir(".."); err != nil {
				message := fmt.Sprintf("%s", err)
				return checkResult{Passed: false, Message: message}, nil
			}
		}
	}

	// Determines depth of project directories
	dir := strings.Split(t.GithubDir, "/")
	dirDepth := len(dir)

	if _, err := os.Stat(t.GithubDir); os.IsNotExist(err) { // Directory doesn't exist
		for i := 0; i < dirDepth; i++ {
			//Creates the new directory
			if err = createDirectory(dir[i]); err != nil {
				message := fmt.Sprintf("%s", err)
				return checkResult{Passed: false, Message: message}, nil
			}

			// Moves to the directory
			if err := os.Chdir(dir[i]); err != nil {
				message := fmt.Sprintf("%s", err)
				return checkResult{Passed: false, Message: message}, nil
			}
		}

		// Returns to the repo root
		for i := 0; i < dirDepth; i++ {
			if err := os.Chdir(".."); err != nil {
				message := fmt.Sprintf("%s", err)
				return checkResult{Passed: false, Message: message}, nil
			}
		}
	}

	if _, err := os.Stat(path.Join(t.GithubDir, t.GithubFile)); os.IsNotExist(err) { // File doesn't exist
		// Creates the file
		if err = createFile(path.Join(t.GithubDir, t.GithubFile)); err != nil {
			message := fmt.Sprintf("%s", err)
			return checkResult{Passed: false, Message: message}, nil
		}
	}

	return checkResult{Passed: true}, nil
}
