// package validator

// import (
// 	"fmt"
// 	errutils "go_app/pkg/errors"
// 	"log"
// 	"reflect"
// 	"strings"

// 	"github.com/go-playground/validator/v10"
// )

// type CustomValidator struct {
// 	v         *validator.Validate
// 	passwdErr error
// }

// func NewCustomValidator() *CustomValidator {
// 	v := validator.New()
// 	cv := &CustomValidator{v: v}

// 	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
// 		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
// 		if name == "-" {
// 			return ""
// 		}
// 		return name
// 	})

// 	if err := v.RegisterValidation("password", cv.passwordValidate); err != nil {
// 		log.Fatal(errutils.WrapPathErr(err))
// 	}

// 	return cv
// }

// func (cv *CustomValidator) Validate(i any) error {
// 	err := cv.v.Struct(i)
// 	if err != nil {
// 		fieldErr := err.(validator.ValidationErrors)[0]
// 		return cv.newValidationError(fieldErr.Field(), fieldErr.Value(), fieldErr.Tag(), fieldErr.Param())
// 	}
// 	return nil
// }

// func (cv *CustomValidator) newValidationError(field string, value any, tag string, param string) error {
// 	switch tag {
// 	case "required":
// 		return fmt.Errorf("field %s is required", field)
// 	case "email":
// 		return fmt.Errorf("field %s must be a valid email address", field)
// 	case "password":
// 		return cv.passwdErr
// 	case "min":
// 		return fmt.Errorf("field %s must be at least %s characters", field, param)
// 	case "max":
// 		return fmt.Errorf("field %s must be at most %s characters", field, param)
// 	default:
// 		return fmt.Errorf("field %s is invalid", field)
// 	}
// }
