package rest_util

import (
	"club/pojo"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	user pojo.User
)

func GetUser(c *gin.Context) (pojo.User, error) {
	user.Uid = c.Request.Header.Get("X-Request-UID")
	i, err := strconv.Atoi(c.Request.Header.Get("ID"))

	if err != nil {
		return user, err
	}

	user.Id = i

	return user, nil
}
