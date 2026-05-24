package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Store struct {
	path string
	data storeData
}

type storeData struct {
	Repos map[string]*RepoState `json:"repos"`
}

type RepoState struct {
	PRs map[int]*PRState `json:"prs"`
}

type PRState struct {
	Read          bool            `json:"read"`
	ReviewedFiles map[string]bool `json:"reviewedFiles"`
	LastSeenAt    time.Time       `json:"lastSeenAt"`
}

func NewStore() (*Store, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("config dir: %w", err)
	}
	dir := filepath.Join(configDir, "ghpr-tui")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("creating state dir: %w", err)
	}
	return NewWithPath(filepath.Join(dir, "state.json"))
}

func (s *Store) repoState(repo string) *RepoState {
	rs, ok := s.data.Repos[repo]
	if !ok {
		rs = &RepoState{PRs: make(map[int]*PRState)}
		s.data.Repos[repo] = rs
	}
	if rs.PRs == nil {
		rs.PRs = make(map[int]*PRState)
	}
	return rs
}

func (s *Store) prState(repo string, number int) *PRState {
	rs := s.repoState(repo)
	ps, ok := rs.PRs[number]
	if !ok {
		ps = &PRState{ReviewedFiles: make(map[string]bool)}
		rs.PRs[number] = ps
	}
	if ps.ReviewedFiles == nil {
		ps.ReviewedFiles = make(map[string]bool)
	}
	return ps
}

func (s *Store) IsRead(repo string, number int) bool {
	return s.prState(repo, number).Read
}

func (s *Store) MarkRead(repo string, number int) {
	ps := s.prState(repo, number)
	ps.Read = true
	ps.LastSeenAt = time.Now()
}

func (s *Store) ToggleRead(repo string, number int) {
	ps := s.prState(repo, number)
	ps.Read = !ps.Read
}

func (s *Store) IsFileReviewed(repo string, number int, path string) bool {
	return s.prState(repo, number).ReviewedFiles[path]
}

func (s *Store) ToggleFileReviewed(repo string, number int, path string) {
	ps := s.prState(repo, number)
	ps.ReviewedFiles[path] = !ps.ReviewedFiles[path]
}

func (s *Store) MarkFileReviewed(repo string, number int, path string) {
	ps := s.prState(repo, number)
	ps.ReviewedFiles[path] = true
}

func (s *Store) MarkAllReviewed(repo string, number int, paths []string) {
	ps := s.prState(repo, number)
	for _, p := range paths {
		ps.ReviewedFiles[p] = true
	}
}

// ReviewedFileCount returns the number of reviewed files for a PR.
func (s *Store) ReviewedFileCount(repo string, number int) int {
	ps := s.prState(repo, number)
	count := 0
	for _, v := range ps.ReviewedFiles {
		if v {
			count++
		}
	}
	return count
}

func NewWithPath(path string) (*Store, error) {
	s := &Store{
		path: path,
		data: storeData{Repos: make(map[string]*RepoState)},
	}
	raw, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(raw, &s.data)
	}
	if s.data.Repos == nil {
		s.data.Repos = make(map[string]*RepoState)
	}
	return s, nil
}

func (s *Store) Save() error {
	raw, _ := json.MarshalIndent(s.data, "", "  ")
	return os.WriteFile(s.path, raw, 0o644)
}
