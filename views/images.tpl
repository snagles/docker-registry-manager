{{template "base/base.html" .}}
{{define "body"}}
{{template "modal.tpl" .}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li><a href="/">Home</a></li>
        <li><a href="/registries">Registries</a></li>
        <li><a class="registry-name" href="/registries/{{.registryName}}/repositories">{{.registryName}}</a></li>
        <li><a href="/registries/{{.registryName}}/repositories">Repositories</a></li>
        <li><a class="registry-name" href="/registries/{{.registryName}}/repositories/{{.repositoryNameEncode}}/tags">{{.repositoryName}}</a></li>
        <li><a class="registry-name" href="/registries/{{.registryName}}/repositories/{{.repositoryNameEncode}}/tags">Tags</a></li>
        <li class="active">{{.tagName}}</li>
      </ol>
    </div>
    <div class="content-block white-bg">
      <div class="row">
        <ul class="nav nav-tabs" role="tablist">
          <li role="presentation" class="active"><a href="#overview" aria-controls="overview" role="tab" data-toggle="tab">Overview</a></li>
          <li role="presentation"><a href="#stages" aria-controls="stages" role="tab" data-toggle="tab">Stages</a></li>
          <li role="presentation"><a href="#layers" aria-controls="layers" role="tab" data-toggle="tab">Layers</a></li>
          <li role="presentation"><a href="#deploy" aria-controls="deploy" role="tab" data-toggle="tab">Deploy</a></li>
        </ul>
        <div class="tab-content">
          <div role="tabpanel" class="tab-pane active" id="overview">
            <div class="row">
              <div class="col-md-12">
                <div class="col-md-12">
                  <h4>Image Overview</h4>
                  <ul>
                    <li>Operating System: {{.os}}</li>
                    <li>Layers: {{len .image.FsLayers}}</li>
                    <li>Size: {{.image.Size}}</li>
                    <li>Language: </li>
                    <li>Last Updated: {{.tagInfo.TimeAgo}}</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
          <div role="tabpanel" class="tab-pane" id="stages">
            <table class="table">
              <thead>
                <th>Image ID:</th>
                <th>Command:</th>
                {{if $.containsV1Size}}
                <th>Size:</th>
                {{end}}
              </thead>
              <tbody>
                {{range $key, $img := .history}}
                <tr>
                  <td>{{$img.V1Compatibility.IDShort}}</td>
                  <td>{{$img.V1Compatibility.ContainerConfig.CmdClean}}</td>
                  {{if $.containsV1Size}}
                  <td>{{$img.V1Compatibility.SizeStr}}</td>
                  {{end}}
                </tr>
                {{end}}
              </tbody>
            </table>
          </div>
          <div role="tabpanel" class="tab-pane" id="layers">
            <table class="table">
              <thead>
                <th>ID:</th>
                <th>Digest:</th>
                <th>Size:</th>
                <th>Blob:</th>
              </thead>
              <tbody>
                {{range $index, $layer := .layers}}
                <tr>
                  <td>{{$index}}</td>
                  <td>{{$layer.BlobSum}}</td>
                  <td>{{$layer.SizeStr}}</td>
                  <td><a href="{{$.registry.Scheme}}://{{$.registry.Name}}:{{$.registry.Port}}/{{$.registry.Version}}/{{$.repositoryName}}/blobs/{{$layer.BlobSum}}" download class="btn btn-sm btn-success"><span class="glyphicon glyphicon-download-alt"></span> Download</a></td>
                </tr>
                {{end}}
              </tbody>
            </table>
          </div>
          <div role="tabpanel" class="tab-pane" id="deploy">
            <div>Push to {{.tagInfo.Name}}:</div>
            <ol>
              <li><code>cd $PROJECTNAME</code></li>
              <li><code>docker build --rm -t {{.registry.Name}}:{{.registry.Port}}/{{.repositoryName}}:{{.tagInfo.Name}} .</code></li>
              <li><code>docker push {{.registry.Name}}:{{.registry.Port}}/{{.repositoryName}} .</code></li>
            </ol>
          </div>
        </div>
      </div>
    </div>

  </div>

  <script>
  $('#image-tabs li a').click(function (e) {
    e.preventDefault()
    $(this).tab('show')
  })

  </script>

{{end}}
