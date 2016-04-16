{{template "base/base.html" .}}

{{template "base/nav.html" .}}

{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li><a href="/">Home</a></li>
        <li><a href="/registries">Registries</a></li>
        <li class="active">{{.registryName}}</li>
      </ol>
    </div>
    <div class="content-block white-bg">
      <div class="row">
        <h1>{{.registryName}}</h1>
        <hr>
      </div>
      <div class="row">
        <table class="table">
          <thead>
            <th>ID:</th>
            <th>Repository Name:</th>
            <th>Tags</th>
          </thead>
          <tbody>
            {{range $key, $repository := .repositories}}
            <tr>
              <td>{{$key}}</td>
              <td><a href=/registries/{{$.registryName}}/repositories/{{$repository.EncodedURI}}/tags>{{$repository.Name}}</span></td>
              <td>{{$repository.TagCount}}</td>
            </tr>
            {{end}}
          </tbody>
        </table>
      </div>
    </div>
  </div>

{{end}}
