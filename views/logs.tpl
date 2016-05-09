{{template "base/base.html" .}}
{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li><a href="/">Home</a></li>
        <li class="active">Logs</li>
      </ol>
    </div>
    <div class="content-block white-bg">
      <div class="row">
        <h1>Logs</h1>
        <hr>
      </div>
      <table id="datatable" class="table table-striped table-bordered" cellspacing="0" width="100%">
          <thead>
              <tr>
                  <th>Level</th>
                  <th>Message</th>
                  <th>Entry Time</th>
              </tr>
          </thead>
          <tfoot>
              <tr>
                  <th>Level</th>
                  <th>Message</th>
                  <th>Time</th>
              </tr>
          </tfoot>
          <tbody>
              {{range $key := .logs}}
              <tr>
                  <td>{{$key.Level}}</td>
                  <td>{{$key.Msg}}</td>
                  <td>{{$key.Time}}</td>
              </tr>
              {{end}}
          </tbody>
      </table>
    </div>
  </div>

<script>
$(document).ready(function() {
    $('#datatable').DataTable( {
        "order": [[ 2, "desc" ]]
    });
});
</script>

{{end}}
