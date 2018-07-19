{{template "base/base.html" .}}
{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <nav aria-label="breadcrumb">
        <ol class="breadcrumb">
          <li class="breadcrumb-item"><a href="/">Home</a></li>
          <li class="breadcrumb-item"><a href="/registries">Registries</a></li>
          <li class="breadcrumb-item"><a class="registry-name" href="/registries/{{.registryName}}/repositories">{{.registryName}}</a></li>
          <li class="breadcrumb-item"><a href="/registries/{{.registryName}}/repositories" class="registry-name">Repositories</a></li>
          <li class="breadcrumb-item"><a class="registry-name" href="/registries/{{.registryName}}/repositories/">{{.repositoryName}}</a></li>
          <li class="breadcrumb-item"><a class="registry-name" href="/registries/{{.registryName}}/repositories/{{.repositoryNameEncode}}/tags">Tags</a>
          <li class="breadcrumb-item active"><a class="registry-name" aria-current="page">{{.tagName}}</a></li>
        </ol>
      </nav>
    </div>
    <div class="row-fluid">
      <div class="ml-2">
        <h1>{{.tagName}} <small class="text-muted"> {{.repositoryName}}</small></h1>
      </div>
    </div>
    <div class="container-fluid-width row ml-2 mr-4">
      <div class="content-block white-bg col-lg-7 col-md-12 col-sm-12">
          <div class="container">
            <div class="row">
              <ul class="nav nav-tabs col-lg-5" role="tablist">
                <li class="nav-item">
                  <a class="nav-link active" href="#overview" aria-controls="overview" role="tab" data-toggle="tab">Overview</a>
                </li>
                <li class="nav-item">
                  <a class="nav-link" href="#build" aria-controls="build" role="tab" data-toggle="tab">Build</a>
                </li>
                <li class="nav-item">
                  <a class="nav-link disabled" href="#inspect" aria-controls="inspect" role="tab" data-toggle="tab">Inspect</a>
                </li>
              </ul>
              <div id="keywords" class="col border-bottom d-flex justify-content-end align-items-center">
                {{range $keyword, $keywordInfo := .labels}}
                  <div class="ml-1 badge keyword-label {{$keywordInfo.Color}}" data-label-color="{{$keywordInfo.Color}}" data-keyword="{{$keyword}}"> <i class="fa {{$keywordInfo.Icon}}"></i> <span>{{$keyword}}</span></div>
                {{end}}
              </div>
            </div>
          </div>
          <div class="tab-content">
            <div role="tabpanel" class="tab-pane active border-between" id="overview">
              <div class="container">
                <div class="row">
                  <div class="col mt-2">
                    <div class="col-12 border-0 card row">
                      <div class="card-block mt-2">
                        <h7 class="pl-0 col-md-12">Metadata</h7>
                        <div class="card bg-light rounded p-3">
                          <ul class="ml-0 pl-0">
                            <li>
                              <strong>Layers:</strong>
                              {{len .tag.DeserializedManifest.Layers}}
                            </li>
                            <li>
                              <strong>Size:</strong>
                              {{bytefmt .tag.Size}}
                            </li>
                            <li>
                              <strong>OS:</strong>
                              {{.tag.Os}}
                            </li>
                            <li>
                              <strong>Architecture:</strong>
                              {{.tag.Architecture}}
                            </li>
                            <li>
                              <strong>Docker Version:</strong>
                              {{.tag.DockerVersion}}
                            </li>
                            <li>
                              <strong>Created:</strong>
                              {{timeAgo .tag.Created}}
                            </li>
                            <li>
                              <strong>Last Updated:</strong>
                              {{timeAgo .tag.LastModified}}
                            </li>
                          </ul>
                        </div>
                      </div>
                    </div>
                    <div class="col-12 border-0 card row">
                      <div class="card-block mt-2">
                        <h7 class="pl-0 col-md-12">Dockerhub {{if .dockerHub.ImageURL}} <a href="{{.dockerHub.ImageURL}}"> <i style="vertical-align:bottom; color:black;" class="fa fa-external-link" aria-hidden="true"></i> </a> {{end}}</h7>
                        <div class="card bg-light rounded p-3">
                          {{if .dockerHub.Error }}
                            <strong>Unable to compare image to Dockerhub</strong>
                            <small class="text-muted">{{.hubErr}}</small>
                            <div>
                              <span>{{.dockerHub.Error}}</span></div>
                          {{else}}
                            {{ $length := len .dockerHub.DiffLayers }}
                            {{ if eq $length 0 }}
                              <strong>Image is up to date.</strong>
                              <div>
                                <span>Both local and remote images have
                                  {{len .tag.DeserializedManifest.Layers}}
                                  layers and are
                                  {{bytefmt .tag.Size}}.</span></div>
                            {{else}}
                              <div>
                                <strong>Image is not up to date.</strong>
                              </div>
                              <div>
                                <span>There are
                                  {{len .dockerHub.DiffLayers}}
                                  different layers out of
                                  {{len .tag.DeserializedManifest.Layers}}
                                  total layers.</span></div>
                              {{if ne .tag.Size .dockerHub.Size}}
                                {{if gt .tag.Size .dockerHub.Size}}
                                  <div>
                                    <span>The Dockerhub image is
                                      {{bytefmtdiff .dockerHub.Size .tag.Size}}
                                      smaller</span></div>
                                {{else if lt .tag.Size .dockerHub.Size}}
                                  <div>
                                    <span>The Dockerhub image is
                                      {{bytefmtdiff .dockerHub.Size .tag.Size}}
                                      larger</span></div>
                                {{else}}
                                  <div>
                                    <span>The images are the same total size
                                      {{bytefmt .tag.Size}}</span></div>
                                {{end}}
                              {{end}}
                            {{end}}
                          {{end}}
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="col mt-2">
                    <div class="row">
                      <canvas id="stages-chart" width="400px" height="400px"></canvas>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div role="tabpanel" class="tab-pane table-responsive" id="build">
              <table id="datatable" class="table borderless" cellspacing="0" width="100%">
                <thead>
                  <th>Stage</th>
                  <th>Digest</th>
                  <th>Command</th>
                  <th>Size
                    <i class="fa fa-question-circle" aria-hidden="true" data-toggle="tooltip" data-placement="top" title="Compressed tar.gz size"></i>
                  </th>
                  <th>Created</th>
                </thead>
                <tbody>
                  {{range $i, $history := $.tag.History}}
                    <tr data-tag-name="">
                      <td>{{oneIndex $i}}</td>
                      {{if not $history.EmptyLayer}}
                        <td data-toggle="tooltip" data-placement="top" title="{{$history.ManifestLayer.Digest.String}}">{{shortenDigest $history.ManifestLayer.Digest.String}}</td>
                      {{else}}
                        <td data-order="0">N/A</td>
                      {{end}}
                      <td style="max-width:400px;">
                        <code class="bash rounded">
                          {{range $i, $cmd := $history.Commands}}
                            {{ if ne (len $cmd.Keywords) 1 }}
                              <div data-keywords="{{$cmd.KeywordTags}}">{{$cmd.Cmd}}</div>
                            {{ else }}
                              <div>{{$cmd.Cmd}}</div>
                            {{end}}
                          {{end}}
                        </code>
                      </td>
                      {{if not $history.EmptyLayer}}
                        <td data-order="{{$history.ManifestLayer.Size}}"></i>{{bytefmt $history.ManifestLayer.Size}}</td>
                    {{else}}
                      <td data-order="0">0B</td>
                    {{end}}
                    <td>{{timeAgo $history.Created}}</td>
                  </tr>
                {{end}}
              </tbody>
            </table>
          </div>
        </div>
      </div>
      <div class="content-block white-bg offset-lg-1 col-lg-4 col-md-12 col-sm-12">
        <ul class="nav nav-tabs" role="tablist">
          <li class="nav-item">
            <a class="nav-link active" href="#actions" aria-controls="actions" role="tab" data-toggle="tab">Actions</a>
          </li>
          <li class="nav-item">
            <a class="nav-link disabled" href="#timeline" aria-controls="timeline" role="tab" data-toggle="tab">Activity</a>
          </li>
        </ul>
        <div class="tab-content">
          <div role="tabpanel" class="tab-pane active" id="actions">
            <div class="row col-md-12 mt-2">
              <h7 class="pl-0 col-md-12">Pull</h7>
              <ul class="col-md-12 card bg-light p-3">
                <code class="bash rounded"> docker pull {{.registryDisplayName}}/{{.repositoryName}}:{{.tagName}} </code>
              </ul>
            </div>
            <div class="row col-md-12">
              <h7 class="pl-0 col-md-12">Push</h7>
              <ul class="col-md-12 card bg-light p-3">
                <code class="bash rounded">
                  docker tag
                  {{.repositoryName}}:{{.tagName}}
                  {{.registryDisplayName}}/{{.repositoryName}}:{{.tagName}}
                  <br>
                  docker push
                  {{.registryDisplayName}}/{{.repositoryName}}:{{.tagName}}
                </code>
              </ul>
            </div>
            <div class="row col-md-12">
              <h7 class="pl-0 col-md-12">Compare to Local Image</h7>
              <ul class="col-md-12 card bg-light p-3">
                <code class="bash rounded">
                  docker save
                  {{.registryDisplayName}}/{{.repositoryName}}:{{.tagName}}
                  -o registry.tar.gz
                  <br>
                  tar tv -f registry.tar.gz > registry.txt
                  <br>
                  docker save
                  {{.repositoryName}}:{{.tagName}}
                  -o local.tar.gz
                  <br>
                  tar -tv -f local.tar.gz > local.txt
                  <br>
                  diff registry.txt local.txt
                </code>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<script>
  $('#image-tabs li a').click(function (e) {
    e.preventDefault()
    $(this).tab('show')
  });

  $(document).ready(function () {
    $('td code').each(function (i, block) {
      hljs.highlightBlock(block);
    });

    $('ul code').each(function (i, block) {
      hljs.highlightBlock(block);
    });

    $('#datatable').DataTable({
      "order": [
        [0, "asc"]
      ],
      "searching": false,
      "autoWidth": false,
      "info": false,
      "paging": false,
      "columnDefs": [
        {
          "width": "39px",
          "targets": 0
        }, {
          "orderable": false,
          "width": "50px",
          "targets": 1
        }, {
          "orderable": false,
          "targets": 2
        }, {
          "width": "44px",
          "targets": 3
        }, {
          "targets": 4,
          "width": "50px"
        }
      ]
    });
    $(function () {
      $('[data-toggle="tooltip"]').tooltip({container: 'body'})
    })
    $(".keyword-label").click(function () {
      $('.nav-tabs a[href="#build"]').tab('show')
      var labelColor = $(this).data("label-color");
      var keyword = $(this).data("keyword");
      $('div[data-keywords~="' + keyword + '"]').each(function () {
        $(this).addClass(labelColor).css("color", "white").delay(5000).queue(function () {
          $(this).removeClass(labelColor).css("color", "").dequeue()
        });
      })
    });
  });

  var ctx = document.getElementById("stages-chart");
  var pieChart = new Chart(ctx, {
    type: 'pie',
    options: {
      segmentShowStroke: false,
      maintainAspectRatio: false,
      pieceLabel: {
        mode: 'label'
      },
      title: {
        display: true,
        fontSize: 18,
        fontStyle: "normal",
        fontColor: "#212121",
        text: "{{len .tag.History}} Build Stages",
        position: "top"
      },
      tooltips: {
        enabled: true,
        mode: 'single',
        callbacks: {
          label: function (tooltipItems, data) {
            var data = data.datasets[tooltipItems.datasetIndex].info[tooltipItems.index]
            var label = "Stage: " + data.stage + ", Command #" + tooltipItems.index;
            if (data.keywords != null) {
              label += ', Keywords: ' + data.keywords;
            }
            return label
          }
        }
      }
    },
    data: {
      datasets: {{.chart}}
    }
  })
</script>

{{end}}
