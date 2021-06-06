package model

import (
	"club/client"
	"club/pojo"
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	Id                int        `gorm:"primaryKey" json:"id"`
	Uid               *string    `json:"uid"`
	Nickname          *string    `json:"nickname"`
	ClubId            *int       `gorm:"column:club_id" json:"clubId"`
	DisconnectionTIme *time.Time `gorm:"column:disconnection_time" json:"disconnectionTIme"`
}

func (u *User) TableName() string {
	return "club_user"
}

func (u *User) Insert(user *pojo.User) error {

	if err := client.DBEngine.Table(u.TableName()).Create(user).Error; err == nil {
		return err
	}

	return nil
}

func (u *User) SetClubId(userId *int, clubId *int) error {

	if err := client.DBEngine.Table(u.TableName()).Where("id = ?", *userId).Update("club_id", clubId).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) DeleteUserById(id *int) error {

	if err := client.DBEngine.Table(u.TableName()).Delete(&User{}, *id).Error; err != nil {
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *User) ClearClub(clubId *int) error {
	err := client.DBEngine.Table(u.TableName()).Where("club_id = ?", *clubId).Update("club_id", nil).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUserClubById(userId *int) (*int, error) {
	var clubId []*int

	if err := client.DBEngine.Table(u.TableName()).Where("id = ?", *userId).Pluck("club_id", &clubId).Error; err != nil {
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
	}
	return clubId[0], nil
}

func (u *User) FindUserByUid(uid *string) (*User, error) {
	var user User

	if err := client.DBEngine.Table(u.TableName()).Where("uid = ?", *uid).Find(&user).Error; err != nil {
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
	}
	return &user, nil
}

func (u *User) CheckUserExist(id *int) (bool, error) {
	var count int

	if err := client.DBEngine.Table(u.TableName()).Where("id = ?", *id).Count(&count).Error; err != nil {
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

func (u *User) CompareUserAndClub(id *int, clubId *int) (bool, error) {
	var count int

	if err := client.DBEngine.Table(u.TableName()).Where("id = ? AND club_id = ?", *id, *clubId).Count(&count).Error; err != nil {
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

func (u *User) UpdateDisconnectionTime(id *int, time *time.Time) error {

	if err := client.DBEngine.Table(u.TableName()).Where("id = ?", *id).Update("disconnection_time", time).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) DeleteExpired(time *time.Time) error {

	if err := client.DBEngine.Debug().Table(u.TableName()).Where("disconnection_time < ?", time).Delete(&User{}).Error; err != nil {
		if err != nil {
			return err
		}
	}

	return nil
}
