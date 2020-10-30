package gitee

// Commit represents a GitHub commit.
type Commit struct {
	SHA  *string `json:"sha,omitempty"`
	Date *string `json:"date,omitempty"`
}
