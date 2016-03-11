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

func checkRepo(t task) (checkResult, error) {
	// No github repo specified = always ok
	if t.GithubRepo == "" {
		return checkResult{Passed: true}, nil
	}

	// Checking about the expected Git repository name
	out, err := exec.Command("git", "remote", "-v").Output()
	if err != nil {
		return checkResult{}, err
	}
	lines := strings.Split(string(out), "\n")

	sshPattern := fmt.Sprintf(sshGithubFormat, t.GithubRepo)
	sshCompiledPattern, err := regexp.Compile(sshPattern)
	if err != nil {
		return checkResult{}, err
	}
	httpsPattern := fmt.Sprintf(httpsGithubFormat, t.GithubRepo)
	httpsCompiledPattern, err := regexp.Compile(httpsPattern)
	if err != nil {
		return checkResult{}, err
	}

	if !oneOfTheLinesMatches(lines, sshCompiledPattern) && !oneOfTheLinesMatches(lines, httpsCompiledPattern) {
		message := fmt.Sprintf("No remote pointing to a \"%s\" Github repository; make sure to run this application from your project's directory.", t.GithubRepo)
		return checkResult{Passed: false, Message: message}, nil
	}
	return checkResult{Passed: true}, nil
}

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

	// TODO checking the dirs and files

	return checkResult{Passed: true}, nil
}

func main() {
	// Looks for '-check' flag to specify checking files
	check := false
	for _, v := range os.Args {
		if v == "-check" {
			check = true
		}
	}

	// Check that .git exists
	if _, err := os.Stat(".git"); os.IsNotExist(err) { // File doesn't exist
		fmt.Println("Current directory does not seem to be a git repository; make sure to run this application from your project's directory.")
		return
	}

	// Listing current projects to choose one from
	ps, err := getCurrentProjects()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("************************")
	fmt.Println("Your current projects:")
	for _, p := range ps {
		fmt.Printf(" %d. %s  (%s)\n", p.ProjectID, p.ProjectName, p.ProjectTrackAndBlockDisplay)
	}

	fmt.Print("\nWhich project number? ")
	var projectNumber int
	if _, err := fmt.Scanf("%d", &projectNumber); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Project choser, let's proceed to get its information
	url := fmt.Sprintf(projectURLFormat, projectNumber)
	body, err := getWithHolbertonAuth(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var p project
	if err := json.Unmarshal([]byte(body), &p); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Now, let's check what's up, and display our findings
	fmt.Println("************************")
	fmt.Printf("Name of the project: %s\n", p.Name)

	for _, t := range p.Tasks {
		res, err := checkRepo(t)
		if check && res.Passed {
			res, err = checkTaskFiles(t)
		} else if res.Passed {
			res, err = createTaskFiles(t)
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if res.Passed {
			if check == true {
				fmt.Printf("[OK] Checked \"%s\"\n", t.Title)
			} else {
				fmt.Printf("[OK] Created \"%s\"\n", t.Title)
			}
		} else {
			fmt.Printf("[FAILED] Task \"%s\" (%s)\n", t.Title, res.Message)
		}
	}
}
