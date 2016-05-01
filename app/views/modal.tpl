<div id="dockerhub-modal" class="modal fade" role="dialog">
  <div class="modal-dialog" style="z-index:999">
    <!-- Modal content-->
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal">&times;</button>
        <h4 class="modal-title">Push to Dockerhub</h4>
      </div>
      <div class="modal-body">
        <div class ="row">Are you sure you want to push:</div>
          <code>{{.repositoryName}}:{{.tagInfo.Name}}</code>
        <div>to Dockerhub with</div>
          <code>{{.dockerAccountName}}:{{.tagInfo.Name}}</code>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
        <button type="button" class="btn btn-default" data-dismiss="modal">Push</button>
      </div>
    </div>
  </div>
</div>
