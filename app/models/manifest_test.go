package manager

import "testing"

func TestKeyWordTags(t *testing.T) {
	cmd := Command{
		Keywords: []string{"git", "gitlab", "github"},
	}
	tags := cmd.KeywordTags()
	if tags != "git gitlab github" {
		t.Error("Invalid keyword tags returned")
	}
}
