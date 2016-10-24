{{template "base/base.html" .}}
{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li><a href="/">Home</a></li>
        <li><a href="/registries">Registries</a></li>
        <li><a class="registry-name" href="/registries/{{.registryName}}/repositories">{{.registryName}}</a></li>
        <li><a href="/registries/{{.registryName}}/repositories">Repositories</a></li>
        <li><a class="registry-name" href="/registries/{{.registryName}}/repositories/{{.repositoryNameEncode}}/tags">{{.repositoryName}}</a></li>
        <li class="active">Tags</li>
      </ol>
    </div>
    <div class="content-block white-bg">
      <div class="row">
        <h1>{{.repositoryName}}</h1>
        <hr>
      </div>
      <div class="row">
        <form id="delete-tags">
          <table id="datatable" class="table table-striped table-bordered display select" cellspacing="0" width="100%" cellspacing="0" width="100%">
             <thead>
                <tr>
                  <th><input name="select_all" value="1" type="checkbox"></th>
                  <th>Tags:</th>
                  <th>Updated:</th>
                  <th>Size:</th>
                  <th>Layers:</th>
                </tr>
             </thead>
             <tfoot>
                <tr>
                  <th></th>
                  <th>Tags:</th>
                  <th>Updated:</th>
                  <th>Size:</th>
                  <th>Layers:</th>
                </tr>
             </tfoot>
             <tbody>
               {{range $key, $tag := .tags}}
              <tr data-tag-name="{{$tag}}">
                <td></td>
                <td ><a href=/registries/{{$.registryName}}/repositories/{{$.repositoryName}}/tags/{{$tag}}/images>{{$tag}}</span></td>
                <td data-order=""></td>
              </tr>
              {{end}}
            </tbody>
        </table>
        <p><button class="btn btn-danger">Delete</button></p>

</form>
      </div>
    </div>
  </div>

  <script>
  //
  // Updates "Select all" control in a data table
  //
  function updateDataTableSelectAllCtrl(table){
     var $table             = table.table().node();
     var $chkbox_all        = $('tbody input[type="checkbox"]', $table);
     var $chkbox_checked    = $('tbody input[type="checkbox"]:checked', $table);
     var chkbox_select_all  = $('thead input[name="select_all"]', $table).get(0);

     // If none of the checkboxes are checked
     if($chkbox_checked.length === 0){
        chkbox_select_all.checked = false;
        if('indeterminate' in chkbox_select_all){
           chkbox_select_all.indeterminate = false;
        }

     // If all of the checkboxes are checked
     } else if ($chkbox_checked.length === $chkbox_all.length){
        chkbox_select_all.checked = true;
        if('indeterminate' in chkbox_select_all){
           chkbox_select_all.indeterminate = false;
        }

     // If some of the checkboxes are checked
     } else {
        chkbox_select_all.checked = true;
        if('indeterminate' in chkbox_select_all){
           chkbox_select_all.indeterminate = true;
        }
     }
  }

  $(document).ready(function (){
     // Array holding selected row IDs
     var rows_selected = [];
     var table = $('#datatable').DataTable({
        'columnDefs': [{
           'targets': 0,
           'searchable':false,
           'orderable':false,
           'width':'1%',
           'className': 'dt-body-center',
           'render': function (data, type, full, meta){
               return '<input type="checkbox">';
           }
        }],
        'order': [2, 'desc'],
        'rowCallback': function(row, data, dataIndex){
           // Get row ID
           var rowId = data[0];

           // If row ID is in the list of selected row IDs
           if($.inArray(rowId, rows_selected) !== -1){
              $(row).find('input[type="checkbox"]').prop('checked', true);
              $(row).addClass('selected');
           }
        }
     });

     // Handle click on checkbox
     $('#datatable tbody').on('click', 'input[type="checkbox"]', function(e){
        var $row = $(this).closest('tr');

        // Get row data
        var data = table.row($row).data();

        // Get row ID
        var rowId = data[0];

        // Determine whether row ID is in the list of selected row IDs
        var index = $.inArray(rowId, rows_selected);

        // If checkbox is checked and row ID is not in list of selected row IDs
        if(this.checked && index === -1){
           rows_selected.push(rowId);

        // Otherwise, if checkbox is not checked and row ID is in list of selected row IDs
        } else if (!this.checked && index !== -1){
           rows_selected.splice(index, 1);
        }

        if(this.checked){
           $row.addClass('selected');
        } else {
           $row.removeClass('selected');
        }

        // Update state of "Select all" control
        updateDataTableSelectAllCtrl(table);

        // Prevent click event from propagating to parent
        e.stopPropagation();
     });

     // Handle click on table cells with checkboxes
     $('#datatable').on('click', 'tbody td, thead th:first-child', function(e){
        $(this).parent().find('input[type="checkbox"]').trigger('click');
     });

     // Handle click on "Select all" control
     $('thead input[name="select_all"]', table.table().container()).on('click', function(e){
        if(this.checked){
           $('#datatable tbody input[type="checkbox"]:not(:checked)').trigger('click');
        } else {
           $('#datatable tbody input[type="checkbox"]:checked').trigger('click');
        }

        // Prevent click event from propagating to parent
        e.stopPropagation();
     });

     // Handle table draw event
     table.on('draw', function(){
        // Update state of "Select all" control
        updateDataTableSelectAllCtrl(table);
     });

     // Handle form submission event
     $('#delete-tags').on('submit', function(e){
        var form = this;

        var x = document.getElementsByClassName("selected");
        $(x).each(function(index) {
          var tagName = $(this).data("tag-name");
          $.ajax({
            type: "POST",
            url: "/registries/{{.registryName}}/repositories/{{.repositoryNameEncode}}/tags/"+tagName+"/delete",
            success: function() {
                  $("#delete-tags").append("<div class='alert alert-success'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a> <strong>Success!</strong> We've successfully deleted "+tagName+" from the registry. </div>");
                  $(".alert").alert();
                  window.setTimeout(function() { $(".alert").alert('close'); }, 5000);
                  $("tr[data-tag-name="+tagName+"]").remove();
            },
            error: function() {
                  $("#delete-tags").append("<div class='alert alert-danger'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a> <strong>Failure!</strong> We were unable to delete "+tagName+" from the registry. Make sure the delete option is enabled on your registry!</div>");
                  $(".alert").alert();
                  window.setTimeout(function() { $(".alert").alert('close'); }, 5000);
            }
          });
        });

        // Prevent actual form submission
        e.preventDefault();
     });
  });
  </script>

{{end}}
