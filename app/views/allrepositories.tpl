{{template "base/base.html" .}} {{define "body"}}
<div class="right-content-container">
	<div class="header">
		<nav aria-label="breadcrumb">
			<ol class="breadcrumb">
				<li class="breadcrumb-item"><a href="/">Home</a></li>
				<li class="breadcrumb-item"><a href="/registries" class="registry-name">Registries</a></li>
				<li class="breadcrumb-item active"><a class="registry-name" aria-current="page">Repositories</a></li>
			</ol>
		</nav>
	</div>
	<div class="row">
		<div class="col-md-12">
			<h1>All Repositories</h1>
		</div>
	</div>
	<div class="content-block white-bg">
		<div class="row">
			<table id="datatable" class="table table-striped table-borderless" cellspacing="0" width="100%">
				<thead>
					<th>Repository</th>
					<th>Registry</th>
					<th>Size
						<i class="fa fa-question-circle" aria-hidden="true" data-toggle="tooltip" data-placement="top" title="Compressed tar.gz size"></i>
					</th>
					<th>Tags</th>
				</thead>
				<tfoot>
					<th>Repository</th>
					<th>Registry</th>
					<th>Size
						<i class="fa fa-question-circle" aria-hidden="true" data-toggle="tooltip" data-placement="top" title="Compressed tar.gz size"></i>
					</th>
					<th>Tags</th>
				</tfoot>
				<tbody>
					{{range $registryName, $repositories := .repositories}} {{range $key, $repo := $repositories}}
					<tr>
						<td>
							<a href=/registries/{{$registryName}}/repositories/{{urlquery $repo.Name}}/tags>{{$repo.Name}}</span></td>
                  <td>{{$registryName}}</td>
                  <td data-order="{{$repo.Size}}">{{bytefmt $repo.Size}}</td>
                  <td>{{len $repo.Tags}}</td>
                </tr>
              {{end}}
            {{end}}
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <script>
    $(document).ready(function () {
      $('#datatable').DataTable({
        "order": [
          [1, "asc"]
        ],
        "pageLength": 25
      });
      $(function () {
        $('[data-toggle="tooltip"]').tooltip({container: 'body'})
      })
    });
  </script>
{{end}}
