package login

import (
	"club/pojo"
	"club/pojo/rest"
	"club/service/user_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	restResult rest.RestResult
	restError  rest.RestError
)

func Login(c *gin.Context) {
	var login pojo.User

	c.ShouldBindJSON(&login)

	user, err := user_service.Insert(&login)
	if err != nil {
		restError.Description = err.Error()
		restResult.Error = &restError
		c.JSON(http.StatusInternalServerError, restResult)
		return
	}

	restResult.Data = user
	c.JSON(http.StatusOK, restResult)
}
