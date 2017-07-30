package manager

import "regexp"

// KeywordInfo contains information used by the front end to build the label
type KeywordInfo struct {
	Icon   string
	Color  string
	Regexp []string
}

// Keywords takes in any meta information strings and returns the keywordMapping key
// by testing the string against any regex in the KeywordMapping map
func Keywords(s string) (keywords []string) {
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
	for keyword := range keywordMap {
		keywords = append(keywords, keyword)
	}
	return keywords
}

// KeywordMapping keywords, icons, colors, and extensions taken heavily from file-icons atom
// see https://github.com/file-icons/atom for their amazing work
var KeywordMapping = map[string]KeywordInfo{
	`ArtText`:       {`arttext-icon`, `dark-purple`, []string{`\.artx`}},
	`Atom`:          {`atom-icon`, `dark-green`, []string{`\.atom`}},
	`Bower`:         {`bower-icon`, `medium-orange`, []string{`bower[-_]components`, `\.(bowerrc|bower\.json|Bowerfile)`}},
	`Chef`:          {`chef-icon`, `dark-purple`, []string{`\.chef`}},
	`CircleCI`:      {`circleci-icon`, `dark-purple`, []string{`\.circleci`, `circle\.yml`}},
	`Docker`:        {`docker-icon`, `dark-blue`, []string{`docker`}},
	`Dropbox`:       {`dropbox-icon`, `medium-blue`, []string{`(Dropbox|\.dropbox\.cache)`}},
	`Emacs`:         {`emacs-icon`, `medium-purple`, []string{`\.emacs`}},
	`Framework`:     {`dylib-icon`, `medium-yellow`, []string{`\.framework`}},
	`Git`:           {`git-icon`, `medium-red`, []string{`\.git`, `git`}},
	`Github`:        {`github-icon`, `medium-orange`, []string{`\.github`}},
	`Gitlab`:        {`gitlab-icon`, `medium-orange`, []string{`\.gitlab`, `\.gitlab-ci\.yml`}},
	`Meteor`:        {`meteor-icon`, `dark-orange`, []string{`\.meteor`}},
	`Mercurial`:     {`hg-icon`, `dark-orange`, []string{`\.hg`, `mercurial`}},
	`Packet`:        {`package-icon`, `medium-green`, []string{`(bundle|paket)`}},
	`SVN`:           {`svn-icon`, `medium-yellow`, []string{`\.svn`, `subversion`}},
	`Textmate`:      {`textmate-icon`, `medium-green`, []string{`\.tmbundle`}},
	`Vagrant`:       {`vagrant-icon`, `medium-cyan`, []string{`\.vagrant`}},
	`Visual Studio`: {`vs-icon`, `medium-blue`, []string{`\.vscode`}},
	`Xcode`:         {`appstore-icon`, `medium-cyan`, []string{`\.xcodeproj`}},
	`Debian`:        {`debian-icon`, `medium-red`, []string{`[a-zA-Z0-9]+\.deb\s+`}},
	`Go`:            {`go-icon`, `medium-blue`, []string{`[a-zA-Z0-9]+\.go\s+`, `GOPATH`, `GOLANG`, `GOROOT`}},
	`Zimpl`:         {`zimpl-icon`, `medium-orange`, []string{`\.(zimpl|zmpl|zpl) `}},
	`Vue`:           {`vue-icon`, `light-green`, []string{`\.vue `}},
	`Ruby`:          {`ruby-icon`, `medium-red`, []string{`\.(rb|ru|ruby|erb|gemspec|god|mspec|pluginspec|podspec|rabl|rake|opal|rails) `}},
	`Rust`:          {`rust-icon`, `medium-maroon`, []string{`\.(rs|\.rlib)`}},
	`R`:             {`r-icon`, `medium-blue`, []string{`\.(r|Rprofile|rsx|rd) `}},
	`Python`:        {`python-icon`, `dark-blue`, []string{`\.(py|\.ipy|pep|py3|\.pyi) `}},
	`Perl`:          {`perl-icon`, `medium-blue`, []string{`\.(pl|perl|pm) `}},
	`Perl6`:         {`perl6-icon`, `medium-purple`, []string{`\.(pl6|p6l|p6m) `}},
	`PHP`:           {`php-icon`, `dark-blue`, []string{`\.php `}},
	`PHPUnit`:       {`phpunit-icon`, `medium-purple`, []string{`\.phpunit\.xml `}},
	`Objective-C`:   {`objc-icon`, `dark-red`, []string{`\.objc`}},
	`NPM`:           {`npm-icon`, `medium-red`, []string{`package\.json|\.npmignore|\.?npmrc|npm-debug\.log|npm-shrinkwrap\.json|package-lock\.json `}},
	`NodeJS`:        {`node-icon`, `medium-green`, []string{`\.(njs|nvmrc|node|node-version) `, `node_modules`}},
	`NGINX`:         {`nginx-icon`, `dark-green`, []string{`nginx\.conf`}},
	`MATLAB`:        {`matlab-icon`, `medium-yellow`, []string{`\.matlab`}},
	`Kotlin`:        {`kotlin-icon`, `dark-blue`, []string{`\.(kt|ktm|kts) `}},
	`Jenkins`:       {`jenkins-icon`, `medium-red`, []string{`Jenkinsfile`}},
	`Java`:          {`java-icon`, `medium-purple`, []string{`\.java `}},
	`Javascript`:    {`js-icon`, `medium-yellow`, []string{`\.(js|_js|jsb|jsm|jss|es6|es|mjs|sjs|ssjs|xsjs|dust) `}},
	`HTML`:          {`html5-icon`, `medium-orange`, []string{`\.html `}},
	`Gulp`:          {`gulp-icon`, `medium-red`, []string{`gulpfile\.js`, `gulpfile\.coffe`, `gulpfile\.babel\.js `}},
	`Erlang`:        {`erlang-icon`, `medium-red`, []string{`\.erl `}},
	`Dart`:          {`dart-icon`, `medium-cyan`, []string{`\.dart `}},
	`CoffeeScript`:  {`coffee-icon`, `medium-maroon`, []string{`\.coffee `}},
	`Clojure`:       {`clojure-icon`, `medium-cyan`, []string{`\.(clj|cl2|cljc|cljx|hic) `}},
	`ClojureScript`: {`cljs-icon`, `dark-cyan`, []string{`\.cljs `}},
	`CMake`:         {`cmake-icon`, `medium-green`, []string{`\.cmake `}},
	`C`:             {`c-icon`, `medium-blue`, []string{`\.(c|h) `}},
	`C++`:           {`cpp-icon`, `light-blue`, []string{`\.(cpp|hpp) `}},
	`C#`:            {`csharp-icon`, `darker-blue`, []string{`\.(csharp|cs) `}},
	`Ansible`:       {`ansible-icon`, `dark-blue`, []string{`\.(ansible|ansible\.yaml|ansible\.yml) `}},
	`Alpine Linux`:  {`alpine-icon`, `dark-blue`, []string{`\.APKBUILD `, ` apk `}},
}
