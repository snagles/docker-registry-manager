{{template "base/base.html" .}}
{{define "body"}}
  {{template "newregistry.tpl" .}}
  <div class="right-content-container">
    <div class="header">
      <nav aria-label="breadcrumb">
        <ol class="breadcrumb">
          <li class="breadcrumb-item"><a href="/">Home</a></li>
          <li class="breadcrumb-item active" aria-current="page">Registries</li>
        </ol>
      </nav>
    </div>

    <div class="content-block-empty">
      <div class="col-lg-12">
        <ul class="boxes">
          {{range $key, $registry := .registries}}
            <li data-registry="{{$registry.Name}}">
              <a href="/registries/{{$registry.Name}}/repositories">
                <div class="white-bg box col-lg-4 col-md-6 col-sm-12 col-xs-12">
                  <div class="col-lg-12">
                    <div>
                      <div class="box-header row">
                        <h2>{{$registry.Name}}</h2>
                      </div>
                      <div class="box-body row border-between">
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
                      <div class="box-footer row">
                        {{if eq $registry.Status "UP" }}
                          <span class="badge badge-success text-capitalize">{{$registry.Status}}</span>
                        {{else}}
                          <span class="badge badge-danger text-capitalize">{{$registry.Status}}</span>
                        {{ end }}
                        {{if ne $registry.IP "" }}
                          <span class="badge badge-info">{{$registry.IP}}</span>
                        {{ end }}
                        <span class="badge badge-info text-uppercase">{{$registry.Scheme}}</span>
                        {{if ne $registry.Pushes 0 }}
                          <span class="badge badge-info">{{$registry.Pushes}} Pushes</span>
                        {{ end }}
                        {{if ne $registry.Pulls 0 }}
                          <span class="badge badge-info">{{$registry.Pulls}} Pulls</span>
                        {{ end }}
                        {{if ne $registry.TTL 0 }}
                          <span class="badge badge-info">Refresh Rate: {{$registry.TTL}}</span>
                        {{ end }}
                          <span class="badge badge-info">Refreshed: {{timeAgo $registry.LastRefresh}}</span>
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
