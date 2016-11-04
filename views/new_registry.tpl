<div id="new-registry-modal" class="modal fade" role="dialog">
  <div class="modal-dialog" style="z-index:999">
    <!-- Modal content-->
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal">&times;</button>
        <h4 class="modal-title">Add new registry</h4>
      </div>
      <div class="modal-body">
        <form id="registry-form" action="/registries/add" method="post">
            <fieldset class="form-group">
              <label for="host-input">Host</label>
              <input type="text" class="form-control" id="host-input" name="host" placeholder="ex: 192.168.1.1 or testhost.com">
            </fieldset>
            <fieldset class="form-group">
              <label for="port-input">Port</label>
              <input type="text" class="form-control" id="port-input" name="port" placeholder="ex: 5000">
            </fieldset>
            <fieldset class="form-group">
              <label for="scheme-input">Scheme</label>
              <input type="text" class="form-control" id="scheme-input" name="scheme" placeholder="ex: https">
            </fieldset>
            <div class="modal-footer">
              <button style="float:left;" type="button" id="test" class="btn btn-warning">Test</button>
              <input type="submit" class="btn btn-success" id="add-registry" value="Submit">
              <button type="button" class="btn btn-danger" data-dismiss="modal">Cancel</button>
            </div>
        </form>
      </div>
    </div>
  </div>
</div>

<script>
$("#test").click(function() {
   var data = $('#registry-form').serialize();
  $.ajax({
      type: 'POST',
      url: '/registries/test',
      data: data,
      dataType: 'json',
      success: function (data) {
          $.each(data, function(index, element) {
            if (index == "is_available") {
              if (element == true) {
                $(".modal-body").append("<div class='alert alert-success'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a> <strong>Success!</strong> We've successfully made a connection to the registry. </div>");
              } else {
                $(".modal-body").append("<div class='alert alert-danger'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a> <strong>Failure!</strong> We could not connect to the registry. </div>");
              }
              $(".alert").alert();
              window.setTimeout(function() { $(".alert").alert('close'); }, 5000);
            }
          });
      }
  });
});
</script>
