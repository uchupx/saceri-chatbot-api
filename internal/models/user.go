package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserModel struct {
	Id        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string        `bson:"username" json:"username"`
	Password  string        `bson:"password" json:"-"` // "-" excludes from JSON response
	Name      string        `bson:"name" json:"name"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
	DeletedAt time.Time     `bson:"deleted_at,omitempty" json:"deleted_at,omitempty" go.mongodb.org/mongo-driver/bson.D`
	OauthID   string        `bson:"oauth_id,omitempty" json:"oauth_id"`
	IsActive  bool          `bson:"is_active" json:"is_active"`
}
