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
						<input type="text" class="form-control" id="port-input" name="port" placeholder="default: 5000">
					</fieldset>
					<fieldset class="form-group">
						<label for="scheme-input">Scheme</label>
						<div class="form-check form-check-inline">
							<label id="http" class="radio-inline">
								<input type="radio" name="scheme" value="http">HTTP
							</label>
							<label id="https" class="radio-inline">
								<input type="radio" name="scheme" value="https">HTTPS
							</label>
						</div>
						<div id="use-insecure" class="alert alert-danger" style="margin-top:10px; display:none;">
							<label class="checkbox-inline">
								<input type="checkbox" name="skip-tls-validation">Skip TLS Validation (required for self signed certs)
							</label>
						</div>
					</fieldset>
					<fieldset class="form-group">
						<label for="interval-input">Refresh Interval (seconds)</label>
						<input type="text" class="form-control" id="interval-input" name="interval" placeholder="default: 10">
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
