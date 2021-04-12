package model

import (
	"club/client"
	"club/pojo"
)

type User struct {
	Id       int    `gorm:"primaryKey"`
	Uid      string ``
	Nickname string ``
	ClubId   *int   `gorm:"column:club_id"`
}

func (u *User) TableName() string {
	return "club_user"
}

func (u *User) Insert(user *pojo.User) error {

	err := client.DBEngine.Table(u.TableName()).Create(user).Error
	if err == nil {
		return err
	}

	return nil
}

func (u *User) SetClubId(userId int, clubId *int) error {
	err := client.DBEngine.Table(u.TableName()).Where("id = ?", userId).Update("club_id", clubId).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *User) ClearClub(clubId int) error {
	err := client.DBEngine.Table(u.TableName()).Where("club_id = ?", clubId).Update("club_id", nil).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *User) FindUserByUid(uid string) (*User, error) {
	var user User

	if err := client.DBEngine.Table(u.TableName()).Where("uid = ?", uid).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
