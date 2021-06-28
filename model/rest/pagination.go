package rest

type Pagination struct {
	Offset *int `json:"offset"`
	Limit  *int `json:"limit"`
	Total  *int `json:"total"`
	Page   *int `json:"page"`
}
