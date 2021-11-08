package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	ptime "github.com/yaa110/go-persian-calendar"
)

type AppTime struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Name      string         `bson:"name"`
	StartDate time.Time      `bson:"start_date"`
	EndDate   time.Time      `bson:"end_date"`
	Ranges    []Range        `bson:"ranges"`
	WeekDays  []time.Weekday `bson:"week_days"`
}

type Range struct {
	StartTime time.Time `bson:"start_time"`
	EndTime   time.Time `bson:"end_time"`
}

type Reserve struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Date          time.Time `bson:"date"`
	StartTime     time.Time `bson:"start_time"`
	EndTime       time.Time `bson:"end_time"`
	MemberName    string    `bson:"member_name"`
	MemberMobile  string    `bson:"member_mobile"`
	MemberMessage string    `bson:"member_message"`
	MemberEmail   string    `bson:"member_email"`
	Price         int       `bson:"price"`
	RegisterDate  time.Time `bson:"register_date"`
	Done          bool      `bson:"done"`
}

type ExceptionDate struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Date    time.Time `bson:"date"`
	Ranges  []Range   `bson:"ranges"`
	Holiday bool      `bson:"holiday"`
}

type ShowDate struct {
	Date         time.Time    `json:"date"`
	Active       bool         `json:"active"`
	Ranges       []Range      `json:"ranges"`
	Reserves     []Reserve    `json:"reserves"`
	WeekDay      time.Weekday `json:"week_day"`
	PersianMonth string       `json:"persian_month"`
}

func InternalServerError() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       "",
	}, nil
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	database := getClient().Database("test")

	colAppTime := database.Collection("app_time")
	var appTime AppTime
	err := colAppTime.FindOne(ctx, bson.M{}).Decode(&appTime)
	if err != nil {
		return InternalServerError()
	}

	colExceptionDate := database.Collection("exception_date")
	var exceptionDates []ExceptionDate
	cur, err := colExceptionDate.Find(ctx, bson.M{
		"date": bson.M{
			"$gte": appTime.StartDate,
			"$lte": appTime.EndDate,
		},
	})
	if err != nil {
		return InternalServerError()
	}

	err = cur.All(ctx, &exceptionDates)
	if err != nil {
		return InternalServerError()
	}

	colReserve := database.Collection("reserve")
	var reserves []Reserve
	cur, err = colReserve.Find(ctx, bson.M{
		"date": bson.M{
			"$gte": appTime.StartDate,
			"$lte": appTime.EndDate,
		},
	})
	if err != nil {
		return InternalServerError()
	}

	err = cur.All(ctx, &reserves)
	if err != nil {
		return InternalServerError()
	}

	//loop on dayes between two date
	var days = int(appTime.EndDate.Sub(appTime.StartDate).Hours() / 24)
	var showDates []ShowDate
	for i := 0; i <= days; i++ {
		var date = appTime.StartDate.AddDate(0, 0, i)
		var showDate *ShowDate = nil
		for _, item := range exceptionDates {
			if item.Date.Equal(date) {
				showDate = &ShowDate{
					Date:    item.Date,
					Active:  !item.Holiday,
					Ranges:  item.Ranges,
					WeekDay: item.Date.Weekday(),
				}
				break
			}
		}
		if showDate == nil {
			showDate = &ShowDate{
				Date:    date,
				Active:  true,
				Ranges:  appTime.Ranges,
				WeekDay: date.Weekday(),
			}
		}

		for _, item := range reserves {
			if item.Date.Equal(date) {
				showDate.Reserves = append(showDate.Reserves, item)
			}
		}

		showDate.PersianMonth = getPersianMonth(showDate.Date)
		showDates = append(showDates, *showDate)
	}

	//sort by date
	sort.Slice(showDates, func(i, j int) bool {
		return showDates[i].Date.Before(showDates[j].Date)
	})

	//convert showDates to json
	json, err := json.Marshal(showDates)
	if err != nil {
		return InternalServerError()
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "json"},
		Body:       string(json),
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

func HandlerTest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Hello World",
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(Handler)
	lambda.Start(HandlerTest)
}

//get persian month from gregoian date
func getPersianMonth(date time.Time) string {
	pt := ptime.New(date)
	month := int(pt.Month())
	return monthToString(month)
}

func monthToString(month int) string {
	if month == 1 {
		return "فروردین"
	} else if month == 2 {
		return "اردیبهشت"
	} else if month == 3 {
		return "خرداد"
	} else if month == 4 {
		return "تیر"
	} else if month == 5 {
		return "مرداد"
	} else if month == 6 {
		return "شهریور"
	} else if month == 7 {
		return "مهر"
	} else if month == 8 {
		return "آبان"
	} else if month == 9 {
		return "آذر"
	} else if month == 10 {
		return "دی"
	} else if month == 11 {
		return "بهمن"
	} else if month == 12 {
		return "اسفند"
	} else {
		return "خطا"
	}
}
