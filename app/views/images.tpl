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
        <li><a class="registry-name" href="/registries/{{.registryName}}/repositories/{{.repositoryNameEncode}}">{{.repositoryName}}</a></li>
        <li><a class="registry-name" href="/registries/{{.registryName}}/repositories/{{.repositoryNameEncode}}">Tags</a></li>
        <li><a href="/registries/{{.registryName}}/repositories/{{.repositoryNameEncode}}/tags/{{.tagInfo.Name}}">{{.tagInfo.Name}}</a></li>
        <li class="active">Images</li>
      </ol>
    </div>
    <div class="content-block white-bg">
      <div class="row">
        <ul class="nav nav-tabs" role="tablist">
          <li role="presentation" class="active"><a href="#overview" aria-controls="overview" role="tab" data-toggle="tab">Overview</a></li>
          <li role="presentation"><a href="#stages" aria-controls="stages" role="tab" data-toggle="tab">Dockerfile</a></li>
          <li role="presentation"><a href="#layers" aria-controls="layers" role="tab" data-toggle="tab">Layers</a></li>
          <li role="presentation"><a href="#private-registry" aria-controls="private-registry" role="tab" data-toggle="tab">Private Registry</a></li>
          <li role="presentation"><a href="#dockerhub" aria-controls="dockerhub" role="tab" data-toggle="tab">Dockerhub</a></li>
        </ul>
        <div class="tab-content">
          <div role="tabpanel" class="tab-pane active" id="overview">
            <div class="row">
              <div class="col-md-4">
                <div class="col-md-12">
                  <h4>Image</h4>
                  <ul>
                    <li>Operating System: {{.os}}</li>
                    <li>Architecture: {{.arch}}</li>
                    <li>Layers: {{.tagInfo.Layers}}</li>
                    <li>Size: {{.tagInfo.Size}}</li>
                  </ul>
                </div>
              </div>
              <div class="col-md-4">
                <div class="col-md-12">
                  <h4>Registry Host</h4>
                  <ul>
                    <li>Name: {{.registry.Name}}</li>
                    <li>IP: {{.registry.IP}}</li>
                    <li>Port: {{.registry.Port}}</li>
                    <li>Version: {{.registry.Version}}</li>
                  </ul>
                </div>
              </div>
              <div class="col-md-4">
                <div class="col-md-12">
                  <h4>Metadata</h4>
                  <ul>
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
          <div role="tabpanel" class="tab-pane" id="private-registry">
            <div>Push to {{.tagInfo.Name}}:</div>
            <ol>
              <li><code>cd $PROJECTNAME</code></li>
              <li><code>docker build --rm -t {{.registry.Name}}:{{.registry.Port}}/{{.repositoryName}}:{{.tagInfo.Name}} .</code></li>
              <li><code>docker push {{.registry.Name}}:{{.registry.Port}}/{{.repositoryName}} .</code></li>
            </ol>

            <div>Download {{.tagInfo.Name}}:</div>
            <div>Push to another private registry {{.tagInfo.Name}}:</div>
            <a href="{{$.registry.Scheme}}://{{$.registry.Name}}:{{$.registry.Port}}/{{$.registry.Version}}/{{$.repositoryName}}/{{$.tagInfo.Name}}" download><i class="fa fa-download"></i></a>
          </div>
          <div role="tabpanel" class="tab-pane" id="dockerhub">
            <h3>Push to Dockerhub</h3>
            <form>
              <div class="form-group">
                <label for="account-name">Docker Account Name</label>
                <input type="account-name" class="form-control" id="account-name" placeholder="Docker Account Name">
              </div>
              <div class="form-group">
                <label for="username">Username</label>
                <input type="password" class="form-control" id="username" placeholder="Username">
              </div>
              <div class="form-group">
                <label for="password">Password</label>
                <input type="password" class="form-control" id="password" placeholder="Password">
              </div>
              <button type="push" data-toggle="modal" data-target="#dockerhub-modal" class="btn btn-default">Push</button>
            </form>
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
