{{template "base/base.html" .}}
{{define "body"}}
{{template "new_registry.tpl" .}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li><a href="/">Home</a></li>
        <li><a href="/registries">Registries</a></li>
      </ol>
    </div>

    <div class = "content-block-empty">
      {{range $key, $registry := .registries}}
      <div class="col-lg-4 col-md-6">
        <a href="/registries/{{$registry.Name}}/repositories">
          <div class="col-lg-12 white-bg box">
            <div class="box-header">
              <span>{{$registry.Name}}</span>
            </div>
            <div class="box-body">
              <ul>
                <li>Host: {{$registry.Name}}</li>
                <li>IP: {{$registry.IP}}</li>
                <li>Port: {{$registry.Port}}</li>
                <li>Version: {{$registry.Version}}</li>
              </ul>
            </div>
            <div class="box-footer">
              <span class="label label-success text-capitalize">{{$registry.Status}}</span>
              <span class="label label-info">{{$registry.RepoCount}} Repositories</span>
              <span class="label label-info">{{$registry.RepoTotalSizeStr}}</span>
            </div>
          </div>
        </a>
      </div>
      {{end}}
      <div class="col-lg-4 col-md-6">
        <div class="col-lg-12 well-box box">
          <div type= "button" class="add-new col-lg-8 col-lg-offset-2" data-toggle="modal" data-target="#new-registry-modal">
            <i class="fa fa-plus add-new-icon"></i>
          </div>
        </div>
      </div>
    </div>
  </div>

{{end}}
