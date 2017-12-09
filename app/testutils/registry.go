package testutils

import (
	"bytes"
	"io"
	"testing"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/testutil"
	digest "github.com/opencontainers/go-digest"
)

var TestImages = []string{"foo/aaaa", "foo/bbbb", "foo/cccc"}
var TestTags = []string{"bar-a", "bar-b", "bar-c"}

// SetupRegistry create a new registry with a set url, returns the url and the test env for shutdown
func SetupRegistry(t *testing.T) (string, *TestEnv) {
	env := newTestEnv(t, false)

	baseURL, err := env.builder.BuildBaseURL()
	if err != nil {
		t.Fatalf("unexpected error building catalog url: %v", err)
	}

	// push misc images to the registry
	for _, image := range TestImages {
		for _, tag := range TestTags {
			tagName, _ := reference.WithName(image)
			tagRef, _ := reference.WithTag(tagName, tag)
			manifestURL, err := env.builder.BuildManifestURL(tagRef)
			manifest := &schema2.Manifest{
				Versioned: manifest.Versioned{
					SchemaVersion: 2,
					MediaType:     schema2.MediaTypeManifest,
				},
				Config: distribution.Descriptor{
					Digest:    "sha256:1a9ec845ee94c202b2d5da74a24f0ed2058318bfa9879fa541efaecba272e86b",
					Size:      3253,
					MediaType: schema2.MediaTypeImageConfig,
				},
				Layers: []distribution.Descriptor{
					{
						Digest:    "sha256:463434349086340864309863409683460843608348608934092322395278926a",
						Size:      6323,
						MediaType: schema2.MediaTypeLayer,
					},
					{
						Digest:    "sha256:630923423623623423352523525237238023652897356239852383652aaaaaaa",
						Size:      6863,
						MediaType: schema2.MediaTypeLayer,
					},
				},
			}

			// Push a config, and reference it in the manifest
			sampleConfig := []byte(`{
								"architecture": "amd64",
								"history": [
								  {
								    "created": "2015-10-31T22:22:54.690851953Z",
								    "created_by": "/bin/sh -c #(nop) ADD file:a3bc1e842b69636f9df5256c49c5374fb4eef1e281fe3f282c65fb853ee171c5 in /"
								  },
								  {
								    "created": "2015-10-31T22:22:55.613815829Z",
								    "created_by": "/bin/sh -c #(nop) CMD [\"sh\"]"
								  }
								],
								"rootfs": {
								  "diff_ids": [
								    "sha256:c6f988f4874bb0add23a778f753c65efe992244e148a1d2ec2a8b664fb66bbd1",
								    "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"
								  ],
								  "type": "layers"
								}
							}`)
			sampleConfigDigest := digest.FromBytes(sampleConfig)

			uploadURLBase, _ := startPushLayer(t, env, tagName)
			pushLayer(t, env.builder, tagName, sampleConfigDigest, uploadURLBase, bytes.NewReader(sampleConfig))

			manifest.Config.Digest = sampleConfigDigest
			manifest.Config.Size = int64(len(sampleConfig))

			// Push 2 random layers
			expectedLayers := make(map[digest.Digest]io.ReadSeeker)

			for i := range manifest.Layers {
				rs, dgstStr, err := testutil.CreateRandomTarFile()

				if err != nil {
					t.Fatalf("error creating random layer %d: %v", i, err)
				}
				dgst := digest.Digest(dgstStr)

				expectedLayers[dgst] = rs
				manifest.Layers[i].Digest = dgst

				uploadURLBase, _ := startPushLayer(t, env, tagName)
				pushLayer(t, env.builder, tagName, dgst, uploadURLBase, rs)
			}

			// -------------------
			// Push the manifest with all layers pushed.
			deserializedManifest, err := schema2.FromStruct(*manifest)
			if err != nil {
				t.Fatalf("could not create DeserializedManifest: %v", err)
			}
			_, canonical, err := deserializedManifest.Payload()
			if err != nil {
				t.Fatalf("could not get manifest payload: %v", err)
			}
			dgst := digest.FromBytes(canonical)
			digestRef, _ := reference.WithDigest(tagName, dgst)
			manifestDigestURL, err := env.builder.BuildManifestURL(digestRef)

			putManifest(t, "putting manifest no error", manifestURL, schema2.MediaTypeManifest, manifest)

			// --------------------
			// Push by digest -- should get same result
			putManifest(t, "putting manifest by digest", manifestDigestURL, schema2.MediaTypeManifest, manifest)
		}
	}

	return baseURL, env
}
