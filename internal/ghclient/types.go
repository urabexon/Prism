package ghclient

import "time"

type CheckBucket int

const (
	CheckBucketPass CheckBucket = iota
	CheckBucketFail
	CheckBucketPending
	CheckBucketSkip
	CheckBucketCancel
	checkBucketSentinel
)

func ParseCheckBucket(s string) CheckBucket {
	switch s {
	case "pass":
		return CheckBucketPass
	case "fail":
		return CheckBucketFail
	case "pending":
		return CheckBucketPending
	case "skipping":
		return CheckBucketSkip
	case "cancel":
		return CheckBucketCancel
	default:
		return CheckBucketPending
	}
}

type Check struct {
	Name        string    `json:"name"`
	State       string    `json:"state"`
	Bucket      string    `json:"bucket"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Workflow    string    `json:"workflow"`
	Event       string    `json:"event"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
}

func (c Check) BucketType() CheckBucket {
	return ParseCheckBucket(c.Bucket)
}

func (c Check) Duration() time.Duration {
	if c.CompletedAt.IsZero() || c.StartedAt.IsZero() {
		return 0
	}
	return c.CompletedAt.Sub(c.StartedAt)
}

type CheckSummary struct {
	Total   int
	Pass    int
	Fail    int
	Pending int
	Skip    int
	Cancel  int
}

func (cs CheckSummary) HasChecks() bool {
	return cs.Total > 0
}

func (cs CheckSummary) AllPass() bool {
	return cs.Total > 0 && cs.Pass == cs.Total
}

func (cs CheckSummary) AnyFail() bool {
	return cs.Fail > 0
}

type PR struct {
	Number       int          `json:"number"`
	Title        string       `json:"title"`
	Author       string       `json:"author"`
	State        string       `json:"state"`
	IsDraft      bool         `json:"isDraft"`
	Additions    int          `json:"additions"`
	Deletions    int          `json:"deletions"`
	UpdatedAt    time.Time    `json:"updatedAt"`
	URL          string       `json:"url"`
	HeadRef      string       `json:"headRefName"`
	BaseRef      string       `json:"baseRefName"`
	Labels       []string     `json:"labels"`
	CheckSummary CheckSummary
}

type prJSON struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Author    struct {
		Login string `json:"login"`
	} `json:"author"`
	State       string    `json:"state"`
	IsDraft     bool      `json:"isDraft"`
	Additions   int       `json:"additions"`
	Deletions   int       `json:"deletions"`
	UpdatedAt   time.Time `json:"updatedAt"`
	URL         string    `json:"url"`
	HeadRefName string    `json:"headRefName"`
	BaseRefName string    `json:"baseRefName"`
	Labels      []struct {
		Name string `json:"name"`
	} `json:"labels"`
	StatusCheckRollup []statusCheckJSON `json:"statusCheckRollup"`
}

type statusCheckJSON struct {
	Typename   string `json:"__typename"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	State      string `json:"state"` // for StatusContext typename
}

func computeCheckSummary(checks []statusCheckJSON) CheckSummary {
	var cs CheckSummary
	cs.Total = len(checks)
	for _, c := range checks {
		switch c.Typename {
		case "CheckRun":
			switch c.Conclusion {
			case "SUCCESS":
				cs.Pass++
			case "FAILURE", "TIMED_OUT", "ACTION_REQUIRED":
				cs.Fail++
			case "CANCELLED":
				cs.Cancel++
			case "SKIPPED":
				cs.Skip++
			case "NEUTRAL":
				cs.Pass++
			default:
				// No conclusion yet = pending/in progress
				cs.Pending++
			}
		case "StatusContext":
			switch c.State {
			case "SUCCESS":
				cs.Pass++
			case "FAILURE", "ERROR":
				cs.Fail++
			case "PENDING", "EXPECTED":
				cs.Pending++
			default:
				cs.Pending++
			}
		default:
			cs.Pending++
		}
	}
	return cs
}

func prFromJSON(p prJSON) PR {
	labels := make([]string, len(p.Labels))
	for i, l := range p.Labels {
		labels[i] = l.Name
	}
	return PR{
		Number:       p.Number,
		Title:        p.Title,
		Author:       p.Author.Login,
		State:        p.State,
		IsDraft:      p.IsDraft,
		Additions:    p.Additions,
		Deletions:    p.Deletions,
		UpdatedAt:    p.UpdatedAt,
		URL:          p.URL,
		HeadRef:      p.HeadRefName,
		BaseRef:      p.BaseRefName,
		Labels:       labels,
		CheckSummary: computeCheckSummary(p.StatusCheckRollup),
	}
}

type ReviewComment struct {
	ID              int       `json:"id"`
	Body            string    `json:"body"`
	Path            string    `json:"path"`
	Line            int       `json:"line"`
	StartLine       int       `json:"start_line"`
	Side            string    `json:"side"`
	StartSide       string    `json:"start_side"`
	InReplyToID     int       `json:"in_reply_to_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	User            struct {
		Login string `json:"login"`
	} `json:"user"`
}

type CommentThread struct {
	Root    ReviewComment
	Replies []ReviewComment
}

func GroupCommentThreads(comments []ReviewComment) []CommentThread {
	rootMap := make(map[int]*CommentThread)
	var roots []int

	for _, c := range comments {
		if c.InReplyToID == 0 {
			ct := CommentThread{Root: c}
			rootMap[c.ID] = &ct
			roots = append(roots, c.ID)
		}
	}

	for _, c := range comments {
		if c.InReplyToID != 0 {
			if ct, ok := rootMap[c.InReplyToID]; ok {
				ct.Replies = append(ct.Replies, c)
			}
		}
	}

	threads := make([]CommentThread, 0, len(roots))
	for _, id := range roots {
		threads = append(threads, *rootMap[id])
	}
	return threads
}

func CommentsForFile(threads []CommentThread, path string) []CommentThread {
	var result []CommentThread
	for _, t := range threads {
		if t.Root.Path == path {
			result = append(result, t)
		}
	}
	return result
}

type LineType int

const (
	LineContext LineType = iota
	LineAdded
	LineRemoved
)

type DiffLine struct {
	Type    LineType
	Content string
	OldNum  int
	NewNum  int
}

type Hunk struct {
	Header string
	Lines  []DiffLine
}

type FileDiff struct {
	OldPath  string
	NewPath  string
	Hunks    []Hunk
	IsBinary bool
	IsNew    bool
	IsDelete bool
	IsRename bool
}

type ParsedDiff struct {
	Files []FileDiff
}

type File struct {
	Path      string `json:"path"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

type PRDetail struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	State     string    `json:"state"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	URL       string    `json:"url"`
}

type prDetailJSON struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	State  string `json:"state"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	URL       string    `json:"url"`
}

func prDetailFromJSON(p prDetailJSON) PRDetail {
	return PRDetail{
		Number:    p.Number,
		Title:     p.Title,
		Body:      p.Body,
		State:     p.State,
		Author:    p.Author.Login,
		CreatedAt: p.CreatedAt,
		URL:       p.URL,
	}
}
