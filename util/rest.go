package util

import (
	"club/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) (*model.User, error) {
	var user model.User

	user.Uid = c.Request.Header.Get("X-Request-UID")
	i, err := strconv.Atoi(c.Request.Header.Get("ID"))

	if err != nil {
		return nil, err
	}

	user.Id = i

	return &user, nil
}
