package manager

import (
	"strings"
	"time"

	"github.com/docker/distribution"
)

// V1Compatibility contains meta information for when each layer contained
// its own configuration for each stage
type V1Compatibility struct {
	Architecture string `json:"architecture"`
	Config       struct {
		Hostname     string      `json:"Hostname"`
		Domainname   string      `json:"Domainname"`
		User         string      `json:"User"`
		AttachStdin  bool        `json:"AttachStdin"`
		AttachStdout bool        `json:"AttachStdout"`
		AttachStderr bool        `json:"AttachStderr"`
		Tty          bool        `json:"Tty"`
		OpenStdin    bool        `json:"OpenStdin"`
		StdinOnce    bool        `json:"StdinOnce"`
		Env          []string    `json:"Env"`
		Cmd          []string    `json:"Cmd"`
		ArgsEscaped  bool        `json:"ArgsEscaped"`
		Image        string      `json:"Image"`
		Volumes      interface{} `json:"Volumes"`
		WorkingDir   string      `json:"WorkingDir"`
		Entrypoint   interface{} `json:"Entrypoint"`
		OnBuild      interface{} `json:"OnBuild"`
		Labels       struct {
		} `json:"Labels"`
	} `json:"config"`
	Container       string `json:"container"`
	ContainerConfig struct {
		Hostname     string      `json:"Hostname"`
		Domainname   string      `json:"Domainname"`
		User         string      `json:"User"`
		AttachStdin  bool        `json:"AttachStdin"`
		AttachStdout bool        `json:"AttachStdout"`
		AttachStderr bool        `json:"AttachStderr"`
		Tty          bool        `json:"Tty"`
		OpenStdin    bool        `json:"OpenStdin"`
		StdinOnce    bool        `json:"StdinOnce"`
		Env          []string    `json:"Env"`
		Cmd          []string    `json:"Cmd"`
		ArgsEscaped  bool        `json:"ArgsEscaped"`
		Image        string      `json:"Image"`
		Volumes      interface{} `json:"Volumes"`
		WorkingDir   string      `json:"WorkingDir"`
		Entrypoint   interface{} `json:"Entrypoint"`
		OnBuild      interface{} `json:"OnBuild"`
		Labels       struct {
		} `json:"Labels"`
	} `json:"container_config"`
	Created       time.Time `json:"created"`
	DockerVersion string    `json:"docker_version"`
	History       []struct {
		Created       time.Time                `json:"created"`
		Author        string                   `json:"author,omitempty"`
		CreatedBy     string                   `json:"created_by,omitempty"`
		Comment       string                   `json:"comment,omitempty"`
		EmptyLayer    bool                     `json:"empty_layer,omitempty"`
		ManifestLayer *distribution.Descriptor `json:"manifest_layer"`
		ShellType     string
		Commands      []Command
	} `json:"history"`
	Os     string `json:"os"`
	Rootfs struct {
		Type      string   `json:"type"`
		DiffIDs   []string `json:"diff_ids,omitempty"`
		BaseLayer string   `json:"base_layer,omitempty"`
	} `json:"rootfs"`
}

// Command represents any command in the Dockerfile, using the commands keywords are parsed from the extensions
type Command struct {
	Cmd      string
	Keywords []string
}

// KeywordTags returns space delimited list of keywords for the given tag
func (c *Command) KeywordTags() string {
	return strings.Join(c.Keywords, " ")
}
