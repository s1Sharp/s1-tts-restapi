package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserIdGetTasksRequest primitive.ObjectID

type UserResponse struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	Email         string             `json:"email,omitempty" bson:"email,omitempty"`
	Role          string             `json:"role,omitempty" bson:"role,omitempty"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
	LastAttemptAt time.Time          `bson:"last_attempt_at"`
}
