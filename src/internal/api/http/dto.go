package http

type Address struct {
	ID      string `json:"id"`
	Country string `json:"country"`
	City    string `json:"city"`
	Street  string `json:"street"`
}
