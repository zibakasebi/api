package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	var loginRequest LoginRequest
	loginRequest.Username = os.Getenv("username")
	loginRequest.Password = os.Getenv("password")
	//convert struct to json
	jsonStr, err := json.Marshal(loginRequest)
	if err != nil {
		t.Error("error:", err)
	}

	var token = os.Getenv("token")

	tests := []struct {
		request    events.APIGatewayProxyRequest
		expect     string
		err        error
		statusCode int
	}{
		{
			// Test valid login
			request: events.APIGatewayProxyRequest{
				Body: string(jsonStr),
			},
			expect:     fmt.Sprintf(`{"token": "%s"}`, token),
			err:        nil,
			statusCode: 200,
		},
		{
			// Test invalid login
			request:    events.APIGatewayProxyRequest{},
			expect:     "",
			err:        nil,
			statusCode: 401,
		},
	}

	for i, test := range tests {
		response, err := Handler(context.Background(), test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
		log.Printf("Test %d: %s", i, response.Body)
	}
}
