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
        <form>
          <fieldset class="form-group">
            <label for="log-level">Log Level</label>
            <select class="form-control" stlye ="width:200px" id="log-level">
              <option>1 - Fatal</option>
              <option>2 - Panic</option>
              <option>3 - Error </option>
              <option>4 - Warn</option>
              <option>5 - Info</option>
              <option>6 - Debug</option>
            </select>
          </fieldset>
        </form>
        Empty logs
        Backup logs
        Upload logs
      </div>
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
{{end}}
