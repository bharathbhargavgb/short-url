package main

import (
    "github.com/aws/aws-lambda-go/lambda"
)

type url struct {
    tinyID string `json:"tinyID"`
    URL string `json:"URL"`
}

func shorten() (*url, error) {
    link, err := getItem("dq")
    if err != nil {
        return nil, err
    }
    return link, nil
}

func main() {
    lambda.Start(shorten)
}
