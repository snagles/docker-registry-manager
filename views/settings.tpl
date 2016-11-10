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
        <h1>About</h1>
        <hr>
      </div>
      <div class="row">
        <p>
          <a href="https://travis-ci.org/snagles/docker-registry-manager"><img src="https://travis-ci.org/snagles/docker-registry-manager.svg?branch=master" alt="Build Status" title="" /></a>
          <a href="https://coveralls.io/github/snagles/docker-registry-manager?branch=master"><img src="https://coveralls.io/repos/github/snagles/docker-registry-manager/badge.svg?branch=master" alt="Coverage Status" title="" /></a>
          <a href="https://godoc.org/github.com/snagles/docker-registry-manager"><img src="https://godoc.org/github.com/snagles/docker-registry-manager?status.svg" alt="GoDoc" title="" /></a>
          <a href="https://github.com/snagles/docker-registry-manager/commit/{{.releaseVersion}}"><img src="https://img.shields.io/badge/Release-{{.releaseVersion}}-green.svg" alt="Version" title="" /></a>
        </p>
        <p>Wiki: <small> https://github.com/snagles/docker-registry-manager/tree/master/resources/docs/wiki</small></p>
        <p>Source: <small> https://github.com/snagles/docker-registry-manager</small></p>
      <!--  <p>Feature Requests: <small> https://github.com/snagles/docker-registry-manager</small></p> -->
      <!--  <p>Bug Report: <small> https://github.com/snagles/docs/wiki/bug-report</small></p> -->
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
        <h1>Statistics</h1>
        <hr>
      </div>
      <div class="row">
        <ul class="nav nav-tabs" role="tablist" style="margin-bottom:20px;">
          <li role="presentation" class="active"><a href="#stats" aria-controls="overview" role="tab" data-toggle="tab">Requests</a></li>
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
    <div class="row content-block white-bg" id="logs">
      <div class="row">
        <h1>Logs</h1>
        <hr>
      </div>
      <div class="row">
        <div class="col-lg-12">
          Debug Logging
          <label id="debug">
            <input id = "debug-level" type="checkbox" data-size="small" checked data-onstyle="success"></input>
          </label>
            <button id="archive-logs" type="button" class="btn btn-default"><i class="fa fa-archive"></i> Archive Logs</button>
            <a href="/logs" download="logs.json" class="btn btn-default" download><i class="fa fa-download"></i> Download Logs</a>
            <button id="clear-logs" type="button" class="btn btn-danger"><i class="fa fa-trash"></i> Clear Logs</button>
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


    // Get the stats data table stats
    $.getJSON( "/settings/stats", function( data ) {
      $.each( data, function( i, item ) {
        console.log(item);
        var $tr = $('<tr>').append(
          $("<td>").text(item.method),
          $('<td>').text(item.request_url),
          $('<td data-order='+item.avg_s+'>').text(item.avg_time),
          $('<td data-order='+item.max_s+'>').text(item.max_time),
          $('<td data-order='+item.min_s+'>').text(item.min_time),
          $('<td data-order='+item.total_s+'>').text(item.total_time),
          $('<td>').text(item.times)
        ).appendTo('#stats-datatable');
      });
      $('#stats-datatable').DataTable( {
          "order": [[ 2, "desc" ]],
          "pageLength": 10,
      });
    });

    // Initialize bootstrap toggle for the debug levels
    $('#debug-level').bootstrapToggle();
    // Get the current log level
    var debug = 0;
    $.get("logs/level", function(data) {
      currentLogLevel = data.LogLevel
      if (data.LogLevel == "debug") {
        debug = 1;
        $('#debug-level').bootstrapToggle("on");
      }
      else {
        debug = 0;
        $('#debug-level').bootstrapToggle("off");
      }
    });
    $("#debug").click(function(e){
      $.ajax({type: "GET",
          url: "/logs/toggle-debug",
          success:function(result){
            if (debug == 0) {
              $('#debug-level').bootstrapToggle("on");
              debug = 1;
            }
            else {
              $('#debug-level').bootstrapToggle("off");
              debug = 0;
            }
          }});
    });

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
       ],
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
