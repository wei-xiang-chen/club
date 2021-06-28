package rest

type RestResult struct {
	Data  interface{} `json:"data"`
	Page  *Pagination `json:"page"`
	Error *RestError  `json:"error"`
}
