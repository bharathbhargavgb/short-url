var baseURL = "https://base.url/"

$('#url-shortener-form').submit(function(event) {
  event.preventDefault();

  resetOutput();
  toggleProgress();

  var reqBody = {
    shortID: $("#custom-id").val(),
    URI: $("#url-input").val()
  };

  $.ajax({
    url: 'https://angh4tqiu8.execute-api.us-east-1.amazonaws.com/staging/shorten',
    type: 'POST',
    dataType: 'json',
    data: JSON.stringify(reqBody),
    headers: {
      'Content-Type': 'application/json',
    },
    success: function (data) {
      var shortURL = baseURL + data.shortID;
      toggleProgress();
      $("#output-success").show();
      $("#shortID").html(shortURL);
    },
    error: function(data) {
      toggleProgress();
      $("#output-error").show();
      $("#error-msg").html("Unable to shorten URL");
      console.log("failure " + data);
    }
  });
});

function toggleProgress() {
  $(".loader").toggle();
}

function resetOutput() {
  $("#shortID").empty();
  $("#error-msg").empty();
  $("#output-success").hide();
  $("#output-error").hide();
}


$('#custom-cb').change(function() {
  $(".customize-id").toggle();
  $(".custom-URL").toggle();
});

$('.copy-btn').click(function() {
  var shortURL = $('#shortID').html();
  copyToClipboard(shortURL);
});

function copyToClipboard(content) {
  navigator.clipboard.writeText(content).then(function() {
    // success: no-op
  }, function(err) {
    console.error('Async: Could not copy text: ', err);
  });
}
