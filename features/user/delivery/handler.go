package delivery

import (
	ck "bengcall/config"
	"bengcall/features/user/domain"
	"bengcall/utils/common"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type userHandler struct {
	srv domain.Service
}

func New(e *echo.Echo, srv domain.Service) {
	handler := userHandler{srv: srv}
	e.GET("/users", handler.MyProfile(), middleware.JWT([]byte(ck.JwtKey)))
	e.PUT("/users", handler.UpdateProfile(), middleware.JWT([]byte(ck.JwtKey)))
	e.DELETE("/users", handler.Deactivate(), middleware.JWT([]byte(ck.JwtKey)))
	e.POST("/register", handler.Register())
	e.POST("/login", handler.Login())
}

func (uh *userHandler) MyProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, role := common.ExtractToken(c)
		if role != 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("this account costumer"))
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			res, err := uh.srv.MyProfile(uint(userID))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}
			return c.JSON(http.StatusOK, SuccessResponse("success get my profile", ToResponse(res, "user")))
		}
	}
}

func (uh *userHandler) UpdateProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, role := common.ExtractToken(c)
		if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else if role == 0 {
			var input EditFormat
			if err := c.Bind(&input); err != nil {
				return c.JSON(http.StatusBadRequest, FailResponse("cannot bind update data"))
			}

			file, fileheader, _ := c.Request().FormFile("images")

			if strings.TrimSpace(input.Email) == "" && strings.TrimSpace(input.Password) == "" && strings.TrimSpace(input.Fullname) == "" && file == nil {
				return c.JSON(http.StatusBadRequest, FailResponse("please insert one field"))
			}

			cnv := ToDomain(input)
			res, err := uh.srv.UpdateProfile(cnv, file, fileheader, userID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}
			return c.JSON(http.StatusAccepted, SuccessResponse("success update user", ToResponse(res, "edit")))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("this account costumer"))
		}
	}
}

func (uh *userHandler) Deactivate() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, role := common.ExtractToken(c)

		if role == 0 {
			res, err := uh.srv.Deactivate(userID)
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
				return c.JSON(http.StatusAccepted, SuccessDeleteResponse("Success deactivate account"))
			}
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("this account costumer"))
		}

	}
}

func (uh *userHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input UserFormat

		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, FailResponse("cannot bind input"))
		}
		// if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
		// 	return c.JSON(http.StatusBadRequest, FailResponse("input empty"))
		// }

		cnv := ToDomain(input)
		res, err := uh.srv.Register(cnv)
		if err != nil {
			if strings.Contains(err.Error(), "password") {
				c.JSON(http.StatusBadRequest, FailResponse(err.Error()))
			} else if strings.Contains(err.Error(), " email") {
				c.JSON(http.StatusBadRequest, FailResponse(err.Error()))
			} else if strings.Contains(err.Error(), "already") {
				c.JSON(http.StatusBadRequest, FailResponse(err.Error()))
			} else {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}
		} else if res.ID != 0 {
			return c.JSON(http.StatusCreated, SuccessResponse("Success create new user", ToResponse(res, "reg")))
		}
		return nil
	}
}

func (uh *userHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {

		var input LoginFormat
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, FailResponse("cannot bind input"))
		}
		if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
			return c.JSON(http.StatusBadRequest, FailResponse("password or email empty"))
		}
		cnv := ToDomain(input)
		res, err := uh.srv.Login(cnv)
		fmt.Println(res.ID)
		if err != nil {
			if strings.Contains(err.Error(), "found") {
				c.JSON(http.StatusBadRequest, FailResponse(err.Error()))
			} else if strings.Contains(err.Error(), "wrong") {
				c.JSON(http.StatusBadRequest, FailResponse(err.Error()))
			} else {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}
		} else if res.ID != 0 {
			res.Token = common.GenerateToken(uint(res.ID), res.Role)
			return c.JSON(http.StatusAccepted, SuccessLogin("Success to login", ToResponse(res, "login")))

		}
		return nil
	}
}
