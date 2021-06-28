package middleware

import (
	club_err "club/model/error"
	"club/model/rest"
	"club/ws/user_ws"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *gin.Context) error
type WsHandlerFunc func(c *gin.Context) (*int, error)

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

func WsErrorHandler(handler WsHandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		var restResult rest.RestResult
		var restError rest.RestError
		var err error

		userId, err := handler(c)

		if err != nil {
			log.Printf("error: %v", err.Error())
			switch err.(type) {
			case club_err.AppError:
				restError.Message = err.Error()
				restResult.Error = &restError

				if userId != nil {
					json, _ := json.Marshal(restResult)
					m := user_ws.Message{Data: []byte(json), User: *userId}
					user_ws.H.Broadcast <- m
				}

				c.JSON(http.StatusBadRequest, restResult)
				return
			default:
				restError.Description = "Websocket connection error."
				restResult.Error = &restError

				if userId != nil {
					json, _ := json.Marshal(restResult)
					m := user_ws.Message{Data: []byte(json), User: *userId}
					user_ws.H.Broadcast <- m
				}

				c.JSON(http.StatusInternalServerError, restResult)
				return
			}
		}
	}
}
