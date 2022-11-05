package main

import (
	"bengcall/config"
	dUser "bengcall/features/user/delivery"
	rUser "bengcall/features/user/repository"
	sUser "bengcall/features/user/services"
	dVehicle "bengcall/features/vehicle/delivery"
	rVehicle "bengcall/features/vehicle/repository"
	sVehicle "bengcall/features/vehicle/services"
	"bengcall/utils/database"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	cfg := config.NewConfig()
	db := database.InitDB(cfg)
	validator := validator.New()

	mdlUser := rUser.New(db)
	mdlVehicle := rVehicle.New(db)
	serUser := sUser.New(mdlUser, validator)
	serVehicle := sVehicle.New(mdlVehicle)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	dUser.New(e, serUser)
	dVehicle.New(e, serVehicle)

	e.Logger.Fatal(e.Start(":8000"))

}
