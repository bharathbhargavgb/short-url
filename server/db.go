package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DBStore struct {
    TableName string
    Client dynamodbiface.DynamoDBAPI
}

func getStorage(tableName string) *DBStore {
    return &DBStore {
        TableName: tableName,
        Client: dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1")),
    }
}


func (db *DBStore) getItem(tinyIDInput string) (*shortURI, error) {
    input := &dynamodb.GetItemInput {
        TableName: aws.String(db.TableName),
        Key: map[string]*dynamodb.AttributeValue {
            "TinyID": {
                S: aws.String(tinyIDInput),
            },
        },
    }

    result, err := db.Client.GetItem(input)
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

func (db *DBStore) putItem(shortItem *shortURI) error {
    input := &dynamodb.PutItemInput{
        TableName: aws.String(db.TableName),
        Item: map[string]*dynamodb.AttributeValue{
            "TinyID": {
                S: aws.String(shortItem.TinyID),
            },
            "URI": {
                S: aws.String(shortItem.URI),
            },
        },
    }

    _, err := db.Client.PutItem(input)
    return err
}
