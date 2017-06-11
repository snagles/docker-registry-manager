package utils

import "regexp"

type KeywordInfo struct {
	Icon   string
	Color  string
	Regexp []string
}

// Keywords takes in any meta information strings and returns the keywordMapping key
// by testing the string against any regex in the KeywordMapping map
func Keywords(s string) []string {
	// use the keyword map to deduplicate keywords for same line
	keywordMap := make(map[string]struct{})
	for key, info := range KeywordMapping {

		// compare the passed string against every regexp possibility
		for _, regex := range info.Regexp {
			// don't check again if the keyword has already been added
			if _, ok := keywordMap[key]; !ok {
				// add i flag to make case insenstive
				r := regexp.MustCompile(`(?i)` + regex)
				if r.Match([]byte(s)) {
					keywordMap[key] = struct{}{}
				}
			}
		}
	}

	// return the slice of keys for the keyword map
	var keywords []string
	for keyword := range keywordMap {
		keywords = append(keywords, keyword)
	}
	return keywords
}

// KeywordMapping keywords, icons, colors, and extensions taken heavily from file-icons atom
// see https://github.com/file-icons/atom for their amazing work
var KeywordMapping = map[string]KeywordInfo{
	"ArtText":       {"arttext-icon", "dark-purple", []string{`.artx`}},
	"Atom":          {"atom-icon", "dark-green", []string{`.atom`}},
	"Bower":         {"bower-icon", "medium-orange", []string{`bower[-_]components`}},
	"Chef":          {"chef-icon", "dark-purple", []string{`.chef`}},
	"CircleCI":      {"circleci-icon", "dark-purple", []string{`.circleci`}},
	"Docker":        {"docker-icon", "dark-blue", []string{"docker"}},
	"Dropbox":       {"dropbox-icon", "medium-blue", []string{"(Dropbox|.dropbox.cache)"}},
	"Emacs":         {"emacs-icon", "medium-purple", []string{".emacs"}},
	"Framework":     {"dylib-icon", "medium-yellow", []string{".framework"}},
	"Git":           {"git-icon", "medium-red", []string{".git"}},
	"Github":        {"github-icon", "medium-orange", []string{".github"}},
	"Gitlab":        {"gitlab-icon", "medium-orange", []string{".gitlab"}},
	"Meteor":        {"meteor-icon", "dark-orange", []string{".meteor"}},
	"Mercurial":     {"hg-icon", "dark-orange", []string{".hg", "mercurial"}},
	"NodeJS":        {"node-icon", "medium-green", []string{"node_modules"}},
	"Packet":        {"package-icon", "medium-green", []string{"(bundle|paket)"}},
	"SVN":           {"svn-icon", "medium-yellow", []string{".svn", "subversion"}},
	"Textmate":      {"textmate-icon", "medium-green", []string{".tmbundle"}},
	"Vagrant":       {"vagrant-icon", "medium-cyan", []string{".vagrant"}},
	"Visual Studio": {"vs-icon", "medium-blue", []string{".vscode"}},
	"Xcode":         {"appstore-icon", "medium-cyan", []string{".xcodeproj"}},
	"Debian":        {"debian-icon", "medium-red", []string{`[a-zA-Z0-9]+.deb\s+`}},
	"Go":            {"go-icon", "medium-blue", []string{`[a-zA-Z0-9]+.go\s+`, `GOPATH`, `GOLANG`, `GOROOT`}},
	"Zimpl":         {"zimpl-icon", "medium-orange", []string{`.(zimpl|zmpl|zpl)`}},
}
