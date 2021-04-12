package model

import (
	"club/client"
	"club/pojo"

	"github.com/jinzhu/gorm"
)

type Club struct {
	Id         int    `gorm:"primaryKey"`
	ClubName   string `gorm:"column:club_name"`
	Topic      string ``
	Owner      int    ``
	population int    ``
}

func (c *Club) TableName() string {
	return "club_clubs"
}

func (c *Club) Insert(club *pojo.Club) error {

	err := client.DBEngine.Table(c.TableName()).Create(&club).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *Club) GetList(topic string, clubName string, offset int, limit int) ([]*Club, error) {
	var clubs []*Club

	tx := client.DBEngine.Table(c.TableName())

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

	err := client.DBEngine.Table(c.TableName()).Delete(id).Error
	if err != nil {
		return err
	}

	return nil
}
