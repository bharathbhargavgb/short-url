package main

import (
    "fmt"
    "testing"
    "errors"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type fakeDynamoDBClient struct {
    dynamodbiface.DynamoDBAPI
}

type shortItemKey struct {
    TinyID string
}

func createFakeStorage() *DBStore {
    return &DBStore {
        TableName: "DummyTable",
        Client: &fakeDynamoDBClient{},
    }
}

func (m *fakeDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
    shortItem := new(shortItemKey)
    dynamodbattribute.UnmarshalMap(input.Key, shortItem)

    if shortItem.TinyID == "goog" {
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
    return nil, errors.New("Key not found")
}

func (m *fakeDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
    shortItemInput := new(shortURI)
    err := dynamodbattribute.UnmarshalMap(input.Item, shortItemInput)
    if err != nil {
        return nil, errors.New("Unable to unmarshalMap from struct")
    }
    if shortItemInput.TinyID == "goog" && shortItemInput.URI == "https://www.google.com" {
        return nil, nil
    }
    return nil, errors.New("Malformed input")
}

func TestGetItem(t *testing.T) {
    db := createFakeStorage()
    want := "https://www.google.com"
    shortLink, err := db.getItem("goog")
    if err != nil {
        t.Errorf("getItem returned error: %q", err)
    }
    if want != shortLink.URI {
        t.Errorf("Got %q, Want %q", shortLink.URI, want)
    }
}

func TestGetItemMissingKey(t *testing.T) {
    db := createFakeStorage()
    shortLink, err := db.getItem("randomString")
    if err == nil {
        t.Errorf("Want `Key not found`, Got %q", err)
    }
    if shortLink != nil {
        t.Error(fmt.Sprintf("Want nil, Got %+v", *shortLink))
    }
}

func TestPutItem(t *testing.T) {
    db := createFakeStorage()
    shortItem := &shortURI {
        TinyID: "goog",
        URI: "https://www.google.com",
    }
    err := db.putItem(shortItem)
    if err != nil {
        t.Error("Error inserting to DB")
    }
}

