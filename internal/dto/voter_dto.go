package dto

type BindNIMRequest struct {
	NIM  string `json:"nim"`
	Nama string `json:"nama"`
}

type BindNIMResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
