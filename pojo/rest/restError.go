package rest

type RestError struct {
	Message     *string `json:"message"`
	Description *string `json:"description"`
}
