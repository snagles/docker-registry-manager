{{template "base/base.html" .}}
{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li><a href="/">Home</a></li>
        <li><a href="/registries">Registries</a></li>
        <li class="active">All</li>
      </ol>
    </div>
    <div class="content-block white-bg">
      <div class="row">
        <h1>All Repositories</h1>
        <hr>
      </div>
      <div class="row">
        <table id="datatable" class="table table-striped table-bordered" cellspacing="0" width="100%">
          <thead>
            <th>Repository</th>
            <th>Registry</th>
            <th>Tags</th>
          </thead>
          <tfoot>
            <th>Repository Name:</th>
            <th>Registry:</th>
            <th>Tags</th>
          </tfoot>
          <tbody>
            {{range $key, $repository := .repositories}}
              {{range $key, $repo := $repository}}
              <tr>
                <td><a href=/registries/{{$repo.Registry}}/repositories/{{$repo.EncodedURI}}/tags>{{$repo.Name}}</span></td>
                <td>{{$repo.Registry}}</td>
                <td>{{$repo.TagCount}}</td>
              </tr>
              {{end}}
            {{end}}
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <script>
  $(document).ready(function() {
      $('#datatable').DataTable( {
          "order": [[ 1, "asc" ]],
          "pageLength": 25
      } );
      $(function () {
        $('[data-toggle="tooltip"]').tooltip({container: 'body'})
      })
  });
  </script>
{{end}}
