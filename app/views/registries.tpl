{{template "base/base.html" .}}
{{define "body"}}
  {{template "newregistry.tpl" .}}
  {{range $key, $registry := .registries}}
    {{template "editregistry.tpl" $registry}}
  {{end}}
  <div class="right-content-container">
    <div class="header">
      <nav aria-label="breadcrumb">
        <ol class="breadcrumb">
          <li class="breadcrumb-item">
            <a href="/">Home</a>
          </li>
          <li class="breadcrumb-item active" aria-current="page">Registries</li>
        </ol>
      </nav>
    </div>

    <div class="content-block-empty">
      <div class="card-deck">
        {{range $key, $registry := .registries}}
          <div class="card col-lg-5 pl-0 pr-0 gutterless">
            <div class="card-header bg-light">
              <h3 class="card-title">
                {{$registry.Name}}
                <small class="text-muted">{{$registry.IP}}</small>
              </h3>
            </div>
            <div class="card-body">
              <div class="container">
                <div class="row">
                  <div class="col-lg-4 col-md-6 text-center border-right">
                    <div class="row-fluid no-gutters">
                      <div class="col" style="position: relative;">
                        <canvas id="{{$registry.Name}}-repositories-chart"></canvas>
                      </div>
                    </div>
                    <div class="row-fluid no-gutters">
                      <span class="metric-value">{{len $registry.Repositories}}</span>
                    </div>
                    <h8>Repos</h8>
                  </div>
                  <div class="col-lg-4 col-md-6 text-center border-right">
                    <div class="row-fluid no-gutters">
                      <div class="col" style="position: relative;">
                        <canvas id="{{$registry.Name}}-tags-chart"></canvas>
                      </div>
                    </div>
                    <div class="row-fluid no-gutters">
                      <span class="metric-value">{{$registry.TagCount}}</span>
                    </div>
                    <h8>Tags</h8>
                  </div>
                  <div class="col-lg-4 col-md-6 text-center">
                    <div class="row-fluid no-gutters">
                      <div class="col" style="position: relative;">
                        <canvas id="{{$registry.Name}}-layers-chart"></canvas>
                      </div>
                    </div>
                    <div class="row-fluid no-gutters">
                      <span class="metric-value">{{$registry.LayerCount}}</span>
                    </div>
                    <h8>Layers</h8>
                  </div>
                </div>
                <div class="row mt-4">
                </div>
              </div>
            </div>
            <div class="card-footer">
              <div class="row d-flex align-items-center">
                <div class="col-5">
                  <small class="text-muted align-baseline">Last updated {{timeAgo $registry.LastRefresh}}</small>
                </div>
                <div class="col-7 text-right">
                  <a href="#" data-toggle="modal" data-target="#edit-registry-modal-{{$registry.Name}}" class="btn btn-info">Edit</a>
                  <a href="/registries/{{$registry.Name}}/repositories" class="btn btn-orange">View</a>
                </div>
              </div>
            </div>
          </div>
       {{end}}
        <div class="col-lg-5 card d-flex bg-light justify-content-center align-text-middle" style="min-height:275px">
          <div type="button" class="add-new align-self-center" data-toggle="modal" data-target="#new-registry-modal">
            <i class="fa fa-plus add-new-icon"></i>
          </div>
        </div>
      </div>
    </div>
  </div>

  <script>
  function padChart(times, values) {
    if (times.length == 1) {
      var data = {
          labels: [times[0], times[0]],
          datasets: [{ data: [values[0], values[0]] }]
      }
      return data
    } else {
      var data = {
          labels: times,
          datasets: [{ data: values }]
      }
      return data
    }
  }
  {{range $key, $registry := .registries}}
    const repos = document.getElementById('{{$registry.Name}}-repositories-chart').getContext('2d');
    const chart1 = new Chart(repos, {
      type: 'line',
      data: padChart({{$registry.HistoryTimes}}, {{$registry.HistoryRepos}}),
      options: {
        responsive: true,
        legend: { display: false },
        elements: {
          line: {
            borderColor: '#000000',
            backgroundColor: '#FFF',
            borderWidth: 1
          },
          point: {
            radius: 0
          }
        },
        tooltips: {
          enabled: true
        },
        scales: {
          yAxes: [ { display: false } ],
          xAxes: [ {
            display: false
          }]
        }
      }
    });

    const tags = document.getElementById('{{$registry.Name}}-tags-chart').getContext('2d');
    const chart2 = new Chart(tags, {
      type: 'line',
      data: padChart({{$registry.HistoryTimes}}, {{$registry.HistoryTags}}),
      options: {
        responsive: true,
        legend: { display: false },
        elements: {
          line: {
            borderColor: '#000000',
            backgroundColor: '#FFF',
            borderWidth: 1
          },
          point: {
            radius: 0
          }
        },
        tooltips: {
          enabled: true
        },
        scales: {
          yAxes: [ { display: false } ],
          xAxes: [ { display: false } ]
        }
      }
    });

    const layers = document.getElementById('{{$registry.Name}}-layers-chart').getContext('2d');
    const chart3 = new Chart(layers, {
      type: 'line',
      data: padChart({{$registry.HistoryTimes}}, {{$registry.HistoryLayers}}),
      options: {
        responsive: true,
        legend: { display: false },
        elements: {
          line: {
            borderColor: '#000000',
            backgroundColor: '#FFF',
            borderWidth: 1
          },
          point: {
            radius: 0
          }
        },
        tooltips: {
          enabled: true
        },
        scales: {
          yAxes: [ { display: false } ],
          xAxes: [ { display: false } ]
        }
      }
    });
    {{end}}
  </script>
{{end}}
