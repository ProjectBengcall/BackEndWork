package delivery

import (
	ck "bengcall/config"
	"bengcall/features/transaction/domain"
	"bengcall/utils/common"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type transactionHandler struct {
	srv domain.Service
}

func New(e *echo.Echo, srv domain.Service) {
	handler := transactionHandler{srv: srv}
	e.GET("/history", handler.HistoryTransaction(), middleware.JWT([]byte(ck.JwtKey)))
	e.GET("/transaction/me", handler.MyTransaction(), middleware.JWT([]byte(ck.JwtKey)))
	e.GET("/transaction/:id", handler.DetailTransaction(), middleware.JWT([]byte(ck.JwtKey)))
	e.GET("/admin/transaction", handler.AllTransaction(), middleware.JWT([]byte(ck.JwtKey)))
	e.POST("/transaction", handler.NewTransaction(), middleware.JWT([]byte(ck.JwtKey)))
	e.POST("/transaction/success", handler.TransactionSuccess())
	e.PUT("/comment/:id", handler.AddComment(), middleware.JWT([]byte(ck.JwtKey)))
	e.PUT("/admin/transaction/:id", handler.UpdateStatus(), middleware.JWT([]byte(ck.JwtKey)))
	e.DELETE("/admin/transaction/:id", handler.CancelTransaction(), middleware.JWT([]byte(ck.JwtKey)))
}

func (th *transactionHandler) NewTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, role := common.ExtractToken(c)
		if role == 0 || role == 1 {
			var input TransactionFormat
			if err := c.Bind(&input); err != nil {
				return c.JSON(http.StatusBadRequest, FailResponse("cannot bind input"))
			}

			cnv := ToDomain(input)
			cns := ToDom(input)
			cnv.UserID = userID
			res, err := th.srv.Transaction(cnv, cns)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}
			return c.JSON(http.StatusCreated, SuccessResponse("Success creating new service", ToResponse(res, "post")))
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		}
	}
}

func (th *transactionHandler) TransactionSuccess() echo.HandlerFunc {
	return func(c echo.Context) error {
		var trx SuccessFormat
		if err := c.Bind(&trx); err != nil {
			return c.JSON(http.StatusBadRequest, FailResponse("cannot bind input"))
		}

		cnv := ToSucc(trx)
		fmt.Println(cnv.Status)
		fmt.Println(cnv.Order)

		if cnv.Status == "capture" || cnv.Status == "settlement" {
			ID, _ := strconv.Atoi(cnv.Order)
			err := th.srv.Success(uint(ID))
			//res := helper.Create(uint(ID)) //proses pengiriman invoice ke email
			//lo.Println(res)
			if err != nil {
				log.Error(err)
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}
			return c.JSON(http.StatusCreated, FailResponse("Transaction Success"))
		}
		return c.JSON(http.StatusBadRequest, FailResponse("Transaction Error"))
	}
}

func (th *transactionHandler) UpdateStatus() echo.HandlerFunc {
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
				var input StatusFormat
				if err := c.Bind(&input); err != nil {
					return c.JSON(http.StatusBadRequest, FailResponse("cannot bind data"))
				}

				cnv := ToDomain(input)
				res, err := th.srv.Status(cnv, uint(ID))
				if err != nil {
					return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
				}
				return c.JSON(http.StatusCreated, SuccessResponse("Success update transaction status", ToResponse(res, "stts")))
			}
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		}
	}
}

func (th *transactionHandler) AddComment() echo.HandlerFunc {
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
				ID, _ := strconv.Atoi(c.Param("id"))
				var input CommentFormat
				if err := c.Bind(&input); err != nil {
					return c.JSON(http.StatusBadRequest, FailResponse("cannot bind data"))
				}
				if strings.TrimSpace(input.Comment) == "" {
					return c.JSON(http.StatusBadRequest, FailResponse("input empty"))
				} else {
					cnv := ToDomain(input)
					res, err := th.srv.Comment(cnv, uint(ID))
					if err != nil {
						return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
					}

					return c.JSON(http.StatusCreated, SuccessResponse("Success add comment", ToResponse(res, "cmmt")))
				}

			}
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		}
	}
}

func (th *transactionHandler) HistoryTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, role := common.ExtractToken(c)
		if role == 0 || role == 1 {
			res, err := th.srv.History(userID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}
			return c.JSON(http.StatusOK, SuccessResponse("Success get all transaction data", ToResponse(res, "history")))
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		}
	}
}

func (th *transactionHandler) MyTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, role := common.ExtractToken(c)
		if role == 0 || role == 1 {
			res, err := th.srv.My(userID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}
			return c.JSON(http.StatusOK, SuccessResponse("Success get my transaction", ToResponse(res, "progress")))
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		}
	}
}

func (th *transactionHandler) DetailTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, role := common.ExtractToken(c)
		if role == 0 || role == 1 {
			ID, _ := strconv.Atoi(c.Param("id"))
			res, dtl, err := th.srv.Detail(uint(ID))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}
			return c.JSON(http.StatusOK, SuccessResponse("Success get detail transaction data", ToResponses(res, dtl, "ally")))
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		}
	}
}

func (th *transactionHandler) AllTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, role := common.ExtractToken(c)
		if role == 1 {
			res, err := th.srv.All()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}
			return c.JSON(http.StatusOK, SuccessResponse("Success get all transaction data", ToResponse(res, "all")))
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		}
	}
}

func (th *transactionHandler) CancelTransaction() echo.HandlerFunc {
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
				err := th.srv.Cancel(uint(ID))
				if err != nil {
					return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
				}
				return c.JSON(http.StatusAccepted, FailResponse("Success Cancel Transaction Service"))
			}
		} else if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse("cannot validate token"))
		}
	}
}
