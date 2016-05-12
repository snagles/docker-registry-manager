{{template "base/base.html" .}}
{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li><a href="/">Home</a></li>
        <li class="active">Settings</li>
      </ol>
    </div>
    <div class="row content-block white-bg">
      <div class="row">
        <h1>General</h1>
        <hr>
      </div>
      <div class="row">
      </div>
    </div>
    <div class="row content-block white-bg">
      <div class="row">
        <h1>Upgrade</h1>
        <hr>
      </div>
      <div class="row">
      </div>
    </div>
    <div class="row content-block white-bg">
      <div class="row">
        <h1>Tasks</h1>
        <hr>
      </div>
      <div class="row">
      </div>
    </div>
    <div class="row content-block white-bg" id="logs">
      <div class="row">
        <h1>Logs</h1>
        <hr>
      </div>
      <div class="row">
        <div class="col-lg-12">
          <form>
            <fieldset class="form-group">
              <label for="log-level">Log Level</label>
              <select class="form-control" id="log-level">
                <option>1 - Fatal</option>
                <option>2 - Panic</option>
                <option>3 - Error </option>
                <option>4 - Warn</option>
                <option>5 - Info</option>
                <option>6 - Debug</option>
              </select>
              <div class="text-muted">Higher numbers include itself, and all log level numbers below. </div>
            </fieldset>
            <fieldset class="form-group">
              <label for="log-level">Log Retention</label>
              <label class="radio-inline">
                <input type="radio" name="retention-days" id="7-days" value="7"> 7 Days
              </label>
              <label class="radio-inline">
                <input type="radio" name="retention-days" id="14-days" value="14-days"> 14 Days
              </label>
              <label class="radio-inline">
                <input type="radio" name="retention-days" id="31-days" value="31-days"> 31 Days
              </label>
              <label class="radio-inline">
                <input type="radio" name="retention-days" id="forever-days" value="forever"> Forever
              </label>
            </fieldset>
            <button id="archive-logs" type="button" class="btn btn-default"><i class="fa fa-archive"></i> Archive Logs</button>
            <a href="/logs" download="logs.json" class="btn btn-default" download><i class="fa fa-download"></i> Download Logs</a>
            <button id="clear-logs" type="button" class="btn btn-danger"><i class="fa fa-trash"></i> Clear Logs</button>
        </form>
      </div>
    </div>
    <hr>
    <table id="datatable" class="table table-striped table-bordered" cellspacing="0" width="100%">
        <thead>
            <tr>
                <th>Level</th>
                <th>Message</th>
                <th>Time</th>
            </tr>
        </thead>
        <tfoot>
            <tr>
                <th>Level</th>
                <th>Message</th>
                <th>Time</th>
            </tr>
        </tfoot>
    </table>
  </div>
    <div class="row content-block white-bg">
      <div class="row">
        <h1>About</h1>
        <hr>
      </div>
      <div class="row">
        <p>Version: <small> 1.06</small></p>
        <p>Wiki: <small> https://github.com/stefannaglee/docker-registry-manager/tree/master/resources/docs/wiki</small></p>
        <p>Chat: <small> https://chat.stefannaglee.com/</small></p>
        <p>Source: <small> https://github.com/stefannaglee/docker-registry-manager</small></p>
        <p>GoDocs: <small> https://stefannaglee/docker-registry-manager/godocs</small></p>
        <p>Feature Requests: <small> https://github.com/stefannaglee/docker-registry-manager</small></p>
        <p>Bug Report: <small> https://github.com/stefannaglee/docs/wiki/bug-report</small></p>
      </div>
    </div>

    <!--
    <div class="row content-block white-bg">
      <div class="row">
        <h1>Activity</h1>
      </div>
      <div class="row">
        Desktop notifications
        Desktop notifications
      </div>
    </div>
    -->
  </div>
<script>
  $(document).ready(function() {
      var table = $('#datatable').DataTable( {
          "ajax": {
              url: '/logs',
              dataSrc: '',
          },
          "order": [[ 2, "desc" ]],
          "pageLength": 10,
          "columns": [
            { "data": "level" },
            { "data": "msg"},
            { "data": "time"}
         ]
      } );

    $("#clear-logs").click(function(e){
      e.preventDefault();
      $.ajax({type: "POST",
          url: "/logs/clear",
          success:function(result){
            table.ajax.reload();
            $("#logs").append("<div class='alert alert-success'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a> <strong>Success!</strong> Cleared logs. </div>");
            $(".alert").alert();
            window.setTimeout(function() { $(".alert").alert('close'); }, 5000);
          }});
    });
    $("#archive-logs").click(function(e){
      e.preventDefault();
      $.ajax({type: "POST",
          url: "/logs/archive",
          success:function(result){
            table.ajax.reload();
            $("#logs").append("<div class='alert alert-success'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a> <strong>Success!</strong> Archived logs in /logs. </div>");
            $(".alert").alert();
            window.setTimeout(function() { $(".alert").alert('close'); }, 5000);
          }});
    });
  });
</script>
{{end}}
