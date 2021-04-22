package main

import (
    "testing"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type fakeDynamoDBClient struct {
    dynamodbiface.DynamoDBAPI
}

func createFakeStorage() *DBStore {
    return &DBStore {
        TableName: "DummyTable",
        Client: &fakeDynamoDBClient{},
    }
}

func (m *fakeDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
    return &dynamodb.GetItemOutput {
        Item: map[string]*dynamodb.AttributeValue {
           "TinyID": {
                S: aws.String("goog"),
           },
           "URI": {
               S: aws.String("https://www.google.com"),
           },
        },
    }, nil
}

func TestGetItem(t *testing.T) {
    db := createFakeStorage()
    want := "https://www.google.com"
    shortLink, _ := db.getItem("dummyString")
    if want != shortLink.URI {
        t.Errorf("got %q, want %q", shortLink.URI, want)
    }
}


