package manager

import (
	manifestV1 "github.com/docker/distribution/manifest/schema1"
)

type Tag struct {
	*manifestV1.SignedManifest
	ID   string
	Name string
}
