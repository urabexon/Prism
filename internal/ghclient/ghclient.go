package ghclient

import (
	"encoding/json"
	"os/exec"
)

type PR struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func ListPRs() ([]PR, error) {
	cmd := exec.Command("gh", "pr", "list", "--json", "number,title,author")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var prs []PR
	if err := json.Unmarshal(out, &prs); err != nil {
		return nil, err
	}
	return prs, nil
}
