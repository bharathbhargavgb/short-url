package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "fmt"
    "regexp"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

var tinyIDRegexp = regexp.MustCompile(`[a-zA-Z]{1,8}`)
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

type url struct {
    TinyID string `json:"TinyID"`
    URL string `json:"URL"`
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    switch req.HTTPMethod {
        case "GET":
            return expand(req)
        case "POST":
            return shorten(req)
        default:
            return clientError(http.StatusMethodNotAllowed)
    }
}

func expand(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Get the `TinyID` query string parameter from the request and
    // validate it.
    TinyID := req.QueryStringParameters["TinyID"]
    if !tinyIDRegexp.MatchString(TinyID) {
        return clientError(http.StatusBadRequest)
    }

    // Fetch the url record from the database based on the TinyID value.
    link, err := getItem(TinyID)
    if err != nil {
        return serverError(err)
    }
    if link == nil {
        return clientError(http.StatusNotFound)
    }

    // The APIGatewayProxyResponse.Body field needs to be a string, so
    // we marshal the url record into JSON.
    js, err := json.Marshal(link)
    if err != nil {
        return serverError(err)
    }

    // Return a response with a 200 OK status and the JSON book record
    // as the body.
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil
}

func shorten(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
        return clientError(http.StatusNotAcceptable)
    }

    link := new(url)
    err := json.Unmarshal([]byte(req.Body), link)
    if err != nil {
        errorLogger.Println("Unmarshal went wrong", err)
        return clientError(http.StatusUnprocessableEntity)
    }
    errorLogger.Println("Unmarshalled struct ", link)

    if !tinyIDRegexp.MatchString(link.TinyID) {
        errorLogger.Println("Regex does not match: ", link.TinyID)
        return clientError(http.StatusBadRequest)
    }
    if link.TinyID == "" || link.URL == "" {
        errorLogger.Println("TinyID or URL is empty ", link.TinyID, link.URL)
        return clientError(http.StatusBadRequest)
    }

    err = putItem(link)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: 201,
        Headers:    map[string]string{"Location": fmt.Sprintf("/shorten?TinyID=%s", link.TinyID)},
    }, nil
}

// Add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func serverError(err error) (events.APIGatewayProxyResponse, error) {
    errorLogger.Println(err.Error())

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Body:       http.StatusText(http.StatusInternalServerError),
    }, nil
}

// Similarly add a helper for send responses relating to client errors.
func clientError(status int) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: status,
        Body:       http.StatusText(status),
    }, nil
}

func main() {
    lambda.Start(router)
}
