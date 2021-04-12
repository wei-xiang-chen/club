package user_service

import (
	"club/model"
	"club/pojo"

	uuid "github.com/satori/go.uuid"
)

var (
	userModel model.User
)

func Insert(user *pojo.User) (*pojo.User, error) {

	user.Uid = uuid.NewV4().String()

	err := userModel.Insert(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
