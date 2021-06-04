var baseURL = "https://base.url/"

$("#submit").click(function() {
  console.log("onclick");
  var reqBody = {
    shortID: "",
    URI: $("#url").val()
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
        $("#shortID").html("<a href=\"" + shortURL + "\">" + shortURL + "</a>");
      },
      failure: function(data) {
        console.log("failure " + data);
      }
  });
});

