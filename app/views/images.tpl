{{template "base/base.html" .}}
{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li>
          <a href="/">Home</a>
        </li>
        <li>
          <a href="/registries">Registries</a>
        </li>
        <li>
          <a class="registry-name" href="/registries/{{.registryName}}/repositories">{{.registryName}}</a>
        </li>
        <li>
          <a href="/registries/{{.registryName}}/repositories">Repositories</a>
        </li>
        <li>
          <a class="registry-name" href="/registries/{{.registryName}}/repositories/{{.repositoryNameEncode}}/tags">{{.repositoryName}}</a>
        </li>
        <li>
          <a class="registry-name" href="/registries/{{.registryName}}/repositories/{{.repositoryNameEncode}}/tags">Tags</a>
        </li>
        <li class="active">{{.tagName}}</li>
      </ol>
    </div>
    <div class="row">
      <div class="col-md-12">
        <h1>{{.tagName}} <small> {{.repositoryName}}</small></h1>
      </div>
    </div>
    <div class="content-block white-bg col-lg-7 col-md-12 col-sm-12">
      <div>
        <ul class="nav nav-tabs" role="tablist">
          <li role="presentation" class="active">
            <a href="#overview" aria-controls="overview" role="tab" data-toggle="tab">Overview</a>
          </li>
          <li role="presentation">
            <a href="#build" aria-controls="build" role="tab" data-toggle="tab">Build</a>
          </li>
          <li role="presentation" class="disabled">
            <a href="#inspect" aria-controls="inspect" role="tab" data-toggle="tab">Inspect</a>
          </li>
          <div id="keywords" style="float:right;">
            {{range $keyword, $keywordInfo := .labels}}
              <a class="label keyword-label {{$keywordInfo.Color}}" data-label-color="{{$keywordInfo.Color}}" data-keyword="{{$keyword}}">
                <i class="fa {{$keywordInfo.Icon}}"></i>
                {{$keyword}}</a>
            {{end}}
          </div>
        </ul>
      </div>
      <div class="row">
        <div class="tab-content">
          <div role="tabpanel" class="col-md-12 tab-pane active border-between" id="overview">
            <div class="col-md-6">
              <div class="row col-md-10">
                <h4>Overview</h4>
                <ul class="well well-md">
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
              <div class="row col-md-10">
                <h4>Dockerhub
                  <a href="{{.dockerHub.ImageURL}}">
                    <i style="vertical-align:bottom; color:black;" class="fa fa-external-link" aria-hidden="true"></i>
                  </a>
                </h4>
                <div class="well well-md">
                  {{if .dockerHub.Error }}
                    <strong>Unable to compare image to Dockerhub</strong>
                    <small>{{.hubErr}}</small>
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
            <div class="col-md-6">
              <canvas id="stages-chart" width="400px" height="400px"></canvas>
            </div>
          </div>
          <div role="tabpanel" class="tab-pane table-responsive" id="build">
            <table id="datatable" class="table" cellspacing="0" width="100%">
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
                      <code>
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
  </div>
  <div class="content-block white-bg col-md-4 col-md-offset-1 col-sm-12">
    <ul class="nav nav-tabs" role="tablist">
      <li role="presentation" class="active">
        <a href="#actions" aria-controls="actions" role="tab" data-toggle="tab">Actions</a>
      </li>
      <li role="presentation" class="disabled">
        <a href="#timeline" aria-controls="timeline" role="tab" data-toggle="tab">Activity</a>
      </li>
    </ul>
    <div role="tabpanel" class="tab-pane active" id="actions">
      <div class="row col-md-12">
        <h4>Pull</h4>
        <ul class="well well-md">
          <code> docker pull {{.registryName}}/{{.repositoryName}}:{{.tagName}} </code>
        </ul>
      </div>
      <div class="row col-md-12">
        <h4>Push</h4>
        <ul class="well well-md">
          <code>
            docker tag
            {{.repositoryName}}:{{.tagName}}
            {{.registryName}}/{{.repositoryName}}:{{.tagName}}
            <br>
            docker push
            {{.registryName}}/{{.repositoryName}}:{{.tagName}}
          </code>
        </ul>
      </div>
      <div class="row col-md-12">
        <h4>Compare to Local Image</h4>
        <ul class="well well-md">
          <code class="bash">
            docker save
            {{.registryName}}/{{.repositoryName}}:{{.tagName}}
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
