package delivery

import (
	"bengcall/features/vehicle/domain"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type vehicleHandler struct {
	srv domain.Service
}

func New(e *echo.Echo, srv domain.Service) {
	handler := vehicleHandler{srv: srv}
	e.POST("/vehicles", handler.AddVehicle())
	// e.GET("/books", handler.ShowAllBook())
	e.DELETE("/vehicles/:id", handler.DeleteVehicle())

}

func (vh *vehicleHandler) AddVehicle() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input AddFormat
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

func (bs *vehicleHandler) DeleteVehicle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ID, err := strconv.Atoi(c.Param("id"))
		if err = bs.srv.DeleteVehicle(uint(ID)); err != nil {
			return c.JSON(http.StatusBadRequest, FailResponse(err.Error()))
		}
		return c.JSON(http.StatusAccepted, FailResponse("success delete vehicle"))
	}
}

func (bs *vehicleHandler) GetVehicle() echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := bs.srv.GetVehicle()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
		}

		return c.JSON(http.StatusOK, SuccessResponse("success get all user", ToResponse(res, "all")))
	}
}
