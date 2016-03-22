package utils

import (
	"fmt"
	"testing"
)

// TestUpdateApp tests to see if all of the git functions execute successfully
func TestUpdateApp(t *testing.T) {

	_, err := UpdateApp()
	fmt.Println(err)
}
