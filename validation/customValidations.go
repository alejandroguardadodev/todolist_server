package validation

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// The function `customValidationRequiredEnum` checks if a custom type enum value is not empty for
// validation.
func customValidationRequiredEnum(fl validator.FieldLevel) bool {
	enum := (fl.Field().Interface()).(CustomTypeBehavior)
	return !enum.Empty()
}

// The function customValidationCheckEnumVal performs custom validation on an enum value.
func customValidationCheckEnumVal(fl validator.FieldLevel) bool {
	enum := (fl.Field().Interface()).(CustomTypeBehavior)
	log.Println(enum)
	return enum.CheckValue()
}

// The customValidationDate function in Go checks if a field value matches yyyy-mm-dd format regex
// pattern.
func customValidationDate(fl validator.FieldLevel) bool {

	val := fl.Field().String()

	if len(strings.Trim(val, " ")) == 0 {
		return true
	}

	strRegexDate := "\\d{4}-\\d{1,2}-\\d{1,2}"

	match, _ := regexp.MatchString(strRegexDate, val)

	return match
}

// The function `customValidationDateLessThanNow` checks if a given date is before the current time.
func customValidationDateAfterOrEqualThanToday(fl validator.FieldLevel) bool {

	val := fl.Field().String()

	if len(strings.Trim(val, " ")) == 0 {
		return true
	}

	date, err := time.Parse(time.DateOnly, val)

	if err != nil {
		log.Println(err)
		return false
	}

	return !date.Before(time.Now())
}
