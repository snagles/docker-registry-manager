{{template "base/base.html" .}}
{{define "body"}}
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
        <h3>{{.repositoryName}}:{{.tagInfo.Name}}</h3>
        <hr>

      </div>
      <div class="row">
        <ul class="nav nav-tabs" role="tablist">
          <li role="presentation" class="active"><a href="#overview" aria-controls="overview" role="tab" data-toggle="tab">Overview</a></li>
          <li role="presentation"><a href="#stages" aria-controls="stages" role="tab" data-toggle="tab">Dockerfile</a></li>
          <li role="presentation"><a href="#activity" aria-controls="activity" role="tab" data-toggle="tab">Activity</a></li>
          <li role="presentation"><a href="#push" aria-controls="push" role="tab" data-toggle="tab">Push</a></li>
        </ul>
        <div class="tab-content">
          <div role="tabpanel" class="tab-pane active" id="overview">
            <h5>{{.os}}</h5>
            <h5>{{.arch}}</h5>
            <h5>{{.tagInfo.Layers}}</h5>
            <h5>{{.tagInfo.Size}}</h5>
            <h5>{{.tagInfo.TimeAgo}}</h5>
            <h5>{{.tagInfo.UpdatedTime}}</h5>
          </div>
          <div role="tabpanel" class="tab-pane" id="stages">
            <table class="table">
              <thead>
                <th>Intermediate Image ID:</th>
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
          <div role="tabpanel" class="tab-pane" id="activity"></div>
          <div role="tabpanel" class="tab-pane" id="push"></div>
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
