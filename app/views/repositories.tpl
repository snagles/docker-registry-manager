{{template "base/base.html" .}}
{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li>
          <a href="/">Home</a>
        </li>
        <li>
          <a href="/registries">Registries</a>
        </li>
        <li>
          <a class="registry-name" href="/registries/{{.registryName}}/repositories">{{.registryName}}</a>
        </li>
        <li class="active">Repositories</li>
      </ol>
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
