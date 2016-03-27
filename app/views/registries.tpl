{{template "base/base.html" .}}

{{template "base/nav.html" .}}

{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li><a href="/">Home</a></li>
        <li><a href="/registries">Registries</a></li>
      </ol>
    </div>

    <div class = "content-block-empty">
      {{range $key, $value := .registries}}
      <div class="col-md-6">
        <div class="col-md-12 white-bg box">
          <div class = "box-header">
            <a href="/registries/{{$value.Name}}/repositories">{{$value.Name}}</span>
          </div>
        </div>
      </div>
      {{end}}
      <div class="col-md-6">
        <div class="col-md-12 white-bg box">
          <div class="row">
            <div class = "box-header">
              <span>Add New</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

{{end}}
