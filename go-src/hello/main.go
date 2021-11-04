package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// name := request.QueryStringParameters["name"]

	col := getClient().Database("test").Collection("test-col")

	insertResult, err := col.InsertOne(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("error")
	}

	response := fmt.Sprintf("id : %s!", insertResult.InsertedID.(primitive.ObjectID).Hex())

	return &events.APIGatewayProxyResponse{
		StatusCode:        200,
		MultiValueHeaders: http.Header{"Set-Cookie": {"Ding", "Ping"}},
		Headers:           map[string]string{"Content-Type": "text/html; charset=UTF-8"},
		Body:              response,
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
