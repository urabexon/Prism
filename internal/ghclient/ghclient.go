package ghclient

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func ListPRs() ([]PR, error) {
	cmd := exec.Command("gh", "pr", "list", "--json", "number,title,author,state,isDraft,additions,deletions,updatedAt,url,headRefName,baseRefName,labels,statusCheckRollup")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var prsJSON []prJSON
	if err := json.Unmarshal(out, &prsJSON); err != nil {
		return nil, err
	}
	prs := make([]PR, len(prsJSON))
	for i, p := range prsJSON {
		prs[i] = prFromJSON(p)
	}
	return prs, nil
}

func GetPRDetail(number int) (*PRDetail, error) {
	cmd := exec.Command("gh", "pr", "view", fmt.Sprintf("%d", number), "--json", "number,title,body,state,author,createdAt,url")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var detailJSON prDetailJSON
	if err := json.Unmarshal(out, &detailJSON); err != nil {
		return nil, err
	}
	detail := prDetailFromJSON(detailJSON)
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

func GetDiff(prNumber int) (string, error) {
	cmd := exec.Command("gh", "pr", "diff", fmt.Sprintf("%d", prNumber))
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
