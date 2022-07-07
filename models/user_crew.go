package models

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	UserCrewCreate struct {
		CrewID string `json:"crew_id"`
	}
	UserCrew struct {
		ID       string         `bson:"_id" json:"id"`
		UserID   string         `bson:"user_id" json:"user_id"`
		Name     string         `bson:"name" json:"name"`
		Email    string         `bson:"email" json:"email"`
		Roles    []vcago.Role   `bson:"roles" json:"roles"`
		CrewID   string         `bson:"crew_id" json:"crew_id"`
		Modified vcago.Modified `bson:"modified" json:"modified"`
	}
	UserCrewUpdate struct {
		ID     string `bson:"_id" json:"id"`
		UserID string `bson:"user_id" json:"user_id"`
		Name   string `bson:"name" json:"name"`
		Email  string `bson:"email" json:"email"`
		CrewID string `bson:"crew_id" json:"crew_id"`
	}
	UserCrewParam struct {
		ID string `param:"id"`
	}
)

func NewUserCrew(userID string, crewID string, name string, email string) *UserCrew {
	return &UserCrew{
		ID:       uuid.NewString(),
		UserID:   userID,
		Name:     name,
		Email:    email,
		CrewID:   crewID,
		Modified: vcago.NewModified(),
	}
}

func (i *UserCrewCreate) CrewFilter() bson.D {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.CrewID)
	return bson.D(*match)
}

func (i *UserCrewUpdate) Filter(token *vcapool.AccessToken) bson.D {
	return bson.D{{Key: "_id", Value: i.ID}, {Key: "user_id", Value: token.ID}}
}