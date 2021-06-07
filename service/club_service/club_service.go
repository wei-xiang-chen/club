package club_service

import (
	"club/model"
	"club/pojo"
	appError "club/pojo/error"
)

var (
	clubModel model.Club
	userModel model.User
)

func GetList(clubId *int, topic *string, clubName *string, offset *int, limit *int) ([]*model.Club, error) {

	clubs, err := clubModel.GetList(clubId, topic, clubName, offset, limit)
	if err != nil {
		return nil, err
	}

	return clubs, nil
}

func Insert(club *pojo.Club) error {

	err := clubModel.Insert(club)
	if err != nil {
		return err
	}

	err = userModel.SetClubId(club.Owner, club.Id)
	if err != nil {
		return err
	}

	return nil
}

func Join(userId *int, clubId *int) error {

	clubExist, err := clubModel.CheckClubExist(clubId)
	if err != nil {
		return err
	}
	if !clubExist {
		return appError.AppError{Message: "The club does not exist."}
	}

	originalClubId, err := userModel.GetUserClubById(userId)
	if err != nil {
		return err
	}
	if originalClubId != nil && *originalClubId != *clubId {
		return appError.AppError{Message: "The user already in the room. Please leave the room first."}
	}

	err = userModel.SetClubId(userId, clubId)
	if err != nil {
		return err
	}

	return nil
}

func Leave(userId *int) error {

	club, err := clubModel.FindByOwner(userId)
	if err != nil {
		return err
	}

	if club != nil {
		err = club.DeleteClubById(&club.Id)
		if err != nil {
			return err
		}

		err = userModel.ClearClub(&club.Id)
		if err != nil {
			return err
		}
	} else {
		err = userModel.SetClubId(userId, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
