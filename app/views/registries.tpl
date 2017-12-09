{{template "base/base.html" .}}
{{define "body"}}
  {{template "newregistry.tpl" .}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li>
          <a href="/">Home</a>
        </li>
        <li>
          <a href="/registries">Registries</a>
        </li>
      </ol>
    </div>

    <div class="content-block-empty">
      <div class="col-lg-12">
        <ul class="boxes">
          {{range $key, $registry := .registries}}
            <li data-registry="{{$registry.Name}}">
              <a href="/registries/{{$registry.Name}}/repositories">
                <div class="white-bg box col-lg-4 col-md-6 col-sm-12 col-xs-12">
                  <div class="col-lg-12">
                    <div class="box-container">
                      <div class="box-header">
                        <h2>{{$registry.Name}}</h2>
                      </div>
                      <div class="box-body col-md-12 border-between">
                        <div class="col-md-4 metric">
                          <h3 class="metric-value">{{len $registry.Repositories}}</h3>
                          <small>{{ $repoCount := len $registry.Repositories }}
                            {{ if eq $repoCount 1 }}
                              Repository
                            {{else}}
                              Repositories
                            {{ end }}
                          </small>
                        </div>
                        <div class="col-md-4 metric">
                          <h3 class="metric-value">{{$registry.TagCount}}</h3>
                          <small>{{ $tagCount := $registry.TagCount }}
                            {{ if eq $tagCount 1 }}
                              Tag
                            {{else}}
                              Tags
                            {{ end }}
                          </small>
                        </div>
                        <div class="col-md-4 metric">
                          <h3 class="metric-value">{{$registry.LayerCount}}</h4>
                          <small>{{ $layerCount := $registry.LayerCount }}
                            {{ if eq $layerCount 1 }}
                              Layer
                            {{else}}
                              Layers
                            {{ end }}
                          </small>
                        </div>
                      </div>
                      <div class="box-footer">
                        {{if eq $registry.Status "UP" }}
                          <span class="label label-success text-capitalize">{{$registry.Status}}</span>
                        {{else}}
                          <span class="label label-danger text-capitalize">{{$registry.Status}}</span>
                        {{ end }}
                        {{if ne $registry.IP "" }}
                          <span class="label label-info">{{$registry.IP}}</span>
                        {{ end }}
                        <span class="label label-info text-capitalize">{{$registry.Version}}</span>
                        <span class="label label-info text-uppercase">{{$registry.Scheme}}</span>
                        {{if ne $registry.Pushes 0 }}
                          <span class="label label-info">{{$registry.Pushes}} Pushes</span>
                        {{ end }}
                        {{if ne $registry.Pulls 0 }}
                          <span class="label label-info">{{$registry.Pulls}} Pulls</span>
                        {{ end }}
                        {{if ne $registry.TTL 0 }}
                          <span class="label label-info">Refresh: {{$registry.TTL}}</span>
                        {{ end }}
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
                  <div type="button" class="add-new" data-toggle="modal" data-target="#new-registry-modal">
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
