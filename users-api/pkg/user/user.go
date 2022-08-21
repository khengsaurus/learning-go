package user

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/khengsaurus/go-tutorials/users-api/pkg/consts"
	"github.com/khengsaurus/go-tutorials/users-api/pkg/models"
	"github.com/khengsaurus/go-tutorials/users-api/pkg/utils"
)

func FetchUser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*models.User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(consts.ErrorFailedToFetchRecord)
	}

	item := new(models.User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(consts.ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]models.User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(consts.ErrorFailedToFetchRecord)
	}

	item := new([]models.User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		return nil, errors.New(consts.ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func PostHelper(u models.User, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*models.User, error) {
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(consts.ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(consts.ErrorCouldNotDynamoPutItem)
	}
	return &u, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*models.User, error) {
	var u models.User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(consts.ErrorInvalidUserData)
	}
	if !utils.IsEmailValid(u.Email) {
		return nil, errors.New(consts.ErrorInvalidEmail)
	}
	currentUser, _ := FetchUser(u.Email, tableName, dynaClient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(consts.ErrorUserAlreadyExists)
	}

	return PostHelper(u, tableName, dynaClient)
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*models.User,
	error,
) {
	var u models.User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(consts.ErrorInvalidEmail)
	}

	currentUser, _ := FetchUser(u.Email, tableName, dynaClient)
	if currentUser != nil && len(currentUser.Email) == 0 {
		return nil, errors.New(consts.ErrorUserDoesNotExist)
	}

	return PostHelper(u, tableName, dynaClient)
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(consts.ErrorCouldNotDeleteItem)
	}

	return nil
}
