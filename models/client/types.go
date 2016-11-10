package client

import "time"

type list struct {
	Name string
	Tags []string
}

type catalog struct {
	Repositories []string
}

// Image contains all information related to the image
type Image struct {
	Name           string
	Tag            string
	SchemaVersion  int
	Architecture   string
	TagID          uint
	ContainsV1Size bool
	History        []History `json:"history"`
	FsLayers       []struct {
		BlobSum string `json:"blobSum"`
		Size    int64  `json:"-"`
		SizeStr string `json:"-"`
	} `json:"fsLayers"`
}

// History contains the v1 compatibility string and marshaled json
type History struct {
	V1CompatibilityStr string `json:"V1Compatibility"`
	V1Compatibility    V1Compatibility
}

// V1Compatibility contains all information grabbed from the V1Compatibility field from registry v1
type V1Compatibility struct {
	ID              string    `json:"id"`
	IDShort         string    `json:"-"`
	Parent          string    `json:"parent"`
	Created         time.Time `json:"created"`
	Container       string    `json:"container"`
	ContainerConfig struct {
		Hostname        string        `json:"Hostname"`
		Domainname      string        `json:"Domainname"`
		User            string        `json:"User"`
		AttachStdin     bool          `json:"AttachStdin"`
		AttachStdout    bool          `json:"AttachStdout"`
		AttachStderr    bool          `json:"AttachStderr"`
		ExposedPorts    interface{}   `json:"ExposedPorts"`
		PublishService  string        `json:"PublishService"`
		Tty             bool          `json:"Tty"`
		OpenStdin       bool          `json:"OpenStdin"`
		StdinOnce       bool          `json:"StdinOnce"`
		Env             []string      `json:"Env"`
		Cmd             []string      `json:"Cmd"`
		CmdClean        string        `json:"-"`
		Image           string        `json:"Image"`
		Volumes         interface{}   `json:"Volumes"`
		VolumeDriver    string        `json:"VolumeDriver"`
		WorkingDir      string        `json:"WorkingDir"`
		Entrypoint      interface{}   `json:"Entrypoint"`
		NetworkDisabled bool          `json:"NetworkDisabled"`
		MacAddress      string        `json:"MacAddress"`
		OnBuild         []interface{} `json:"OnBuild"`
		Labels          struct {
		} `json:"Labels"`
	} `json:"container_config"`
	DockerVersion string `json:"docker_version"`
	Config        struct {
		Hostname        string        `json:"Hostname"`
		Domainname      string        `json:"Domainname"`
		User            string        `json:"User"`
		AttachStdin     bool          `json:"AttachStdin"`
		AttachStdout    bool          `json:"AttachStdout"`
		AttachStderr    bool          `json:"AttachStderr"`
		ExposedPorts    interface{}   `json:"ExposedPorts"`
		PublishService  string        `json:"PublishService"`
		Tty             bool          `json:"Tty"`
		OpenStdin       bool          `json:"OpenStdin"`
		StdinOnce       bool          `json:"StdinOnce"`
		Env             []string      `json:"Env"`
		Cmd             []string      `json:"Cmd"`
		Image           string        `json:"Image"`
		Volumes         interface{}   `json:"Volumes"`
		VolumeDriver    string        `json:"VolumeDriver"`
		WorkingDir      string        `json:"WorkingDir"`
		Entrypoint      interface{}   `json:"Entrypoint"`
		NetworkDisabled bool          `json:"NetworkDisabled"`
		MacAddress      string        `json:"MacAddress"`
		OnBuild         []interface{} `json:"OnBuild"`
		Labels          struct {
		} `json:"Labels"`
	} `json:"config"`
	Architecture string `json:"architecture"`
	Os           string `json:"os"`
	Size         int    `json:"Size"`
	SizeStr      string `json:"-"`
}
