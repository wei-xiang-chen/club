package pojo

type Club struct {
	Id       *int    `json:"id"`
	ClubName *string `json:"clubName"`
	Topic    *string `json:"topic"`
	Owner    *int    `json:"owner"`
}
