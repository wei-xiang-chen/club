package login

import (
	"club/model"
	"club/model/rest"
	"club/service/user_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) error {
	var restResult rest.RestResult
	var login model.User

	c.ShouldBindJSON(&login)

	user, err := user_service.Insert(&login)
	if err != nil {
		return nil
	}

	restResult.Data = user
	c.JSON(http.StatusOK, restResult)
	return nil
}
