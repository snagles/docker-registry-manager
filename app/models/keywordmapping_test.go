package manager

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKeywordMapping(t *testing.T) {
	Convey("Test strings should return the expected keywords", t, func(c C) {
		for _, test := range keywordTestStrings {
			keywords := Keywords(test.TestString)
			c.So(keywords, ShouldNotBeNil)
			c.So(keywords, ShouldResemble, test.Keywords)
		}
	})

}

// TODO: Add more test cases
var keywordTestStrings = []struct {
	TestString string
	Keywords   []string
}{
	{`curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz &&`, []string{"Go"}},
	{`echo "$GOLANG_DOWNLOAD_SHA1 golang.tar.gz" | sha1sum -c - &&`, []string{"Go"}},
	{`./file.zimpl --no-clean 2>&1`, []string{"Zimpl"}},
}
