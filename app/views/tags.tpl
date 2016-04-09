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
        <li class="active">Tags</li>
      </ol>
    </div>
    <div class="content-block white-bg">
      <div class="row">
        <h1>{{.repositoryName}}</h1>
        <hr>
      </div>
      <div class="row">
        <table class="table">
          <thead>
            <th>Tags:</th>
            <th>Created:</th>
            <th>Size:</th>
            <th>Layers:</th>
            <th>Delete:</th>
          </thead>
          <tbody>
            {{range $key, $tag := .tags}}
            <tr>
              <td><a href=/registries/{{$.registryName}}/repositories/{{$.repositoryName}}/tags/{{$tag.Name}}/images>{{$tag.Name}}</span></td>
              <td>{{$tag.TimeAgo}}</td>
              <td>{{$tag.Size}}</td>
              <td>{{$tag.Layers}}</td>
            </tr>
            {{end}}
          </tbody>
        </table>
      </div>
    </div>
  </div>

{{end}}
