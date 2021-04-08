package model

import (
	"club/client"
	"club/pojo"
	"fmt"

	"github.com/jinzhu/gorm"
)

type User struct {
	Id       int    `gorm:"primaryKey"`
	Uid      string ``
	Nickname string ``
	RoomId   *int   `gorm:"column:room_id"`
}

var (
	db *gorm.DB
)

const (
	TABLE_NAME string = "club_user"
)

func (u *User) Insert(login *pojo.Login) error {

	fmt.Println(login.Nickname)
	res := client.DBEngine.Table(TABLE_NAME).Create(login)
	if res.Error == nil {
		fmt.Println(res.Error)
		return res.Error
	}
	return nil
}
