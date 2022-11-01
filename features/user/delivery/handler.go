package delivery

import (
	ck "bengcall/config"
	"bengcall/features/user/domain"
	"bengcall/utils/common"
	"bengcall/utils/helper"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
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

func (us *userHandler) UpdateProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input EditFormat
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, FailResponse("cannot bind data"))
		}

		file, _ := c.FormFile("images")
		if file != nil {
			res, err := helper.UploadProfile(c)
			if err != nil {
				return err
			}
			log.Print(res)
			input.Images = res
		}

		id, _ := common.ExtractToken(c)
		userID := uint(id)
		cnv := ToDomain(input)
		res, err := us.srv.UpdateProfile(cnv, userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
		}

		return c.JSON(http.StatusCreated, SuccessResponse("Success update user", ToResponse(res, "edit")))
	}
}

func (uh *userHandler) Deactivate() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, role := common.ExtractToken(c)
		if role != 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			err := uh.srv.Deactivate(userID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}
			return c.JSON(http.StatusAccepted, FailResponse("success deactivate account"))
		}
	}
}

func (uh *userHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input UserFormat
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, FailResponse("cannot bind input"))
		}

		input.Images = "https://ecommerce-alta.s3.ap-southeast-1.amazonaws.com/profile/KJeT8FtTYYFq9MRbiv3u-profile.jpg"
		input.Role = 0
		cnv := ToDomain(input)
		res, err := uh.srv.Register(cnv)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
		}
		return c.JSON(http.StatusCreated, SuccessResponse("Success create new user", ToResponse(res, "reg")))
	}
}

func (uh *userHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input LoginFormat
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, FailResponse("cannot bind input"))
		}

		cnv := ToDomain(input)
		res, err := uh.srv.Login(cnv)
		fmt.Println(res.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
		}

		res.Token = common.GenerateToken(uint(res.ID), res.Role)

		return c.JSON(http.StatusAccepted, SuccessLogin("Success to login", ToResponse(res, "login")))
	}
}
