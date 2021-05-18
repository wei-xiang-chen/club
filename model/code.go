package model

import (
	"club/client"

	"github.com/jinzhu/gorm"
)

type Code struct {
	Type    string
	Option  string
	Comment *string
	Sort    *int
	Active  *string
}

func (c *Code) TableName() string {
	return "club_code"
}

func (c *Code) GetByTypes(types []string) (*[]Code, error) {
	var codes []Code

	tx := client.DBEngine.Table(c.TableName()).Select("club_code.type, club_code.option")
	if len(types) != 0 {
		tx = tx.Where("type IN (?)", types)
	}

	if err := tx.Order("type, sort").Find(&codes).Error; err != nil {
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
	}
	return &codes, nil
}

func (c *Code) CheckCode(_type string, option string) (*int, error) {
	var count int

	err := client.DBEngine.Table(c.TableName()).Where("type = ? AND option = ?", _type, option).Count(&count).Error
	if err != nil {
		return nil, err
	}

	return &count, nil
}
