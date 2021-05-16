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
	pageDefaultS  string = "1"
	limitDefaultS string = "10"
)

func GetList(c *gin.Context) error {
	var page, limit int
	var restResult rest.RestResult

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
		return err
	}

	limit, err = strconv.Atoi(limitString)
	if err != nil {
		return err
	}

	offset := (page - 1) * limit
	clubs, err := club_service.GetList(topic, clubName, offset, limit)
	if err != nil {
		return err
	}

	restResult.Data = clubs
	c.JSON(http.StatusOK, restResult)
	return nil
}

func Create(c *gin.Context) error {
	var club pojo.Club
	var restResult rest.RestResult

	c.ShouldBindJSON(&club)

	user, err := rest_util.GetUser(c)
	if err != nil {
		return err
	}

	club.Owner = user.Id
	err = club_service.Insert(&club)
	if err != nil {
		return err
	}

	restResult.Data = club
	c.JSON(http.StatusOK, restResult)
	return nil
}

func Join(c *gin.Context) error {

	clubIdString := c.Param("clubId")
	clubId, err := strconv.Atoi(clubIdString)

	if err != nil {
		return err
	}

	user, err := rest_util.GetUser(c)
	if err != nil {
		return err
	}

	err = club_service.Join(user.Id, clubId)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func Leave(c *gin.Context) error {

	user, err := rest_util.GetUser(c)
	if err != nil {
		return err
	}

	err = club_service.Leave(user.Id)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}
