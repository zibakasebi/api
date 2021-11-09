package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Setting struct {
	AppTime_StartDate time.Time      `bson:"app_time_start_date"`
	AppTime_EndDate   time.Time      `bson:"app_time_end_date"`
	AppTime_Ranges    []Range        `bson:"app_time_ranges"`
	AppTime_WeekDays  []time.Weekday `bson:"app_time_week_days"`

	Kavenegar string `bson:"kavenegar"`
	ZarinPal  string `bson:"zarinPal"`
	Price     int    `bson:"price"`
}

type Range struct {
	StartTime string `bson:"start_time"`
	EndTime   string `bson:"end_time"`
}

// type Reserve struct {
// 	ID primitive.ObjectID `bson:"_id,omitempty"`

// 	Date          time.Time `bson:"date"`
// 	StartTime     time.Time `bson:"start_time"`
// 	EndTime       time.Time `bson:"end_time"`
// 	MemberName    string    `bson:"member_name"`
// 	MemberMobile  string    `bson:"member_mobile"`
// 	MemberMessage string    `bson:"member_message"`
// 	MemberEmail   string    `bson:"member_email"`
// 	Price         int       `bson:"price"`
// 	RegisterDate  time.Time `bson:"register_date"`
// 	Done          bool      `bson:"done"`
// }

// type ExceptionDate struct {
// 	ID primitive.ObjectID `bson:"_id,omitempty"`

// 	Date    time.Time `bson:"date"`
// 	Ranges  []Range   `bson:"ranges"`
// 	Holiday bool      `bson:"holiday"`
// }

// type ShowDate struct {
// 	Date         time.Time    `json:"date"`
// 	Active       bool         `json:"active"`
// 	Ranges       []Range      `json:"ranges"`
// 	Reserves     []Reserve    `json:"reserves"`
// 	WeekDay      time.Weekday `json:"week_day"`
// 	PersianMonth string       `json:"persian_month"`
// }

func InternalServerError(message string) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       message,
	}, nil
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

	if request.HTTPMethod == "GET" {
		return HandlerGet(ctx, request)
	} else if request.HTTPMethod == "PUT" {
		return HandlerEdit(ctx, request)
	} else {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotImplemented,
		}, nil
	}
}

func HandlerGet(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	database := getClient().Database("test")

	colSetting := database.Collection("setting")
	var setting Setting

	err := colSetting.FindOne(ctx, bson.M{}).Decode(&setting)
	if err != nil && err != mongo.ErrNoDocuments {
		return InternalServerError(err.Error())
	}

	if err == mongo.ErrNoDocuments {
		setting.AppTime_StartDate = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		setting.AppTime_EndDate = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		setting.AppTime_Ranges = []Range{}
		setting.AppTime_WeekDays = []time.Weekday{0, 1, 2, 3, 4, 5, 6}
		setting.Kavenegar = ""
		setting.ZarinPal = ""
		setting.Price = 0
	}

	//convert struct to json string
	jsonSetting, err := json.Marshal(setting)
	if err != nil {
		return InternalServerError(err.Error())
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
		Headers:    map[string]string{"Content-Type": "json"},
		Body:       string(jsonSetting),
	}, nil
}

func HandlerEdit(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Print("body:", request.Body)

	database := getClient().Database("test")

	colSetting := database.Collection("setting")
	var setting Setting
	//convert json string to struct
	err := json.Unmarshal([]byte(request.Body), &setting)
	if err != nil {
		log.Print(err)
		return InternalServerError(err.Error())
	}

	replaceOption := options.Replace()
	replaceOption.SetUpsert(true)

	_, err = colSetting.ReplaceOne(ctx, bson.M{}, setting, replaceOption)
	if err != nil {
		log.Print(err)
		return InternalServerError(err.Error())
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
		Headers:    map[string]string{"Content-Type": "json"},
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
