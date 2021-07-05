package code

import (
	"club/model/rest"
	"club/service/code_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Code(c *gin.Context) error {
	var restResult rest.RestResult

	types := c.QueryArray("type")

	codes, err := code_service.Code(types)
	if err != nil {
		return err
	}

	restResult.Data = codes
	c.JSON(http.StatusOK, restResult)
	return nil
}
