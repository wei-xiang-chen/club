package user_service

import (
	"club/dao"
	"club/model"

	uuid "github.com/satori/go.uuid"
)

var (
	userModel dao.User
)

func Insert(user *model.User) (*model.User, error) {

	user.Uid = uuid.NewV4().String()

	err := userModel.Insert(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
