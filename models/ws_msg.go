package models

type WsMsg struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}
