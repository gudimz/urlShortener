package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Validator struct {
	v *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{v: validator.New()}
}

func (v *Validator) Validate(i interface{}) error {
	err := v.v.Struct(i)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
