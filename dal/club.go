package dal

import (
	"club/clients"
	"club/models"
	appError "club/models/error"
	"time"

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

func (c *Club) Insert(club *models.Club) error {

	err := clients.DBEngine.Table(c.TableName()).Create(&club).Error
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

func (c *Club) GetList(clubId *int, topic *string, clubName *string, offset *int, limit *int) ([]*Club, error) {
	var clubs []*Club

	tx := clients.DBEngine.Table(c.TableName()).Select("club_clubs.id, club_clubs.club_name, club_clubs.topic, club_clubs.population, club_user.id, club_user.nickname").Joins("LEFT JOIN club_user ON club_clubs.owner = club_user.id")

	if clubId != nil {
		tx = tx.Where("club_clubs.id = ?", *clubId)
	}
	if topic != nil {
		tx = tx.Where("topic = ?", *topic)
	}
	if clubName != nil {
		tx = tx.Where("club_name LIKE ?", "%"+*clubName+"%")
	}

	err := tx.Offset(*offset).Limit(*limit).Find(&clubs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return clubs, nil
}

func (c *Club) FIndExpired(time *time.Time) (*[]Club, error) {
	var clubs []Club

	err := clients.DBEngine.Table(c.TableName()).Select("club_clubs.id").Joins("LEFT JOIN club_user ON club_clubs.owner = club_user.id").Where("club_user.disconnection_time <= ?", *time).Find(&clubs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &clubs, nil
}

func (c *Club) FindByOwner(owner *int) (*Club, error) {
	var club Club

	err := clients.DBEngine.Table(c.TableName()).Where("owner = ?", *owner).Find(&club).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &club, nil
}

func (c *Club) CheckClubExist(id *int) (bool, error) {
	var count int

	if err := clients.DBEngine.Table(c.TableName()).Where("id = ?", *id).Count(&count).Error; err != nil {
		if err != nil {
			return false, err
		}
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (c *Club) DeleteClubById(id *int) error {

	if err := clients.DBEngine.Table(c.TableName()).Delete(&Club{}, *id).Error; err != nil {
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Club) UpdatePopulation(id *int, population *int) error {

	if err := clients.DBEngine.Table(c.TableName()).Where("id = ?", *id).Update("population", population).Error; err != nil {
		return err
	}

	return nil
}
