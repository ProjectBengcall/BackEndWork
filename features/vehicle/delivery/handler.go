package delivery

import (
	"bengcall/config"
	"bengcall/features/vehicle/domain"
	"bengcall/utils/common"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type vehicleHandler struct {
	srv domain.Service
}

func New(e *echo.Echo, srv domain.Service) {
	handler := vehicleHandler{srv: srv}
	e.POST("/admin/vehicle", handler.AddVehicle(), middleware.JWT([]byte(config.JwtKey)))
	e.GET("/vehicle", handler.GetVehicle(), middleware.JWT([]byte(config.JwtKey)))
	e.GET("/vehicleservice", handler.GetService(), middleware.JWT([]byte(config.JwtKey)))
	e.DELETE("/admin/vehicle/:id", handler.DeleteVehicle(), middleware.JWT([]byte(config.JwtKey)))

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
			if strings.TrimSpace(input.Name_vehicle) == "" {
				return c.JSON(http.StatusBadRequest, FailResponse("input empty"))
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
		if role == 1 {
			ID, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, FailResponse("id product must integer"))
			}

			res, err := bs.srv.DeleteVehicle(uint(ID))
			log.Println("res data :", res)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return c.JSON(http.StatusBadRequest, FailResponse("not found"))
				} else if strings.Contains(err.Error(), "database") {
					return c.JSON(http.StatusBadRequest, FailResponse("not found"))
				} else {
					return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
				}
			} else {
				return c.JSON(http.StatusAccepted, SuccessDeleteResponse("Success delete vehicle"))
			}
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
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
				if strings.Contains(err.Error(), "found") {
					c.JSON(http.StatusBadRequest, FailResponse(err.Error()))
				} else {
					return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
				}
			} else {
				return c.JSON(http.StatusOK, SuccessResponse("success get all vehicle", ToResponse(res, "all")))
			}
		}
		return nil
	}
}

func (bs *vehicleHandler) GetService() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, _ := common.ExtractToken(c)
		if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			ver, ser, err := bs.srv.GetService()
			if err != nil {
				if strings.Contains(err.Error(), "found") {
					c.JSON(http.StatusBadRequest, FailResponse(err.Error()))
				} else {
					return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
				}
			} else {
				return c.JSON(http.StatusOK, SuccessResponse("success get all vehicle + service", ToResponseGet(ver, ser, "vs")))
			}
		}
		return nil
	}
}
