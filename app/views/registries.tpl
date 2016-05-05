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
      <div class="col-lg-12">
        <ul class="boxes">
          <li>
          <a href="/registries/{{$registry.Name}}/repositories">
            <div class="white-bg box col-lg-4 col-md-6 col-sm-12 col-xs-12">
              <div class="col-lg-12">
                <div class="box-container">
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
              </div>
            </div>
          </a>
          </li>
          {{end}}
          <li>
            <div class="well-box box col-lg-4 col-md-6 col-sm-12 col-xs-12">
              <div class="col-lg-12">
                <div class="box-container">
                  <div type= "button" class="add-new" data-toggle="modal" data-target="#new-registry-modal">
                    <i class="fa fa-plus add-new-icon"></i>
                  </div>
                </div>
              </div>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </div>

{{end}}
