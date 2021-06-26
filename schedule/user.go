package schedule

import (
	"club/model"
	"club/ws/club_ws"
	"log"
	"time"
)

func DeleteExpiredUser() {
	var userModel model.User
	var clubModel model.Club

	for {
		targetTime := time.Now().Add(-1 * time.Minute)
		clubs, err := clubModel.FIndExpired(&targetTime)
		if err != nil {
			log.Printf("error: %v", err)
		}

		for _, club := range *clubs {
			s := club_ws.Subscription{Room: club.Id}
			club_ws.H.CloseRoom <- s

			err = clubModel.DeleteClubById(&club.Id)
			if err != nil {
				log.Printf("error: %v", err)
			}
		}

		err = userModel.DeleteExpired(&targetTime)
		if err != nil {
			log.Printf("error: %v", err)
		}

		time.Sleep(15 * time.Minute)
	}
}
