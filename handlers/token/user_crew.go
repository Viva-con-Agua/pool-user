package token

import (
	"pool-user/dao"
	"pool-user/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type UserCrewHandler struct {
	vcago.Handler
}

var UserCrew = &UserCrewHandler{*vcago.NewHandler("user_crew")}

func (i *UserCrewHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, vcapool.AccessCookieConfig())
	group.PUT("", i.Update, vcapool.AccessCookieConfig())
	group.DELETE("", i.Delete, vcapool.AccessCookieConfig())
}

func (i *UserCrewHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserCrewCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	crew := new(models.Crew)
	if err = dao.CrewsCollection.FindOne(c.Ctx(), body.CrewFilter(), crew); err != nil {
		return
	}
	result := models.NewUserCrew(token.ID, crew.ID, crew.Name, crew.Email)
	if err = dao.UserCrewCollection.InsertOne(c.Ctx(), result); err != nil {
		return
	}
	if err = dao.ActiveCollection.InsertOne(c.Ctx(), models.NewActive(token.ID, crew.ID)); err != nil {
		return
	}
	if err = dao.NVMCollection.InsertOne(c.Ctx(), models.NewNVM(token.ID)); err != nil {
		return
	}
	return c.Created(result)
}

func (i *UserCrewHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserCrewUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if token.ID != body.UserID {
		return vcago.NewPermissionDenied("crew")
	}
	result := new(models.UserCrew)
	if err = dao.UserCrewCollection.UpdateOne(c.Ctx(), body.Filter(token), vmdb.NewUpdateSet(body), result); err != nil {
		return
	}
	//reset active and nvm
	if err = dao.ActiveCollection.UpdateOne(
		c.Ctx(),
		bson.D{{Key: "user_id", Value: body.UserID}},
		vmdb.NewUpdateSet(models.ActiveWithdraw()),
		nil,
	); err != nil && vmdb.ErrNoDocuments(err) {
		return
	}
	//reject nvm state
	if err = dao.NVMCollection.UpdateOne(
		c.Ctx(),
		bson.D{{Key: "user_id", Value: body.UserID}},
		vmdb.NewUpdateSet(models.NVMWithdraw()),
		nil,
	); err != nil && vmdb.ErrNoDocuments(err) {
		return
	}
	return c.Updated(result)
}

func (i *UserCrewHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.UserCrewCollection.DeleteOne(c.Ctx(), bson.D{{Key: "user_id", Value: token.ID}}); err != nil {
		return
	}
	//reset active and nvm
	if err = dao.ActiveCollection.DeleteOne(
		c.Ctx(),
		bson.D{{Key: "user_id", Value: token.ID}},
	); err != nil && vmdb.ErrNoDocuments(err) {
		return
	}
	//reject nvm state
	if err = dao.NVMCollection.DeleteOne(
		c.Ctx(),
		bson.D{{Key: "user_id", Value: token.ID}},
	); err != nil && vmdb.ErrNoDocuments(err) {
		return
	}
	if err = dao.PoolRoleCollection.DeleteMany(
		c.Ctx(),
		bson.D{{Key: "user_id", Value: token.ID}},
	); err != nil && vmdb.ErrNoDocuments(err) {
		return
	}
	err = nil
	return c.Deleted(token.ID)

}