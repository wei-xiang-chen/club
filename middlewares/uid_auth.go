package middlewares

import (
	"club/dal"
	"club/models/rest"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UidAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var restResult rest.RestResult
		var restError rest.RestError

		if uid := c.Request.Header.Get("X-Request-UID"); uid != "" {
			var user *dal.User

			user, err := user.FindUserByUid(&uid)
			if err != nil {
				restError.Description = err.Error()
				restResult.Error = &restError
				c.JSON(http.StatusInternalServerError, restResult)
				c.Abort()
				return
			}

			if user == nil {
				restError.Message = "無效 X-Request-UID"
				restResult.Error = &restError
				c.JSON(http.StatusUnauthorized, restResult)
				c.Abort()
				return
			}

			c.Request.Header.Set("ID", strconv.Itoa(user.Id))
			return
		} else {
			restError.Message = "無 X-Request-UID"
			restResult.Error = &restError

			c.JSON(http.StatusUnauthorized, restResult)
			c.Abort()
			return
		}
	}
}
