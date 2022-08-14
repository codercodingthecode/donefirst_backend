package data

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const TableName = "RegistrationTable"

type RegistrationResponse struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Dob         json.Number `json:"dob"`
	Phone       string      `json:"phone"`
	Email       string      `json:"email"`
	Address     string      `json:"address"`
	PhotoDl     string      `json:"photoDl"`
	Appointment json.Number `json:"appointment"`
}

func CreateTableSession() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return dynamodb.New(sess)
}

// GetRegistration returns the registration data for the given id and operate on DynamoDB
func GetRegistration(id string) (*RegistrationResponse, error) {
	registration := new(RegistrationResponse)
	db := CreateTableSession()

	input := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	result, err := db.GetItem(input)
	if err != nil {
		return registration, err
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, registration)

	if err != nil {
		return registration, err
	}

	return registration, nil
}

// GetRegistrations  returns all the registration data from DynamoDB
func GetRegistrations() ([]RegistrationResponse, error) {
	db := CreateTableSession()
	input := &dynamodb.ScanInput{
		TableName: aws.String(TableName),
	}
	result, err := db.Scan(input)
	if err != nil {
		return nil, err
	}
	var registrations []RegistrationResponse
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &registrations)
	if err != nil {
		return nil, err
	}
	return registrations, nil
}

// SaveRegistration saves the registration data to DynamoDB
func SaveRegistration(registration RegistrationResponse) error {
	data, dataErr := dynamodbattribute.MarshalMap(registration)
	if dataErr != nil {
		return dataErr
	}
	db := CreateTableSession()
	input := &dynamodb.PutItemInput{
		Item:      data,
		TableName: aws.String(TableName),
	}
	_, err := db.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}
