<div id="new-registry-modal" class="modal fade" role="dialog">
	<div class="modal-dialog" role="document" style="z-index:999">
		<!-- Modal content-->
		<div class="modal-content">
			<div class="modal-header">
				<h4 class="modal-title">Add new registry</h4>
				<button type="button" class="close" data-dismiss="modal">&times;</button>
			</div>
			<div class="modal-body">
				<form id="registry-form" action="/registries/add" method="post">
					<fieldset class="form-group">
						<label for="name-input">Name</label>
						<input type="text" class="form-control" id="name-input" name="name" placeholder="registry1">
					</fieldset>
					<fieldset class="form-group">
						<label for="name-input">Display Name (optional)</label>
						<input type="text" class="form-control" id="displayName-input" name="name" placeholder="full-url:port/path/">
					</fieldset>
					<fieldset class="form-group">
						<label for="host-input">Host</label>
						<input type="text" class="form-control" id="host-input" name="host" placeholder="ex: 192.168.1.1 or testhost.com">
					</fieldset>
					<fieldset class="form-group">
						<label for="port-input">Port</label>
						<input type="text" class="form-control" id="port-input" name="port" placeholder="default: 5000">
					</fieldset>
					<fieldset class="form-group">
						<div for="scheme-input">Scheme</div>
						<div>
							<div class="form-check form-check-inline" id="http" class="radio-inline">
								<input class="form-check-input" type="radio" checked name="scheme" id="scheme" value="http">
								<label class="form-check-label" for="http">HTTP</label>
							</div>
							<div  class="form-check form-check-inline" id="https" class="radio-inline">
								<input class="form-check-input" type="radio" name="scheme" id="scheme" value="https">
								<label class="form-check-label" for="https">HTTPS</label>
							</div>
						</div>
						<div id="use-insecure" class="alert alert-danger" style="margin-top:10px; display:none;">
							<div class="form-check form-check-inline" id="http" class="radio-inline">
								<input class="form-check-input" type="checkbox" name="skip-tls-validation">
								<label class="form-check-label" for="skip-tls-validation">Skip TLS Validation (required for self signed certs)</label>
							</div>
						</div>
					</fieldset>
					<fieldset class="form-group">
						<label for="interval-input">Refresh Interval (seconds)</label>
						<input type="text" class="form-control" id="interval-input" name="interval" placeholder="default: 60">
					</fieldset>
					<fieldset class="form-group">
						<div class="form-check form-check-inline" id="http" class="radio-inline">
							<input class="form-check-input" type="checkbox" name="dockerhub-integration">
							<label class="form-check-label" for="dockerhub-integration">Compare images to hub.docker.com</label>
						</div>
						<small class="form-text text-muted">Every image tag is queried using the hub.docker.com API, and then compares layers and sizes</small>
					</fieldset>
					<fieldset class="form-group">
						<div class="form-check form-check-inline" id="http" class="radio-inline">
							<input class="form-check-input" type="checkbox" name="read-only">
							<label class="form-check-label" for="read-only">Read Only Mode</label>
						</div>
					</fieldset>
					<div class="modal-footer d-flex">
					  <div class="mr-auto"><button type="button" id="test" class="d-flex justify-content-start btn btn-warning text-white">Test</button></div>
					  <div><input type="submit" class="btn btn-success" id="add-registry" value="Submit"></div>
					  <div><button type="button" class="btn btn-danger" data-dismiss="modal">Cancel</button></div>
					</div>
				</form>
			</div>
		</div>
	</div>
</div>

<script>
	$(function () {
	  $('[data-toggle="tooltip"]').tooltip()
	});

	$('#https').click(function () {
		$("#use-insecure").show();
	});
	$('#http').click(function () {
		$("#use-insecure").hide();
	});
	$("#test").click(function () {
		var data = $('#registry-form').serialize();
		$.ajax({
			type: 'POST',
			url: '/registries/test',
			data: data,
			dataType: 'json',
			success: function (data) {
				$.each(data, function (index, element) {
					if (index == "is_available") {
						if (element == true) {
							$(".modal-body").append("<div class='alert alert-success'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a> <strong>Success!</strong> We've successfully made a connection to the registry. </div>");
						} else {
							$(".modal-body").append("<div class='alert alert-danger'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a> <strong>Failure!</strong> We could not connect to the registry. </div>");
						}
						$(".alert").alert();
						window.setTimeout(function () {
							$(".alert").alert('close');
						}, 5000);
					}
				});
			}
		});
	});
</script>
