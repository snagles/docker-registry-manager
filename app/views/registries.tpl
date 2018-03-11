{{template "base/base.html" .}}
{{define "body"}}
  {{template "newregistry.tpl" .}}
  <div class="right-content-container">
    <div class="header">
      <nav aria-label="breadcrumb">
        <ol class="breadcrumb">
          <li class="breadcrumb-item">
            <a href="/">Home</a>
          </li>
          <li class="breadcrumb-item active" aria-current="page">Registries</li>
        </ol>
      </nav>
    </div>

    <div class="content-block-empty">
      <div class="card-deck">
        {{range $key, $registry := .registries}}
          <div class="card">
            <img class="card-img-top" src="http://1000logos.net/wp-content/uploads/2017/07/Docker-Logo-500x148.png" alt="Card image cap">
            <div class="card-body">
              <h5 class="card-title">{{$registry.Name}}</h5>
              <h6 class="card-subtitle mb-2 text-muted">{{$registry.IP}}</h6>
              <div class="text-right">
                <a href="#" class="btn btn-info">Edit</a>
                <a href="/registries/{{$registry.Name}}/repositories" class="btn btn-orange">View</a>
              </div>
            </div>
            <div class="card-footer">
              <small class="text-muted">Last updated {{timeAgo $registry.LastRefresh}}</small>
            </div>
          </div>
       {{end}}
        <div class="card d-flex bg-light justify-content-center align-text-middle" style="min-height:275px">
          <div type="button" class="add-new align-self-center" data-toggle="modal" data-target="#new-registry-modal">
            <i class="fa fa-plus add-new-icon"></i>
          </div>
        </div>
      </div>
    </div>
  </div>
{{end}}
