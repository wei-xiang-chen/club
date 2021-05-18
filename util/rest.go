package util

import (
	"club/pojo"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) (*pojo.User, error) {
	var user pojo.User

	user.Uid = c.Request.Header.Get("X-Request-UID")
	i, err := strconv.Atoi(c.Request.Header.Get("ID"))

	if err != nil {
		return nil, err
	}

	user.Id = i

	return &user, nil
}
