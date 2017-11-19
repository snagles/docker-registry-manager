package registry

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/astaxie/beego"
	manifestV2 "github.com/docker/distribution/manifest/schema2"
	"github.com/snagles/docker-registry-manager/app/models"
)

type ImagesController struct {
	beego.Controller
}

// GetImages returns the template for the images page
func (c *ImagesController) GetImages() {

	registryName := c.Ctx.Input.Param(":registryName")
	repositoryName, _ := url.QueryUnescape(c.Ctx.Input.Param(":splat"))
	repositoryNameEncode := url.QueryEscape(repositoryName)
	c.Data["tagName"] = c.Ctx.Input.Param(":tagName")

	registry, _ := manager.AllRegistries.Registries[registryName]
	c.Data["registry"] = registry

	tag, _ := registry.Repositories[repositoryName].Tags[c.Ctx.Input.Param(":tagName")]
	c.Data["tag"] = tag

	labels := make(map[string]manager.KeywordInfo)
	for _, h := range tag.History {
		// run each command through the keyword lookup
		for _, cmd := range h.Commands {
			for _, keyword := range cmd.Keywords {
				labels[keyword] = manager.KeywordMapping[keyword]
			}
		}
	}
	c.Data["labels"] = labels

	// build the js chart dataset
	type segmentInfo struct {
		Stage    int      `json:"stage"`
		Cmd      string   `json:"cmd"`
		Keywords []string `json:"keywords"`
		Size     string
	}

	type dataset struct {
		Label            string   `json:"label"`
		Data             []int64  `json:"data"`
		BackgroundColor  []string `json:"backgroundColor"`
		BorderColor      []string `json:"borderColor"`
		BorderWidth      int64    `json:"borderWidth"`
		CutoutPercentage int64    `json:"cutoutPercentage"`

		// Custom data fields
		Info []segmentInfo `json:"info"`
	}

	var chart []dataset
	colors := []string{"#43A19E", "#7B43A1", "#F2317A", "#FF9824", "#58CF6C"}
	for i, history := range tag.History {
		ds := dataset{}
		for _, cmd := range history.Commands {
			ds.Data = append(ds.Data, 10)
			color := colors[0]
			colors = append(colors[:0], colors[1:]...)
			ds.BackgroundColor = append(ds.BackgroundColor, color)
			// add stage tooltip info
			ds.Info = append(ds.Info, segmentInfo{Stage: i + 1, Keywords: cmd.Keywords, Cmd: cmd.Cmd})
			colors = append(colors, color)
			if len(history.Commands) != 1 {
				ds.BorderColor = []string{"#FFF"}
				ds.BorderWidth = 1
			}
		}
		chart = append(chart, ds)
	}

	for i := len(chart)/2 - 1; i >= 0; i-- {
		opp := len(chart) - 1 - i
		chart[i], chart[opp] = chart[opp], chart[i]
	}

	c.Data["chart"] = chart
	c.Data["registryName"] = registryName
	c.Data["repositoryNameEncode"] = repositoryNameEncode
	c.Data["repositoryName"] = repositoryName

	// Compare the two manifest layers
	c.TplName = "images.tpl"
	hubManifest, err := manager.HubGetManifest(repositoryName, tag.Name)
	if hubManifest == nil || err != nil {
		c.Data["dockerHub"] = struct {
			Error    error
			ImageURL string
		}{
			errors.New("Unable to retrieve information from dockerhub"),
			"",
		}
		return
	} else if hubManifest.SchemaVersion == tag.SchemaVersion {
		diffLayers := make(map[string]struct{})
		var size int64
		if err == nil {
			for _, layer := range tag.DeserializedManifest.Layers {
				diffLayers[layer.Digest.String()] = struct{}{}
			}
			for _, layer := range hubManifest.Layers {
				size += layer.Size
				if _, ok := diffLayers[layer.Digest.String()]; ok {
					delete(diffLayers, layer.Digest.String())
				}
			}
		}

		c.Data["dockerHub"] = struct {
			DiffLayers map[string]struct{}
			Manifest   *manifestV2.DeserializedManifest
			ImageURL   string
			Error      error
			Size       int64
		}{
			diffLayers,
			hubManifest,
			fmt.Sprintf("https://hub.docker.com/r/library/%s/tags/%s/", repositoryName, tag.Name),
			err,
			size,
		}
	} else {
		c.Data["dockerHub"] = struct {
			Error    error
			ImageURL string
		}{
			errors.New("Different manifest scheme versions"),
			fmt.Sprintf("https://hub.docker.com/r/library/%s/tags/%s/", repositoryName, tag.Name),
		}
	}
}
