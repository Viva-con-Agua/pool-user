package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Crew struct {
	ID       string         `json:"id,omitempty" bson:"_id"`
	Name     string         `json:"name" bson:"name"`
	City     string         `json:"city" bson:"city"`
	Country  string         `json:"country" bson:"country"`
	Modified vcago.Modified `json:"modified" bson:"modified"`
}

var CrewsCollection = Database.Collection("crews").CreateIndex("name", true)

func (i *Crew) Create(ctx context.Context) (err error) {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	i.Modified = vcago.NewModified()
	err = UserCollection.InsertOne(ctx, &i)
	return
}

func (i *Crew) Get(ctx context.Context, id string) (err error) {
	err = CrewsCollection.FindOne(ctx, bson.M{"_id": id}, &id)
	return
}

func (i *Crew) Update(ctx context.Context) (err error) {
	i.Modified.Update()
	err = CrewsCollection.UpdateOne(ctx, bson.M{"_id": i.ID}, &i)
	return
}

func (i *Crew) Delete(ctx context.Context) (err error) {
	err = CrewsCollection.DeleteOne(ctx, bson.M{"_id": i.ID})
	return
}

type CrewQuery struct {
	ID   string `query:"id"`
	Name string `query:"name"`
}

func (i *CrewQuery) Filter() bson.M {
	f := vcago.NewMongoFilterM()
	f.Equal("_id", i.ID)
	f.Like("name", i.Name)
	return f.Filter
}

type CrewList []Crew

func (i *CrewList) Get(ctx context.Context, filter bson.M) (err error) {
	err = CrewsCollection.Find(ctx, filter, &i)
	return
}
