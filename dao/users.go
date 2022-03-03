package dao

import (
	"context"
	"encoding/json"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

type UserInsert struct {
	ID            string         `json:"id,omitempty" bson:"_id"`
	Email         string         `json:"email" bson:"email" validate:"required,email"`
	FirstName     string         `bson:"first_name" json:"first_name" validate:"required"`
	LastName      string         `bson:"last_name" json:"last_name" validate:"required"`
	FullName      string         `bson:"full_name" json:"full_name"`
	DisplayName   string         `bson:"display_name" json:"display_name"`
	Roles         vcago.RoleList `json:"system_roles" bson:"system_roles"`
	Country       string         `bson:"country" json:"country"`
	PrivacyPolicy bool           `bson:"privacy_policy" json:"privacy_policy"`
	Confirmd      bool           `bson:"confirmed" json:"confirmed"`
	LastUpdate    string         `bson:"last_update" json:"last_update"`
	Modified      vcago.Modified `json:"modified" bson:"modified"`
}

type User vcapool.User

//NewUser creates an new User from given vcago.User
func NewUser(user *vcago.User) (r *UserInsert) {
	bytes, _ := json.Marshal(&user)
	r = new(UserInsert)
	_ = json.Unmarshal(bytes, &r)
	r.Modified = vcago.NewModified()
	return
}

func ConvertUser(user *vcago.User, modified *vcago.Modified) (r *UserInsert) {
	bytes, _ := json.Marshal(&user)
	r = new(UserInsert)
	_ = json.Unmarshal(bytes, &r)
	r.Modified = *modified
	return

}

//UseUserCollection represents the user collection
var UserCollection = Database.Collection("users").CreateIndex("email", true)

//Create handles vcago.User model that is providing by auth-service.
func (i *UserInsert) Create(ctx context.Context) (err error) {
	err = UserCollection.InsertOne(ctx, &i)
	return
}

func (i *UserInsert) Update(ctx context.Context) (err error) {
	i.Modified.Update()
	err = UserCollection.UpdateOne(ctx, bson.M{"_id": i.ID}, &i)
	return
}

//Get selects an User from database
func (i *User) Get(ctx context.Context, filter bson.M) (err error) {
	if err = UserCollection.FindOne(ctx, filter, &i); err != nil {
		return
	}
	profile := new(Profile)
	err = ProfilesCollection.FindOne(ctx, bson.M{"user_id": i.ID}, profile)
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	address := new(Address)
	err = AddressesCollection.FindOne(ctx, bson.M{"user_id": i.ID}, address)
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	userCrew := new(UserCrew)
	err = UserCrewCollection.FindOne(ctx, bson.M{"user_id": i.ID}, userCrew)
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	userActive := new(UserActive)
	err = UserActiveCollection.FindOne(ctx, bson.M{"user_id": i.ID}, userActive)
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	err = nil
	i.Profile = vcapool.Profile(*profile)
	i.Address = vcapool.Address(*address)
	i.Crew = vcapool.UserCrew(*userCrew)
	i.Active = vcapool.UserActive(*userActive)
	return
}

type UserList []vcapool.User

func (i *UserList) List(ctx context.Context) (err error) {
	pipe := vcago.NewMongoPipe()
	pipe.AddModelAt(AddressesCollection.Name, "_id", "user_id", "address")
	pipe.AddModelAt(ProfilesCollection.Name, "_id", "user_id", "profile")
	pipe.AddModelAt(UserCrewCollection.Name, "_id", "user_id", "crew")
	pipe.AddModelAt(UserActiveCollection.Name, "_id", "user_id", "active")
	err = UserCollection.Aggregate(ctx, pipe.Pipe, i)
	return
}
