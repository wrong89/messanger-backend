package http

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrAddressIsEmpty = errors.New("address is empty")
)

type CreateChatReqDTO struct {
	Address string `json:"address"`
}

func (c CreateChatReqDTO) Validate() error {
	if c.Address == "" {
		return ErrAddressIsEmpty
	}

	return nil
}

type ErrorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func NewErrorDTO(err error) ErrorDTO {
	return ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}
}

func (e ErrorDTO) String() string {
	b, _ := json.MarshalIndent(e, "", "    ")

	return string(b)
}
