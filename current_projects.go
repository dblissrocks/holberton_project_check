package main

import "encoding/json"

type currentProject struct {
	ProjectName                 string `json:"project_name"`
	ProjectTrackAndBlockDisplay string `json:"project_track_and_block_display"`
	ProjectID                   int    `json:"project_id"`
}

func getCurrentProjects() ([]currentProject, error) {
	resp, err := getWithHolbertonAuth(myCurrentProjectsURL)
	if err != nil {
		return nil, err
	}
	var p []currentProject
	if err := json.Unmarshal([]byte(resp), &p); err != nil {
		return nil, err
	}
	return p, nil
}
