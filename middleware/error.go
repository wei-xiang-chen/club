package middleware

import (
	club_err "club/pojo/error"
	"club/pojo/rest"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *gin.Context) error

func ErrorHandler(handler HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		var restResult rest.RestResult
		var restError rest.RestError
		var err error

		err = handler(c)

		if err != nil {
			log.Printf("error: %v", err.Error())
			switch err.(type) {
			case club_err.AppError:
				restError.Message = err.Error()
				restResult.Error = &restError
				c.JSON(http.StatusBadRequest, restResult)
				return
			default:
				restError.Description = err.Error()
				restResult.Error = &restError
				c.JSON(http.StatusInternalServerError, restResult)
				return
			}
		}
	}
}
