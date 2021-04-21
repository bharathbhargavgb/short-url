package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))
const URITable = "URIStore"

func getItem(tinyIDInput string) (*shortURI, error) {
    input := &dynamodb.GetItemInput {
        TableName: aws.String(URITable),
        Key: map[string]*dynamodb.AttributeValue {
            "TinyID": {
                S: aws.String(tinyIDInput),
            },
        },
    }

    result, err := db.GetItem(input)
    if err != nil {
        return nil, err
    }
    if result.Item == nil {
        return nil, nil
    }

    // The result.Item object returned has the underlying type
    // map[string]*AttributeValue. We can use the UnmarshalMap helper
    // to parse this straight into the fields of a struct.
    shortItem := new(shortURI)
    err = dynamodbattribute.UnmarshalMap(result.Item, shortItem)
    if err != nil {
        return nil, err
    }

    return shortItem, nil
}

func putItem(shortItem *shortURI) error {
    input := &dynamodb.PutItemInput{
        TableName: aws.String(URITable),
        Item: map[string]*dynamodb.AttributeValue{
            "TinyID": {
                S: aws.String(shortItem.TinyID),
            },
            "URI": {
                S: aws.String(shortItem.URI),
            },
        },
    }

    _, err := db.PutItem(input)
    return err
}
