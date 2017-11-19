package manager

import (
	"testing"
)

func TestKeywordMapping(t *testing.T) {
	for _, test := range keywordTestStrings {
		keywords := Keywords(test.TestString)
		if keywords == nil {
			t.Errorf("Keywords method returned no results, when expected %s", test.Keyword)
		}
		for _, key := range keywords {
			if key != test.Keyword {
				t.Errorf("Incorrect expected keyword: %s", key)
			}
		}
	}

}

// TODO: Add more test cases
var keywordTestStrings = []struct {
	TestString string
	Keyword    string
}{
	{`curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz &&`, "Go"},
	{`echo "$GOLANG_DOWNLOAD_SHA1 golang.tar.gz" | sha1sum -c - &&`, "Go"},
	{`./file.zimpl --no-clean 2>&1`, "Zimpl"},
}
