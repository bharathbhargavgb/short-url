package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "regexp"
    "errors"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

var DDBTable = os.Getenv("URI_STORE")
var shortIDRegexp = regexp.MustCompile(`[a-zA-Z]{1,8}`)
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

type shortURI struct {
    ShortID string `json:"shortID"`
    URI string `json:"URI"`
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
    shortIDInput := req.QueryStringParameters["shortID"]
    if !shortIDRegexp.MatchString(shortIDInput) {
        return clientError(http.StatusBadRequest)
    }

    dataStore := getStorage(DDBTable)
    shortItem, err := dataStore.getItem(shortIDInput)
    if err != nil {
        return serverError(err)
    }
    if shortItem == nil {
        return clientError(http.StatusNotFound)
    }

    httpResponseBody, err := json.Marshal(shortItem)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(httpResponseBody),
    }, nil
}

func shorten(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
        return clientError(http.StatusNotAcceptable)
    }

    shortItem := new(shortURI)
    err := json.Unmarshal([]byte(req.Body), shortItem)
    if err != nil {
        return clientError(http.StatusUnprocessableEntity)
    }

    if shortItem.URI == "" {
        return clientError(http.StatusBadRequest)
    }

    if shortItem.ShortID != "" && !shortIDRegexp.MatchString(shortItem.ShortID) {
        return clientError(http.StatusBadRequest)
    }

    dataStore := getStorage(DDBTable)
    if shortItem.ShortID == "" {
        shortItem.ShortID = getValidShortID(dataStore, shortItem.URI)
        if shortItem.ShortID == "" {
            return serverError(errors.New("getValidShortID() returned empty"))
        }
    }

    err = dataStore.putItem(shortItem)
    if err != nil {
        return serverError(err)
    }

    httpResponseBody, err := json.Marshal(shortItem)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: 201,
        Headers:    map[string]string{
            "Access-Control-Allow-Origin": "*",
        },
        Body: string(httpResponseBody),
    }, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
    errorLogger.Println(err.Error())

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Headers:    map[string]string{
            "Access-Control-Allow-Origin": "*",
        },
        Body:       http.StatusText(http.StatusInternalServerError),
    }, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: status,
        Headers:    map[string]string{
            "Access-Control-Allow-Origin": "*",
        },
        Body:       http.StatusText(status),
    }, nil
}

func main() {
    lambda.Start(router)
}
