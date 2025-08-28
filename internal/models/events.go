package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type EventModel struct {
	Id             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string        `bson:"name" json:"name"`
	StartEvent     time.Time     `bson:"start_event" json:"start_event"`
	EndEvent       time.Time     `bson:"end_event" json:"end_event"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt      *time.Time    `bson:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time    `bson:"deleted_at,omitempty" json:"deleted_at,omitempty" go.mongodb.org/mongo-driver/bson.D`
	Description    string        `bson:"description" json:"description"`
	DonationTarget *uint         `bson:"donation_target" json:"donation_target"`
}
