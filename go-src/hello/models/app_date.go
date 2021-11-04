package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppDate struct {
	ID primitive.ObjectID

	Date   time.Time
	Enable bool
	Ranges []Range
}
