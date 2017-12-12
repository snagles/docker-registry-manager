{{template "base/base.html" .}}
{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li>
          <a href="/">Home</a>
        </li>
        <li class="active">Logs</li>
      </ol>
    </div>
    <div class="row">
      <div class="col-md-12">
        <h1>Logs</h1>
      </div>
    </div>
    <div class="content-block white-bg" id="logs">
      <table id="log-table" class="table table-striped table-bordered" cellspacing="0" width="100%">
        <thead>
          <tr>
            <th>Level</th>
            <th>Message</th>
            <th>File</th>
            <th>Line</th>
            <th>Source</th>
            <th>Time</th>
          </tr>
        </thead>
        <tfoot>
          <tr>
            <th>Level</th>
            <th>Message</th>
            <th>File</th>
            <th>Line</th>
            <th>Source</th>
            <th>Time</th>
          </tr>
        </tfoot>
        <tbody>
          {{range $i, $entry := .logs}}
            <tr>
              <td>{{$entry.Level}}</td>
              <td>{{$entry.Message}}</td>
              <td>{{$entry.File}}</td>
              <td>{{$entry.Line}}</td>
              <td>{{$entry.Source}}</td>
              <td>{{$entry.Time}}</td>
            </tr>
          {{end}}
        </tbody>
      </table>
    </div>
    <div class="row">
      <div class="col-md-12">
        <h1>Requests</h1>
      </div>
    </div>
    <div class="content-block white-bg">
      <div class="row">
        <ul class="nav nav-tabs" role="tablist" style="margin-bottom:20px;">
          <li role="presentation" class="active">
            <a href="#stats" aria-controls="overview" role="tab" data-toggle="tab">Requests</a>
          </li>
        </ul>
        <div role="tabpanel" class="tab-pane" id="stats">
          <table id="stats-datatable" class="table table-striped table-bordered" cellspacing="0" width="100%">
            <thead>
              <tr>
                <th>Method</th>
                <th>Request URL</th>
                <th>Average Time</th>
                <th>Max Time</th>
                <th>Min Time</th>
                <th>Total Time</th>
                <th>Request Count</th>
              </tr>
            </thead>
            <tbody>
              {{range $i, $stat := .stats}}
                <tr>
                  <td>{{index $stat "method"}}</td>
                  <td>{{index $stat "request_url"}}</td>
                  <td data-order="{{statToSeconds (index $stat `avg_time`)}}">{{index $stat "avg_time"}}</td>
                  <td data-order="{{statToSeconds (index $stat `max_time`)}}">{{index $stat "max_time"}}</td>
                  <td data-order="{{statToSeconds (index $stat `min_time`)}}">{{index $stat "min_time"}}</td>
                  <td data-order="{{statToSeconds (index $stat `total_time`)}}">{{index $stat "total_time"}}</td>
                  <td>{{index $stat "times"}}</td>
                </tr>
              {{end}}
            </tbody>
            <tfoot>
              <tr>
                <th>Method</th>
                <th>Request URL</th>
                <th>Average Time</th>
                <th>Max Time</th>
                <th>Min Time</th>
                <th>Total Time</th>
                <th>Request Count</th>
              </tr>
            </tfoot>
          </table>
        </div>
      </div>
    </div>
  </div>

  <script>
    $(document).ready(function () {
      // use local storage for alerting success/failures after page loads
      if (localStorage.getItem("alert")) {
        $(".right-content-container").prepend(localStorage.getItem("alert"));
        $(".alert").alert();
        window.setTimeout(function () {
          $(".alert").alert('close');
        }, 5000);
        localStorage.clear();
      }
      // Get the stats data table stats
      $('#stats-datatable').DataTable({
        "order": [
          [2, "desc"]
        ],
        "pageLength": 10
      });

      var table = $('#log-table').DataTable({
        "order": [
          [2, "desc"]
        ],
        "pageLength": 10,
        dom: "<'row'<'col-sm-3'l><'col-sm-6'B><'col-sm-3'f>><'row'<'col-sm-12'tr>><'row'<'col-sm-5'i><'col-sm-7'p>>",
        buttons: [
          {
            text: '<button class="btn btn-default dropdown-toggle" type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false"> <i class="fa fa-edit"></i> Log Level: <span id="active-level">{{.activeLevel}} </span><span class="caret"> </span></but' +
                'ton> <ul class="dropdown-menu" aria-labelledby="dropdownMenu1"> {{ range $i, $level :=.allLevels }} <li><a class="dropdown-item" data-level="{{$level}}" href="#">{{$level}}</a></li> {{ end }} </ul>'
          }, {
            text: '<button id="archive-logs" type="button" class="btn btn-default btn-group" style="margin-right:5px"><i class="fa fa-archive"></i> Archive Logs</button>',
            action: function (e, dt, node, config) {
              e.preventDefault();
              $.ajax({
                type: "POST",
                url: "/logs/archive",
                success: function (response) {
                  if (response.Error == null) {
                    window.location.reload(true)
                    $('html, body').animate({
                      scrollTop: 0
                    }, 0);
                    localStorage.setItem("alert", "<div class='alert alert-success'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a> <strong>Success!</strong> Archived logs. </div>")
                    location.reload()
                  }
                }
              });
            }
          }, {
            text: '<button href="/logs/json" download="logs.json" class="btn btn-default btn-group" style="margin-right:5px" download><i class="fa fa-download"></i> Download Logs</button>'
          }, {
            text: '<button id="clear-logs" type="button" class="btn btn-danger btn-group"><i class="fa fa-trash"></i> Clear Logs</button>',
            action: function (e, dt, node, config) {
              e.preventDefault();
              $.ajax({
                type: "DELETE",
                url: "/logs",
                success: function (response) {
                  if (response.Error == null) {
                    $('html, body').animate({
                      scrollTop: 0
                    }, 0);
                    localStorage.setItem("alert", "<div class='alert alert-success'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a> <strong>Success!</strong> Cleared logs. </div>")
                    location.reload()
                  }
                }
              });
            }
          }
        ]
      });
      $('.dropdown-item').click(function (e) {
        e.preventDefault();
        var level = $(this).attr("data-level");
        $.ajax({
          type: "POST",
          url: "/logs/level/" + level,
          success: function (response) {
            if (response.Error == null) {
              window.location.reload(true)
              $('html, body').animate({
                scrollTop: 0
              }, 0);
              localStorage.setItem("alert", "<div class='alert alert-success'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a> <strong>Success!</strong> Updated log level to {{.activeLevel}}</div>")
              location.reload()
            }
          }
        });
      });
    });
  </script>
{{ end }}
