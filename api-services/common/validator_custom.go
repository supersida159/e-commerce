package common

import (
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

// MM-DD-YYYY
func DateCustom(fl validator.FieldLevel) bool {
	switch v := fl.Field(); v.Kind() {
	case reflect.String:
		const shortForm = "01-02-2006"
		_, err := time.Parse(shortForm, v.String())
		if err == nil {
			return true
		}
	default:
		return false
	}

	return false
}

// HH-MM AM/PM
func TimeCustom12h(fl validator.FieldLevel) bool {
	switch v := fl.Field(); v.Kind() {
	case reflect.String:
		const regex = "^(0?[1-9]|1[012]):([0-5][0-9])[ap]m$"
		match, err := regexp.MatchString(regex, v.String())
		if err == nil && match {
			return true
		}
	default:
		return false
	}

	return false
}

// HH-MM
func TimeCustom24h(fl validator.FieldLevel) bool {
	switch v := fl.Field(); v.Kind() {
	case reflect.String:
		const regex = "^(0?[1-9]|1[0-9]|2[0-4]):([0-5][0-9])$"
		match, err := regexp.MatchString(regex, v.String())
		if err == nil && match {
			return true
		}
	default:
		return false
	}

	return false
}

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) ValidateStruct(data interface{}) *AppError {
	err := v.validate.Struct(data)
	if err == nil {
		return nil
	}

	for _, validationErr := range err.(validator.ValidationErrors) {
		field := validationErr.Field()
		tag := validationErr.Tag()
		param := validationErr.Param()

		switch tag {
		case "required":
			return ErrMissingRequiredField(field)
		case "min":
			minValue, _ := strconv.Atoi(param)
			return ErrFieldBelowMinimum(field, minValue)
		case "max":
			maxValue, _ := strconv.Atoi(param)
			return ErrFieldAboveMaximum(field, maxValue)
		case "email":
			return ErrInvalidEmail()
		case "oneof":
			return ErrInvalidInputData(validationErr) // can add oneof error later
		case "url":
			return ErrInvalidURL(field)
		case "datetime":
			return ErrInvalidDateTime(field, param) // passing the expected format
		default:
			return ErrInvalidInputData(validationErr)
		}
	}

	return ErrInvalidInputData(err)
}
