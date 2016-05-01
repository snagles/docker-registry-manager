<div id="new-registry-modal" class="modal fade" role="dialog">
  <div class="modal-dialog" style="z-index:999">
    <!-- Modal content-->
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal">&times;</button>
        <h4 class="modal-title">Add new registry</h4>
      </div>
      <div class="modal-body">
        <form action="/registries/add" method="post">
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
              <input type="submit" class="btn btn-success">
              <button type="button" class="btn btn-danger" data-dismiss="modal">Cancel</button>
            </div>
        </form>
      </div>
    </div>
  </div>
</div>
