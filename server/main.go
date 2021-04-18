package main

import (
    "github.com/aws/aws-lambda-go/lambda"
)

func shorten() (string, error) {
    shortString := "tiny"
    return shortString, nil
}

func main() {
    lambda.Start(shorten)
}
