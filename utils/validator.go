package utils

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func ValidateStruct(s interface{}) []*string {
    var errors []*string
    err := validate.Struct(s)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            msg := err.Field() + " is " + err.Tag()
            errors = append(errors, &msg)
        }
    }
    return errors
}