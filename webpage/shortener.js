var baseURL = "https://base.url/"
$('#url-shortener-form').submit(function(event) {
  event.preventDefault();

  $("#shortID").empty();
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
      $("#shortID").html("<a href=\"" + shortURL + "\">" + shortURL + "</a>");
    },
    error: function(data) {
      toggleProgress();
      $("#shortID").html("Unable to shorten URL");
      console.log("failure " + data);
    }
  });
});

function toggleProgress() {
  $(".loader").toggle();
}


$('#custom-cb').change(function() {
  $(".customize-id").toggle();
  $(".custom-URL").toggle();
});
