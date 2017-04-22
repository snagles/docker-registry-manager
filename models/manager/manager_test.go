package manager

import "testing"

func TestAddRegistry(t *testing.T) {
	err := AddRegistry("http", "localhost", 5000)
	if err != nil {
		t.Error("Did not correctly add a registry!" + err.Error())
	}
}
