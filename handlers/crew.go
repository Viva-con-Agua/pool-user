package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

func CreateCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Crew)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	/*user := new(vcapool.User)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}*/
	if err = body.Create(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("address", body).Created())
}

func GetCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	result := new(dao.Crew)
	if err = result.Get(ctx, c.Param("id")); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("address", result).Selected())
}

func UpdateCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Crew)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	if err = body.Update(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("address", body).Updated())
}

func DeleteCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Crew)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	if err = body.Delete(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("address", body).Deleted())
}

func ListCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.AddressQuery)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	result := new(dao.AddressList)
	if err = result.Get(ctx, body.Filter()); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("address_list", result).Selected())
}
