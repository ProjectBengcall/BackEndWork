package delivery

import (
	"bengcall/config"
	"bengcall/features/vehicle/domain"
	"bengcall/utils/common"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type vehicleHandler struct {
	srv domain.Service
}

func New(e *echo.Echo, srv domain.Service) {
	handler := vehicleHandler{srv: srv}
	e.POST("/vehicles", handler.AddVehicle(), middleware.JWT([]byte(config.JwtKey)))
	e.GET("/books", handler.GetVehicle())
	e.DELETE("/vehicles/:id", handler.DeleteVehicle(), middleware.JWT([]byte(config.JwtKey)))

}

func (vh *vehicleHandler) AddVehicle() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input AddFormat
		userID, role := common.ExtractToken(c)
		if role != 1 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			if err := c.Bind(&input); err != nil {
				return c.JSON(http.StatusBadRequest, FailResponse("cannot bind input"))
			}
			cnv := ToDomain(input)
			res, err := vh.srv.AddVehicle(cnv)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}

			return c.JSON(http.StatusCreated, SuccessResponse("sucses add vehicle", ToResponse(res, "reg")))
		}
	}
}

func (bs *vehicleHandler) DeleteVehicle() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, role := common.ExtractToken(c)
		if role != 1 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			ID, err := strconv.Atoi(c.Param("id"))
			if err = bs.srv.DeleteVehicle(uint(ID)); err != nil {
				return c.JSON(http.StatusBadRequest, FailResponse(err.Error()))
			}
			return c.JSON(http.StatusAccepted, FailResponse("success delete vehicle"))
		}
	}
}

func (bs *vehicleHandler) GetVehicle() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, _ := common.ExtractToken(c)
		if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			res, err := bs.srv.GetVehicle()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}

			return c.JSON(http.StatusOK, SuccessResponse("success get all vehicle", ToResponse(res, "all")))
		}
	}
}
