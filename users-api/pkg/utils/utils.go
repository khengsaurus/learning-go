package utils

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/khengsaurus/go-tutorials/users-api/pkg/consts"
	"github.com/khengsaurus/go-tutorials/users-api/pkg/models"
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func IsEmailValid(email string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]{1,64}@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return len(email) > 3 && len(email) <= 254 && rxEmail.MatchString(email)
}

func ApiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}

func ForwardAwsError(status int, err error) (*events.APIGatewayProxyResponse, error) {
	return ApiResponse(status, ErrorBody{
		aws.String(err.Error()),
	})
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
