package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

type task struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	GithubRepo string `json:"github_repo"`
	GithubDir  string `json:"github_dir"`
	GithubFile string `json:"github_file"`
}

type project struct {
	Name  string `json:"name"`
	Tasks []task
}

type checkResult struct {
	Passed  bool
	Message string // If the task didn't pass
}

// Checks that one of the lines passed matches the compiled regexp
// (Useful to check if got remote well-configured)
func oneOfTheLinesMatches(lines []string, r *regexp.Regexp) bool {
	for _, line := range lines {
		match := r.MatchString(line)
		if match {
			return true
		}
	}
	return false
}

// Checks that a given task is as expected
func checkTaskFiles(t task) (checkResult, error) {

	// Checking about the expected Git repository name
	out, err := exec.Command("git", "remote", "-v").Output()
	if err != nil {
		return checkResult{}, err
	}
	lines := strings.Split(string(out), "\n")
	pattern := fmt.Sprintf("git@github.com:.+/%s.git", t.GithubRepo)
	r, err := regexp.Compile(pattern)
	if err != nil {
		return checkResult{}, err
	}
	if !oneOfTheLinesMatches(lines, r) {
		message := fmt.Sprintf("No SSH remote points to a \"%s\" Github repository; make sure to run this application from your project's directory.", t.GithubRepo)
		return checkResult{Passed: false, Message: message}, nil
	}

	if _, err := os.Stat(t.GithubDir); os.IsNotExist(err) { // Directory doesn't exist
		message := fmt.Sprintf("No \"%s\" directory was found in this directory.", t.GithubDir)
		return checkResult{Passed: false, Message: message}, nil
	}

	if _, err := os.Stat(path.Join(t.GithubDir, t.GithubFile)); os.IsNotExist(err) { // Directory doesn't exist
		message := fmt.Sprintf("No \"%s\" file was found in the \"%s\" directory.", t.GithubFile, t.GithubDir)
		return checkResult{Passed: false, Message: message}, nil
	}

	// TODO checking the dirs and files

	return checkResult{Passed: true}, nil
}

func main() {

	// Check that .git exists
	if _, err := os.Stat(".git"); os.IsNotExist(err) { // File doesn't exist
		fmt.Println("Current directory does not seem to be a git repository; make sure to run this application from your project's directory.")
		return
	}

	body, err := getWithHolbertonAuth("https://intranet.hbtn.io/projects/97.json")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var p project
	if err := json.Unmarshal([]byte(body), &p); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("************************")
	fmt.Printf("Name of the project: %s\n", p.Name)

	for _, t := range p.Tasks {
		res, err := checkTaskFiles(t)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if res.Passed {
			fmt.Printf("[OK] Task \"%s\"\n", t.Title)
		} else {
			fmt.Printf("[FAILED] Task \"%s\" (%s)\n", t.Title, res.Message)
		}
	}
}
