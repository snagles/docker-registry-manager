package utils

import "fmt"

func Update() error {

	isUpToDate, err := IsAppUpToDate()
	fmt.Println(isUpToDate)
	fmt.Println(err)
	return nil
}

// IsAppUpToDate checks to see if the local status of the git tree is up to date with the remote
func IsAppUpToDate() (bool, error) {

	// Fetch origin refs
	RemoteUpdate()

	// Get the local branch info
	localBranchName, err := GetAppBranch()
	if err != nil {
		Log.Error("Could not get app branch name from the local git copy")
		return true, err
	}
	localBranchSHA, err := GetAppSHA()
	if err != nil {
		Log.Error("Could not get app branch SHA from the local git copy")
		return true, err
	}

	// Get the remote info
	remoteBaseSHA, err := GetBaseSHA()
	if err != nil {
		Log.Error("Could not get app branch SHA from the local git copy")
		return true, err
	}
	remoteBranchSHA, err := GetRemoteBranchSHA(localBranchName)
	if err != nil {
		Log.Error("Could not get remote branch SHA for " + localBranchName)
		Log.Error(err)
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
	Log.Error("You have edited things locally with commits. Cannot autoresolve.")
	return true, err
}
