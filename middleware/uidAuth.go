package middleware

import (
	"club/model"
	"club/pojo/rest"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	restResult rest.RestResult
	restError  rest.RestError
)

func UidAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if uid := c.Request.Header.Get("X-Request-UID"); uid != "" {
			var user *model.User

			user, err := user.FindUserByUid(uid)
			if err != nil {
				restError.Description = err.Error()
				restResult.Error = &restError
				c.JSON(http.StatusInternalServerError, restResult)
			}

			if user == nil {
				restError.Message = "無效 X-Request-UID"
				restResult.Error = &restError
				c.JSON(http.StatusUnauthorized, restResult)
			}

			c.Request.Header.Set("ID", strconv.Itoa(user.Id))
			return
		} else {
			restError.Message = "無 X-Request-UID"
			restResult.Error = &restError

			c.JSON(http.StatusUnauthorized, restResult)
			return
		}
	}
}
