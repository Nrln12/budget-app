package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Error     string `json:"error"`
	Key       string `json:"key"`
	Condition string `json:"condition"`
}

func (h *Handler) ValidateRequestBody(c echo.Context, payload interface{}) []*ValidationError {
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())
	var errors []*ValidationError
	err := validate.Struct(payload)
	validationErrors, ok := err.(validator.ValidationErrors)
	if ok {
		reflected := reflect.ValueOf(payload)
		for _, validationError := range validationErrors {
			field, _ := reflected.Type().FieldByName(validationError.StructField())
			key := field.Tag.Get("json")
			if key == "" {
				key = strings.ToLower(validationError.StructField())
			}
			condition := validationError.Tag()
			keyToTitleCase := strings.Replace(key, "_", " ", -1)
			param := validationError.Param()
			errMessage := keyToTitleCase + " field is " + condition
			switch condition {
			case "required":
				errMessage = keyToTitleCase + " is required"
			case "email":
				errMessage = keyToTitleCase + " is invalid email"
			case "min":
				if _, err := strconv.Atoi(param); err == nil {
					errMessage = fmt.Sprintf("%s must be at least %s", keyToTitleCase, param)
				}
			}

			fmt.Println("key: ", key)
			fmt.Println(validationError.ActualTag())
			currentValidationError := &ValidationError{
				Error:     errMessage,
				Key:       key,
				Condition: condition,
			}
			errors = append(errors, currentValidationError)
		}
	}
	return errors
}
