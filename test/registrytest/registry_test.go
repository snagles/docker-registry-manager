package registrytest

import "testing"

func TestNewRegistry(t *testing.T) {
	r := New("https.yml")
	if r == nil {
		t.Fatalf("Failed to create new test registry")
	}

	if err := Start(r); err != nil {
		t.Fatalf("Unable to start registry: %s", err.Error())
	}

	if err := Seed(r); err != nil {
		t.Fatalf("Unable to seed registry: %s", err.Error())
	}

	if err := Stop(r); err != nil {
		t.Fatalf("Unable to stop registry: %s", err.Error())
	}
}
