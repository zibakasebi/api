package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func setupResponseCors() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
			"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		},
	}, nil
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "OPTIONS" {
		return setupResponseCors()
	}

	//convert string to struct
	var loginRequest LoginRequest
	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		fmt.Println("error:", err)
	}

	//get username and password from environment variables
	var username = os.Getenv("username")
	var password = os.Getenv("password")
	var token = os.Getenv("token")

	if username == loginRequest.Username && password == loginRequest.Password {
		return &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"access-control-allow-origin": "*",
			},
			Body: fmt.Sprintf(`{"token": "%s"}`, token),
		}, nil
	} else {
		return &events.APIGatewayProxyResponse{
			StatusCode: 401,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"access-control-allow-origin": "*",
			},
		}, nil
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(Handler)
}
