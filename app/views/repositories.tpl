{{template "base/base.html" .}}
{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <nav aria-label="breadcrumb">
        <ol class="breadcrumb">
          <li class="breadcrumb-item"><a href="/">Home</a></li>
          <li class="breadcrumb-item"><a href="/registries">Registries</a></li>
          <li class="breadcrumb-item"><a href="/registries/{{.registryName}}/repositories" class="registry-name">{{.registryName}}</a></li>
          <li class="breadcrumb-item active"><a class="registry-name" aria-current="page">Repositories</a></li>
        </ol>
      </nav>
    </div>
    <div class="row">
      <div class="col-md-12">
        <h1>{{.registryName}}</h1>
      </div>
    </div>
    <div class="content-block white-bg">
      <div class="row">
        <table id="datatable" class="table table-striped" cellspacing="0" width="100%">
          <thead>
            <th>Repository</th>
            <th>Size
              <i class="fa fa-question-circle" aria-hidden="true" data-toggle="tooltip" data-placement="top" title="Compressed tar.gz size"></i>
            </th>
            <th>Tags</th>
          </thead>
          <tbody>
            {{range $key, $repository := .repositories}}
              <tr>
                <td>
                  <a href="/registries/{{$.registryName}}/repositories/{{urlquery $repository.Name}}/tags">{{$repository.Name}}</span></td>
                <td data-order="{{$repository.Size}}">{{bytefmt $repository.Size}}</td>
                <td>{{len $repository.Tags}}</td>
              </tr>
            {{end}}
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <script>
    $(document).ready(function () {
      $('#datatable').DataTable({
        "order": [
          [1, "asc"]
        ],
        "info": false
      });
      $(function () {
        $('[data-toggle="tooltip"]').tooltip({container: 'body'})
      })
    });
  </script>
{{end}}
