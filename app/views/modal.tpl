<div id="pull-modal" class="modal fade" role="dialog">
  <div class="modal-dialog" style="z-index:999">
    <!-- Modal content-->
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal">&times;</button>
        <h4 class="modal-title">Download Blob</h4>
      </div>
      <div class="modal-body">
        <form class="form-horizontal">
          <div class="form-group">
            <label for="inputEmail3" class="col-sm-2 control-label">Email</label>
            <div class="col-sm-10">
              <input type="email" class="form-control" id="inputEmail3" placeholder="Email">
            </div>
          </div>
          <div class="form-group">
            <label for="inputPassword3" class="col-sm-2 control-label">Password</label>
            <div class="col-sm-10">
              <input type="password" class="form-control" id="inputPassword3" placeholder="Password">
            </div>
          </div>
          <div class="form-group">
            <div class="col-sm-offset-2 col-sm-10">
              <div class="checkbox">
                <label>
                  <input type="checkbox"> Remember me
                </label>
              </div>
            </div>
          </div>
          <div class="form-group">
            <div class="col-sm-offset-2 col-sm-10">
              <button type="submit" class="btn btn-default">Sign in</button>
            </div>
          </div>
        </form>
        <a href="{{$.registry.Scheme}}://{{$.registry.Name}}:{{$.registry.Port}}/{{$.registry.Version}}/{{$.repositoryName}}/{{$.tagInfo.Name}}" download><i class="fa fa-download"></i></a>
        <a href="{{$.registry.Scheme}}://{{$.registry.Name}}:{{$.registry.Port}}/{{$.registry.Version}}/{{$.repositoryName}}/{{$.tagInfo.Name}}" download><i class="fa fa-upload"></i></a>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
      </div>
    </div>
  </div>
</div>
<div id="push-modal" class="modal fade" role="dialog">
  <div class="modal-dialog" style="z-index:999">
    <!-- Modal content-->
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal">&times;</button>
        <h4 class="modal-title">Download Blob</h4>
      </div>
      <div class="modal-body">
          <a href="{{$.registry.Scheme}}://{{$.registry.Name}}:{{$.registry.Port}}/{{$.registry.Version}}/{{$.repositoryName}}/{{$.tagInfo.Name}}" download><i class="fa fa-upload"></i></a>
        <p>Some text in the modal.</p>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
      </div>
    </div>
  </div>
</div>
