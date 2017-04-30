package git

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
)

// ReleaseVersion contains the git tag version
var ReleaseVersion string

func init() {
	sha, err := GetAppSHA()
	if err != nil {
		logrus.Error(err.Error())
		logrus.Error("Could not get git sha version")
	}
	ReleaseVersion = sha
}

// Ref contains the response from the github ref request
type Ref struct {
	Object struct {
		Sha  string `json:"sha"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"object"`
	Ref string `json:"ref"`
	URL string `json:"url"`
}

/* Github oriented functions
 *
 *
 *
 */

// GetRemoteBranchSHA returns the commit SHA of the remote branch from github
func GetRemoteBranchSHA(branch string) (string, error) {

	// Get the branches current SHA
	url := "https://api.github.com/repos/snagles/docker-registry-manager/git/refs/heads/" + branch
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	r := Ref{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return string(body), err
	}

	return r.Object.Sha, nil
}

/* Local git oriented functions
 *
 *
 *
 */

// GetAppBranch returns the current branch name
func GetAppBranch() (string, error) {
	cmdArgs := []string{"rev-parse", "--abbrev-ref", "HEAD"}
	branchName, err := exec.Command("/usr/bin/git", cmdArgs...).Output()
	return string(branchName), err
}

// GetAppSHA returns the current SHA location
func GetAppSHA() (string, error) {
	cmdArgs := []string{"rev-parse", "--short", "HEAD"}
	branchSHA, err := exec.Command("/usr/bin/git", cmdArgs...).Output()
	clean := strings.Replace(string(branchSHA), "\n", "", -1)
	return string(clean), err
}

// GetBaseSHA returns the current base SHA location
// http://stackoverflow.com/questions/3258243/check-if-pull-needed-in-git
func GetBaseSHA() (string, error) {
	cmdArgs := []string{"merge-base", "@", "@{u}"}
	baseSHA, err := exec.Command("/usr/bin/git", cmdArgs...).Output()
	return string(baseSHA), err
}

// GetTag returns the current tag
// https://git-scm.com/book/en/v2/Git-Basics-Tagging
func GetTag() (string, error) {
	cmdArgs := []string{"describe", "--abbrev=0"} // only shows annotated
	//cmdArgs := []string{"describe"} // This shows all tags
	tag, err := exec.Command("/usr/bin/git", cmdArgs...).Output()
	return string(tag), err
}

// RemoteUpdate updates with origin
func RemoteUpdate() error {
	cmdArgs := []string{"remote", "update"}
	err := exec.Command("/usr/bin/git", cmdArgs...).Run()
	return err
}

// IsAppUpToDate checks to see if the local status of the git tree is up to date with the remote
func IsAppUpToDate() (bool, error) {

	// Fetch origin refs
	RemoteUpdate()

	// Get the local branch info
	localBranchName, err := GetAppBranch()
	if err != nil {
		return true, err
	}
	localBranchSHA, err := GetAppSHA()
	if err != nil {
		return true, err
	}

	// Get the remote info
	remoteBaseSHA, err := GetBaseSHA()
	if err != nil {
		return true, err
	}
	remoteBranchSHA, err := GetRemoteBranchSHA(localBranchName)
	if err != nil {
		return true, err
	}

	// Compare the local SHA and remote SHA, if they're the same we are up to date
	// http://stackoverflow.com/questions/3258243/check-if-pull-needed-in-git
	if localBranchSHA == remoteBranchSHA {
		return true, nil
	} else if localBranchSHA == remoteBaseSHA {
		// This means we need to update
		return false, nil
	}

	// If branch is diverged or you need to push
	return true, err
}
