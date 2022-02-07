package user_service

import (
	"club/dal"
	"club/models"

	uuid "github.com/satori/go.uuid"
)

var (
	userModel dal.User
)

func Insert(user *models.User) (*models.User, error) {

	user.Uid = uuid.NewV4().String()

	err := userModel.Insert(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
