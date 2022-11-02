package main

import (
	"bengcall/config"
	dUser "bengcall/features/user/delivery"
	rUser "bengcall/features/user/repository"
	sUser "bengcall/features/user/services"
	"bengcall/utils/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	cfg := config.NewConfig()
	db := database.InitDB(cfg)

	mdlUser := rUser.New(db)
	serUser := sUser.New(mdlUser)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	dUser.New(e, serUser)
	e.Logger.Fatal(e.Start(":8000"))

}
