package model

import (
	"club/client"
	"club/pojo"
	appError "club/pojo/error"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Club struct {
	Id         int     `gorm:"primaryKey" json:"id"`
	ClubName   *string `gorm:"column:club_name" json:"clubName"`
	Topic      *string `json:"topic"`
	Owner      *User   `gorm:"embedded" json:"owner"`
	Population *int    `json:"population"`
}

func (c *Club) TableName() string {
	return "club_clubs"
}

func (c *Club) Insert(club *pojo.Club) error {

	err := client.DBEngine.Table(c.TableName()).Create(&club).Error

	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" {
			return appError.AppError{Message: "您已是房主，請離開房間後再次建立。"}
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Club) GetList(topic string, clubName string, offset int, limit int) ([]*Club, error) {
	var clubs []*Club

	tx := client.DBEngine.Table(c.TableName()).Select("club_clubs.id, club_clubs.club_name, club_clubs.topic, club_clubs.population, club_user.id, club_user.nickname").Joins("LEFT JOIN club_user ON club_clubs.owner = club_user.id")

	if topic != "" {
		tx = tx.Where("topic = ?", topic)
	}
	if clubName != "" {
		tx = tx.Where("club_name = ?", clubName)
	}

	err := tx.Offset(offset).Limit(limit).Find(&clubs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return clubs, nil
}

func (c *Club) FindByOwner(owner int) (*Club, error) {
	var club Club

	err := client.DBEngine.Table(c.TableName()).Where("owner = ?", owner).Find(&club).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &club, nil
}

func (c *Club) Delete(id int) error {

	err := client.DBEngine.Table(c.TableName()).Delete(&Club{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
