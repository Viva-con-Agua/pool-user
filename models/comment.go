package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	ActivityDB struct {
		ID        string `json:"id" bson:"_id"`
		UserID    string `json:"user_id" bson:"user_id"`
		Comment   string `json:"comment" bson:"comment"`
		ModelType string `json:"model_type" bson:"model_type"`
		ModelID   string `json:"model_id" bson:"model_id"`
		Status    string `json:"status" bson:"status"`
		Created   int64  `json:"created" bson:"created"`
	}
	Activity struct {
		ID        string       `json:"id" bson:"_id"`
		UserID    string       `json:"user_id" bson:"user_id"`
		User      UserDatabase `json:"user" bson:"user"`
		Comment   string       `json:"comment" bson:"comment"`
		ModelType string       `json:"model_type" bson:"model_type"`
		ModelID   string       `json:"model_id" bson:"model_id"`
		Status    string       `json:"status" bson:"status"`
		Created   int64        `json:"created" bson:"created"`
	}
)

func NewActivityDB(userID string, modelType string, modelID string, comment string, status string) *ActivityDB {
	return &ActivityDB{
		ID:        uuid.NewString(),
		UserID:    userID,
		Comment:   comment,
		ModelType: modelType,
		ModelID:   modelID,
		Status:    status,
		Created:   time.Now().Unix(),
	}
}
