package ghclient

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Author struct {
	Login string `json:"login"`
}

type PR struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Author Author `json:"author"`
}

type File struct {
	Path      string `json:"path"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
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

type PRDetail struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	State     string `json:"state"`
	Author    Author `json:"author"`
	CreatedAt string `json:"createdAt"`
	URL       string `json:"url"`
}

func GetPRDetail(number int) (*PRDetail, error) {
	cmd := exec.Command("gh", "pr", "view", fmt.Sprintf("%d", number), "--json", "number,title,body,state,author,createdAt,url")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var detail PRDetail
	if err := json.Unmarshal(out, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

func ListFiles(prNumber int) ([]File, error) {
	cmd := exec.Command("gh", "pr", "view", fmt.Sprintf("%d", prNumber),
		"--json", "files", "--jq", ".files")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var files []File
	if err := json.Unmarshal(out, &files); err != nil {
		return nil, err
	}
	return files, nil
}
