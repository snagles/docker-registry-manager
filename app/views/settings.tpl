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
    <div class="row content-block white-bg">
      <div class="row">
        <h1>Logs</h1>
        <hr>
      </div>
      <div class="row">
        <div class="col-lg-12">
          <form class="form ">
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

        </form>
      </div>
    </div>
    <table id="datatable" class="table table-striped table-bordered" cellspacing="0" width="100%">
        <thead>
            <tr>
                <th data-toggle="tooltip" data-placement="top" title="Message type ranging from 'Fatal' (most severe) to 'Debug' (Used for debugging problems)">Level</th>
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
    <button type="button" class="btn btn-default"><i class="fa fa-archive"></i> Backup Logs</button>
    <button type="button" class="btn btn-default"><i class="fa fa-download"></i> Download Logs</button>
    <button type="button" class="btn btn-danger"><i class="fa fa-trash"></i> Clear Logs</button>
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
      $('#datatable').DataTable( {
          "order": [[ 2, "desc" ]],
          "pageLength": 10
      } );
  });
</script>
{{end}}
