package models

import (
	"time"

	"github.com/zibakasebi/api/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollectionTime struct {
	ID primitive.ObjectID

	Name      string
	StartDate time.Time
	EndDate   time.Time
	Ranges    []Range
	WeekDays  []enums.WeekDay
}

type Range struct {
	StartTime time.Time
	EndTime   time.Time
}
