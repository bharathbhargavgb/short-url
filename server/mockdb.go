package main

import (
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
    return &dynamodb.GetItemOutput {
        Item: nil,
    }, nil
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

