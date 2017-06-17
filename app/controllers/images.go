package controllers

import (
	"fmt"
	"net/url"

	"github.com/astaxie/beego"
	manifestV2 "github.com/docker/distribution/manifest/schema2"
	"github.com/snagles/docker-registry-manager/app/models/manager"
	"github.com/snagles/docker-registry-manager/utils"
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

	labels := make(map[string]utils.KeywordInfo)
	for _, h := range tag.History {
		// run each command through the keyword lookup
		for _, cmd := range h.Commands {
			for _, keyword := range cmd.Keywords {
				labels[keyword] = utils.KeywordMapping[keyword]
			}
		}
	}
	c.Data["labels"] = labels

	// build the js chart dataset

	type dataset struct {
		Label            string   `json:"label"`
		Data             []int64  `json:"data"`
		BackgroundColor  []string `json:"backgroundColor"`
		BorderColor      []string `json:"borderColor"`
		BorderWidth      int64    `json:"borderWidth"`
		CutoutPercentage int64    `json:"cutoutPercentage"`
		TestProperty     string   `json:"testProperty"`
	}
	var chart []dataset
	colors := []string{"#43A19E", "#7B43A1", "#F2317A", "#FF9824", "#58CF6C"}
	for _, history := range tag.History {
		ds := dataset{}
		for _, cmd := range history.Commands {
			ds.Data = append(ds.Data, 10)
			color := colors[0]
			colors = append(colors[:0], colors[1:]...)
			ds.BackgroundColor = append(ds.BackgroundColor, color)
			colors = append(colors, color)
			ds.Label = cmd.Cmd
			if len(history.Commands) != 1 {
				ds.BorderColor = []string{"#FFF"}
				ds.BorderWidth = 1
			}
		}
		chart = append(chart, ds)
	}

	c.Data["chart"] = chart
	c.Data["registryName"] = registryName
	c.Data["repositoryNameEncode"] = repositoryNameEncode
	c.Data["repositoryName"] = repositoryName

	// Compare the two manifest layers
	hubManifest, err := manager.HubGetManifest(repositoryName, tag.Name)
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

	// Index template
	c.TplName = "images.tpl"
}
