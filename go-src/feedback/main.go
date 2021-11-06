package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FeedbackRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Message string `json:"message"`
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
	var feedbackRequest FeedbackRequest
	err := json.Unmarshal([]byte(request.Body), &feedbackRequest)
	if err != nil {
		fmt.Println("error:", err, request.Body)
	}

	col := getClient().Database("test").Collection("feedback")

	_, err = col.InsertOne(ctx, feedbackRequest)
	if err != nil {
		return nil, errors.New("error")
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 204, // no content
	}, nil
}

func getClient() *mongo.Client {
	var dburi = os.Getenv("mongo")
	clientOptions := options.Client().ApplyURI(dburi)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(Handler)
}
