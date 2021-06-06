package schedule

import (
	"club/model"
	"log"
	"time"
)

func DeleteExpiredUser() {
	var userModel model.User
	for {
		targetTime := time.Now().Add(-5 * time.Minute)

		err := userModel.DeleteExpired(&targetTime)
		if err != nil {
			log.Printf("error: %v", err)
		}

		time.Sleep(30 * time.Minute)
	}
}
