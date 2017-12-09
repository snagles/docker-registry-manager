{{template "base/base.html" .}}
{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li>
          <a href="/">Home</a>
        </li>
        <li class="active">Events</li>
      </ol>
    </div>
    <div class="row">
      <div class="col-md-12">
        <h1>Events</h1>
        <div></div>
        <div class="content-block white-bg col-md-12" id="logs">
          <table id="events-table" class="table table-striped table-bordered" cellspacing="0" width="100%">
            <thead>
              <tr>
                <th>Registry</th>
                <th>Action</th>
                <th>Repository</th>
                <th>Tag</th>
                <th>Size</th>
                <th>URL</th>
                <th>Source</th>
              </tr>
            </thead>
            <tfoot>
              <tr>
                <th>Registry</th>
                <th>Action</th>
                <th>Repository</th>
                <th>Tag</th>
                <th>Size</th>
                <th>URL</th>
                <th>Source</th>
              </tr>
            </tfoot>
            <tbody>
              {{range $registry, $ids := .events}}
                {{range $id, $event := $ids}}
                  <tr>
                    <td>{{$registry}}</td>
                    <td>{{$event.Action}}</td>
                    <td>{{$event.Target.Repository}}</td>
                    <td>{{$event.Target.Tag}}</td>
                    <td>{{bytefmt $event.Target.Size}}</td>
                    <td>{{$event.Target.URL}}</td>
                    <td>{{$event.Source.Addr}}</td>
                  </tr>
                {{end}}
              {{end}}
            </tbody>
          </table>
        </div>
        <script>
          $('#events-table').DataTable({"pageLength": 50});
        </script>
      {{ end }}
