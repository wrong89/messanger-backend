package http

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrNameIsEmpty     = errors.New("name is empty")
	ErrLoginIsEmpty    = errors.New("login is empty")
	ErrPasswordIsEmpty = errors.New("password is empty")
)

type RegisterReqDTO struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r RegisterReqDTO) Validate() error {
	if r.Login == "" {
		return ErrLoginIsEmpty
	}
	if r.Name == "" {
		return ErrNameIsEmpty
	}
	if r.Password == "" {
		return ErrPasswordIsEmpty
	}

	return nil
}

type LoginReqDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r LoginReqDTO) Validate() error {
	if r.Login == "" {
		return ErrLoginIsEmpty
	}
	if r.Password == "" {
		return ErrPasswordIsEmpty
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
