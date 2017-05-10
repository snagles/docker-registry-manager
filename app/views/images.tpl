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
        </ul>
        <div class="tab-content">
          <div role="tabpanel" class="tab-pane active" id="overview">
            <div class="row">
              <div class="col-md-12">
                <div class="col-md-12">
                  <h4>Image Overview</h4>
                  <ul>
                    <li>Layers: {{.tag.LayerCount}}</li>
                    <li>Last Updated: {{.tag.LastModifiedTimeAgo}}</li>
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
              </thead>
              <tbody>
                {{range $key, $img := .tag.Histories}}
                <tr>
                  <td>{{$img.IDShort}}</td>
                  <td>{{$img.ContainerConfig.CmdClean}}</td>
                </tr>
                {{end}}
              </tbody>
            </table>
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
