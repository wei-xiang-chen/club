package login_service

import (
	"club/model"
	"club/pojo"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

var (
	userModel model.User
)

func Insert(login *pojo.Login) *pojo.Login {

	login.Uid = uuid.NewV4().String()

	err := userModel.Insert(login)

	if nil == err {
		fmt.Println("fail")
	}

	return login
}
