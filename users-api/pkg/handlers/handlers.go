package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/khengsaurus/go-tutorials/users-api/pkg/user"
	"github.com/khengsaurus/go-tutorials/users-api/pkg/utils"
)

var ErrorMethodNotAllowed = "Method not allowed"

func GetUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	if len(email) > 0 {
		result, err := user.FetchUser(email, tableName, dynaClient)
		if err != nil {
			return utils.ForwardAwsError(http.StatusBadRequest, err)
		}
		return utils.ApiResponse(http.StatusOK, result)
	}

	result, err := user.FetchUsers(tableName, dynaClient)
	if err != nil {
		return utils.ForwardAwsError(http.StatusBadRequest, err)
	}
	return utils.ApiResponse(http.StatusOK, result)
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := user.CreateUser(req, tableName, dynaClient)
	if err != nil {
		return utils.ForwardAwsError(http.StatusInternalServerError, err)
	}
	return utils.ApiResponse(http.StatusCreated, result)
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := user.UpdateUser(req, tableName, dynaClient)
	if err != nil {
		return utils.ForwardAwsError(http.StatusInternalServerError, err)
	}
	return utils.ApiResponse(http.StatusOK, result)
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	err := user.DeleteUser(req, tableName, dynaClient)
	if err != nil {
		return utils.ForwardAwsError(http.StatusInternalServerError, err)
	}
	return utils.ApiResponse(http.StatusOK, nil)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return utils.ApiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
