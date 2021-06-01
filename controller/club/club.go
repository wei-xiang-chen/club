package club

import (
	"club/pojo"
	appError "club/pojo/error"
	"club/pojo/rest"
	"club/service/club_service"
	"club/service/code_service"
	"club/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetList(c *gin.Context) error {
	var topic, clubName *string
	var clubId, offset, limit *int
	var restResult rest.RestResult

	if value, ok := c.GetQuery("clubId"); ok {
		c, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		clubId = &c
	}
	if value, ok := c.GetQuery("topic"); ok {
		topic = &value
	}
	if value, ok := c.GetQuery("clubName"); ok {
		clubName = &value
	}
	offset, limit, err := util.GetPagination(c)
	if err != nil {
		return err
	}

	clubs, err := club_service.GetList(clubId, topic, clubName, offset, limit)
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

	if club.ClubName == nil || club.Topic == nil {
		return appError.AppError{Message: "Check request body. Required fields are not filled."}
	}

	err := code_service.CheckCode("clubs_topic", club.Topic)
	if err != nil {
		return err
	}

	user, err := util.GetUser(c)
	if err != nil {
		return err
	}

	club.Owner = &user.Id

	err = club_service.Insert(&club)
	if err != nil {
		return err
	}

	restResult.Data = club
	c.JSON(http.StatusOK, restResult)
	return nil
}

func Join(c *gin.Context) error {

	clubId, err := strconv.Atoi(c.Param("clubId"))

	if err != nil {
		return err
	}

	user, err := util.GetUser(c)
	if err != nil {
		return err
	}

	err = club_service.Join(&user.Id, &clubId)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func Leave(c *gin.Context) error {

	user, err := util.GetUser(c)
	if err != nil {
		return err
	}

	err = club_service.Leave(&user.Id)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}
