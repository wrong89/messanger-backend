package http

type RegisterRes struct {
	ID uint64 `json:"id"`
}

type LoginRes struct {
	Token string `json:"token"`
}
