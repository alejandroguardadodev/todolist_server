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

func dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// The function `customValidationDateLessThanNow` checks if a given date is before the current time.
func customValidationDateAfterOrEqualThanToday(fl validator.FieldLevel) bool {

	val := fl.Field().String()

	if len(strings.Trim(val, " ")) == 0 {
		return true
	}

	if strings.Contains(val, "T") {
		val = strings.Split(val, "T")[0]
	}

	date, err := time.Parse(time.DateOnly, val)

	if err != nil {
		log.Println(err)
		return false
	}

	//og.Println(val + " - " + time.Now().String())

	return !date.Before(time.Now()) || dateEqual(date, time.Now())
}
