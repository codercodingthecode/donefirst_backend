package main

import (
	data "data"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"net/http"
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

// apiResponse is a utility function to create the response for the API Gateway
func apiResponse(status int, body interface{}) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return resp, nil
}

// saveData is a function to save the data to the DynamoDB table
func saveData(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var registration data.RegistrationResponse
	registrationError := json.Unmarshal([]byte(request.Body), &registration)
	if registrationError != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(registrationError.Error())})
	}

	dataErr := data.SaveRegistration(registration)
	if dataErr != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(dataErr.Error())})
	}

	return apiResponse(http.StatusOK, "Data saved")
}

// getData is a function to get the data from the DynamoDB table
func getData(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	payload, _ := request.PathParameters["id"]

	d, dataErr := data.GetRegistration(payload)
	if dataErr != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(dataErr.Error())})
	}

	if d.Id == "" {
		return apiResponse(http.StatusNotFound, ErrorBody{aws.String("No data found")})
	}

	return apiResponse(http.StatusOK, d)
}

// getCollection is a function to get a List of all the data from the DynamoDB table
func getCollection() (events.APIGatewayProxyResponse, error) {
	data, dataErr := data.GetRegistrations()
	if dataErr != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(dataErr.Error())})
	}

	return apiResponse(http.StatusOK, data)
}

// Handler  acting as a router for the requests
func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	switch request.RequestContext.HTTP.Method {
	case "POST":
		return saveData(request)
	case "GET":
		if request.RequestContext.HTTP.Path == "/registration/list" {
			return getCollection()
		} else {
			return getData(request)
		}
	default:
		return events.APIGatewayProxyResponse{Body: "Invalid request", StatusCode: 404}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
