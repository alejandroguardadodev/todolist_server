package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ErrField struct {
	FieldName  string `json:"field_name"`
	ErrorTitle string `json:"error_title"`
	Value      string `json:"value"`
}

var Validate *validator.Validate

func registerCustomValidations() {

}

func getJsonTag(fieldname string, val reflect.Value) string {
	for i := 0; i < val.Type().NumField(); i++ {
		f := val.Type().Field(i)

		if f.Name == fieldname {
			return f.Tag.Get("json")
		}
	}

	return fieldname
}

func ValidationInit() {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	registerCustomValidations()
}

func GetValidateInformation(err error, element any) *[]ErrField {

	errFields := []ErrField{}
	val := reflect.ValueOf(element)

	for _, err := range err.(validator.ValidationErrors) {

		var value string

		if val, ok := err.Value().(string); ok {
			value = val
		}

		errFields = append(errFields, ErrField{
			FieldName:  getJsonTag(err.Field(), val),
			ErrorTitle: err.ActualTag(),
			Value:      value,
		})
	}

	return &errFields
}
