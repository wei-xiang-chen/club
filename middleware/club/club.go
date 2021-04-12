package club

import (
	"club/pojo"
	"club/pojo/rest"
	"club/service/club_service"
	rest_util "club/util/rest"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	restResult rest.RestResult
	restError  rest.RestError

	pageDefaultS  string = "1"
	limitDefaultS string = "10"
)

func GetList(c *gin.Context) {
	var page, limit int

	topic := c.Query("topic")
	clubName := c.Query("clubName")

	pageString := c.Query("page")
	limitString := c.Query("limit")
	if pageString == "" {
		pageString = pageDefaultS
	}
	if limitString == "" {
		limitString = limitDefaultS
	}

	page, err := strconv.Atoi(pageString)
	if err != nil {
		restError.Description = err.Error()
		restResult.Error = &restError
		c.JSON(http.StatusInternalServerError, restResult)
		return
	}

	limit, err = strconv.Atoi(limitString)
	if err != nil {
		restError.Description = err.Error()
		restResult.Error = &restError
		c.JSON(http.StatusInternalServerError, restResult)
		return
	}

	offset := (page - 1) * limit
	clubs, err := club_service.GetList(topic, clubName, offset, limit)
	if err != nil {
		restError.Description = err.Error()
		restResult.Error = &restError
		c.JSON(http.StatusInternalServerError, restResult)
		return
	}

	restResult.Data = clubs
	c.JSON(http.StatusOK, restResult)
}

func Create(c *gin.Context) {
	var club pojo.Club

	c.ShouldBindJSON(&club)

	user, err := rest_util.GetUser(c)
	if err != nil {
		restError.Description = err.Error()
		restResult.Error = &restError
		c.JSON(http.StatusInternalServerError, restResult)
		return
	}

	club.Owner = user.Id
	err = club_service.Insert(&club)
	if err != nil {
		restError.Description = err.Error()
		restResult.Error = &restError
		c.JSON(http.StatusInternalServerError, restResult)
		return
	}

	restResult.Data = club
	c.JSON(http.StatusOK, restResult)
}

func Join(c *gin.Context) {
	clubIdString := c.Param("clubId")
	clubId, err := strconv.Atoi(clubIdString)

	if err != nil {
		restError.Description = err.Error()
		restResult.Error = &restError
		c.JSON(http.StatusInternalServerError, restResult)
		return
	}

	user, err := rest_util.GetUser(c)
	if err != nil {
		restError.Description = err.Error()
		restResult.Error = &restError
		c.JSON(http.StatusInternalServerError, restResult)
		return
	}

	err = club_service.Join(user.Id, clubId)
	if err != nil {
		restError.Description = err.Error()
		restResult.Error = &restError
		c.JSON(http.StatusInternalServerError, restResult)
		return
	}

	c.Status(200)
}

func Leave(c *gin.Context) {
	user, err := rest_util.GetUser(c)
	if err != nil {
		restError.Description = err.Error()
		restResult.Error = &restError
		c.JSON(http.StatusInternalServerError, restResult)
		return
	}

	err = club_service.Leave(user.Id)
	if err != nil {
		restError.Description = err.Error()
		restResult.Error = &restError
		c.JSON(http.StatusInternalServerError, restResult)
		return
	}

	c.Status(200)
}
