package main

import (
	"encoding/json"
	"fmt"
)

type project struct {
	Name  string `json:"name"`
	Tasks []struct {
		ID         int    `json:"id"`
		Title      string `json:"title"`
		GithubRepo string `json:"github_repo"`
		GithubDir  string `json:"github_dir"`
		GithubFile string `json:"github_file"`
	} `json:"tasks"`
}

func main() {
	body, err := getWithHolbertonAuth("https://intranet.hbtn.io/projects/97.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var p project
	if err := json.Unmarshal([]byte(body), &p); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("************")
	fmt.Println(p.Name)

	for _, task := range p.Tasks {
		fmt.Println("************")
		fmt.Println(task.Title)
		fmt.Println(task.GithubRepo)
		fmt.Println(task.GithubDir)
		fmt.Println(task.GithubFile)
	}
}
