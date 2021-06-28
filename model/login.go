package model

type User struct {
	Id       int    `json:"id"`
	Uid      string `json:"uid"`
	Nickname string `json:"nickname"`
}
