package delivery

import (
	ck "bengcall/config"
	"bengcall/features/service/domain"
	"bengcall/utils/common"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type serviceHandler struct {
	srv domain.Service
}

func New(e *echo.Echo, srv domain.Service) {
	handler := serviceHandler{srv: srv}
	e.GET("/service/:id", handler.GetSpesificService(), middleware.JWT([]byte(ck.JwtKey)))
	e.POST("/admin/service", handler.AddServiceType(), middleware.JWT([]byte(ck.JwtKey)))
	e.DELETE("/admin/service/:id", handler.DeleteService(), middleware.JWT([]byte(ck.JwtKey)))
}

func (sh *serviceHandler) GetSpesificService() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, role := common.ExtractToken(c)
		if role == 0 || role == 1 {
			var check int = 0
			for _, char := range c.Param("id") {
				if unicode.IsNumber(char) {
					check += 0
				} else {
					check = 1
				}
			}
			if check == 1 {
				return c.JSON(http.StatusInternalServerError, FailResponse("id not valid"))
			} else {
				vehicle, _ := strconv.Atoi(c.Param("id"))
				res, err := sh.srv.GetSpesific(vehicle)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
				}
				return c.JSON(http.StatusOK, SuccessResponse("Success get spesific service type", ToResponse(res, "servray")))
			}
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		}
	}
}

func (sh *serviceHandler) AddServiceType() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, role := common.ExtractToken(c)
		if role == 1 {
			var input ServiceFormat

			if err := c.Bind(&input); err != nil {
				return c.JSON(http.StatusBadRequest, FailResponse("cannot bind input"))
			}

			if strings.TrimSpace(input.ServiceName) == "" {
				return c.JSON(http.StatusBadRequest, FailResponse("there's input empty"))
			} else {
				cnv := ToDomain(input)
				res, err := sh.srv.AddService(cnv)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
				}
				return c.JSON(http.StatusCreated, SuccessResponse("Success creating new service", ToResponse(res, "service")))
			}
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		}
	}
}

func (sh *serviceHandler) DeleteService() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, role := common.ExtractToken(c)
		if role == 1 {
			var check int = 0
			for _, char := range c.Param("id") {
				if unicode.IsNumber(char) {
					check += 0
				} else {
					check = 1
				}
			}
			if check == 1 {
				return c.JSON(http.StatusInternalServerError, FailResponse("id not valid"))
			} else {
				ID, _ := strconv.Atoi(c.Param("id"))
				err := sh.srv.Delete(uint(ID))
				if err != nil {
					return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
				}
				return c.JSON(http.StatusAccepted, FailResponse("Success delete service type"))
			}
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		}
	}
}
