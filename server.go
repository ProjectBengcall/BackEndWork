package main

import (
	"bengcall/config"
	dService "bengcall/features/service/delivery"
	rService "bengcall/features/service/repository"
	sService "bengcall/features/service/services"
	dTrx "bengcall/features/transaction/delivery"
	rTrx "bengcall/features/transaction/repository"
	sTrx "bengcall/features/transaction/services"
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

	mdlService := rService.New(db)
	serService := sService.New(mdlService)

	mdlTrx := rTrx.New(db)
	serTrx := sTrx.New(mdlTrx)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	dService.New(e, serService)
	dUser.New(e, serUser)
	dVehicle.New(e, serVehicle)
	dTrx.New(e, serTrx)

	e.Logger.Fatal(e.Start(":8000"))

}
