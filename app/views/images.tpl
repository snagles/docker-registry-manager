{{template "base/base.html" .}}

{{template "base/nav.html" .}}

{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li><a href="/">Home</a></li>
        <li><a href="/registries">Registries</a></li>
        <li><a href="/registries/{{.registryName}}/repositories">{{.registryName}}</a></li>
        <li><a href="/registries/{{.registryName}}/repositories/{{.repository.EncodedURI}}">{{.repositoryName}}</a></li>
        <li><a href="/registries/{{.registryName}}/repositories/{{.repositoryNameEncode}}/tags/{{.tagName}}">{{.tagName}}</a></li>
        <li class="active">Images</li>
      </ol>
    </div>
    <div class="content-block white-bg">
      <div class="row">
        <h3>{{.repositoryName}}:{{.tagName}}</h3>
        <hr>
      </div>
      <div class="row">
        <table class="table">
          <thead>
            <th>Intermediate Image ID:</th>
            <th>Command:</th>
            <th>Size:</th>
          </thead>
          <tbody>
            {{range $key, $img := .history}}
            <tr>
              <td>{{$img.V1Compatibility.ID}}</td>
              <td>{{$img.V1Compatibility.ContainerConfig.Cmd}}</td>
              <td>{{$img.V1Compatibility.Size}}bytes</td>
            </tr>
            {{end}}
          </tbody>
        </table>
      </div>
    </div>
  </div>

{{end}}
