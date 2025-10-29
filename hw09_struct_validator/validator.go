package hw09structvalidator

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const ValidateTag = "validate"

const (
	And = "|"
	Or  = ","
)

const (
	Len   = "len:"
	In    = "in:"
	Regex = "regexp:"
	Min   = "min:"
	Max   = "max:"
)

var (
	ErrTypeIsNotStruct    = errors.New("type is not a struct")
	ErrUnknownType        = errors.New("unknown type")
	ErrInvalidTag         = errors.New("invalid check in validate tag")
	ErrConvertStringToInt = errors.New("failed to convert string %v to int: %w")

	ErrIn         = errors.New("'%v' not in array of values: '%v'")
	ErrLen        = errors.New("'%v' does not have length = '%v'")
	ErrRegex      = errors.New("'%v' doesn't match regex: '%v'")
	ErrEmptySlice = errors.New("slice is empty")
	ErrMin        = errors.New("'%v' is less than: '%v'")
	ErrMax        = errors.New("'%v' is greater than: '%v'")
)

type ValidationError struct {
	Field string
	Value interface{}
	Err   error
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("%s: %v occurs err: %s", v.Field, v.Value, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	buffer := bytes.Buffer{}
	for _, val := range v {
		buffer.WriteString(val.Error())
		buffer.WriteString("\n")
	}

	return fmt.Sprint(buffer.String())
}

func Validate(v interface{}) error {
	refValue := reflect.ValueOf(v)
	if refValue.Kind() != reflect.Struct {
		return ErrTypeIsNotStruct
	}

	refType := reflect.TypeOf(v)
	var validationErrs ValidationErrors

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		tag, ok := field.Tag.Lookup(ValidateTag)
		if !ok {
			continue
		}

		if tag != "" {
			var validationErrors ValidationErrors
			var err error

			//nolint:exhaustive
			switch field.Type.Kind() {
			case reflect.Slice:
				validationErrors, err = validateSlice(field.Name, refValue.Field(i), tag)
			case reflect.String, reflect.Int:
				validationErrors, err = validateField(field.Name, refValue.Field(i), tag)
			default:
				return ErrUnknownType
			}

			if err != nil {
				return fmt.Errorf("failed to validate field '%v': %w", field.Name, err)
			}

			validationErrs = append(validationErrs, validationErrors...)
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}

	return nil
}

func validateSlice(fieldName string, value reflect.Value, tag string) (ValidationErrors, error) {
	var validationErrors ValidationErrors
	var err error

	if value.Len() == 0 {
		validationErrors = append(validationErrors, ValidationError{
			Field: fieldName,
			Value: value.String(),
			Err:   ErrEmptySlice,
		})
	}

	for i := 0; i < value.Len(); i++ {
		var validationsErrs ValidationErrors

		//nolint:exhaustive
		switch value.Index(i).Kind() {
		case reflect.String, reflect.Int:
			validationsErrs, err = validateField(fieldName, value.Index(i), tag)
		default:
			return nil, ErrUnknownType
		}
		if err != nil {
			return nil, err
		}

		validationErrors = append(validationErrors, validationsErrs...)
	}

	return validationErrors, nil
}

func validateField(fieldName string, value reflect.Value, tag string) (ValidationErrors, error) {
	var validationErrors ValidationErrors
	var err error

	for _, validationRule := range strings.Split(tag, And) {
		var validErrors ValidationErrors
		switch {
		case value.Kind() == reflect.String && strings.HasPrefix(validationRule, Len):
			rule := strings.TrimPrefix(validationRule, Len)
			validErrors, err = validateLen(fieldName, value, rule)
		case value.Kind() == reflect.String && strings.HasPrefix(validationRule, In):
			rule := strings.TrimPrefix(validationRule, In)
			validErrors = validateIn(fieldName, value, rule)
		case value.Kind() == reflect.Int && strings.HasPrefix(validationRule, In):
			rule := strings.TrimPrefix(validationRule, In)
			validErrors = validateIn(fieldName, value, rule)
		case value.Kind() == reflect.String && strings.HasPrefix(validationRule, Regex):
			rule := strings.TrimPrefix(validationRule, Regex)
			validErrors = validateRegex(fieldName, value, rule)
		case value.Kind() == reflect.Int && strings.HasPrefix(validationRule, Min):
			validErrors, err = validateMin(fieldName, value, validationRule)
		case value.Kind() == reflect.Int && strings.HasPrefix(validationRule, Max):
			validErrors, err = validateMax(fieldName, value, validationRule)
		default:
			return nil, ErrInvalidTag
		}

		validationErrors = append(validationErrors, validErrors...)
	}

	if err != nil {
		return nil, err
	}

	return validationErrors, nil
}

func validateLen(fieldName string, value reflect.Value, rules string) (ValidationErrors, error) {
	var validationErrors ValidationErrors

	for _, rule := range strings.Split(rules, Or) {
		expectedLength, err := strconv.Atoi(rule)
		if err != nil {
			return nil, fmt.Errorf(ErrConvertStringToInt.Error(), rule, err)
		}
		valueString := value.String()
		if expectedLength == len(valueString) {
			return nil, nil
		}
	}

	validationErrors = append(validationErrors, ValidationError{
		Field: fieldName,
		Value: value.String(),
		Err:   fmt.Errorf(ErrLen.Error(), fieldName, rules),
	})

	return validationErrors, nil
}

func validateIn(fieldName string, value reflect.Value, rules string) ValidationErrors {
	var validationErrors ValidationErrors

	for _, rule := range strings.Split(rules, Or) {
		if rule == value.String() {
			return nil
		}
	}

	validationErrors = append(validationErrors, ValidationError{
		Field: fieldName,
		Value: value.String(),
		Err:   fmt.Errorf(ErrIn.Error(), fieldName, rules),
	})

	return validationErrors
}

func validateRegex(fieldName string, value reflect.Value, rule string) ValidationErrors {
	var validationErrors ValidationErrors

	regex := regexp.MustCompile(rule)
	if regex.MatchString(value.String()) {
		return nil
	}

	validationErrors = append(validationErrors, ValidationError{
		Field: fieldName,
		Value: value.String(),
		Err:   fmt.Errorf(ErrRegex.Error(), fieldName, rule),
	})

	return validationErrors
}

func validateMin(fieldName string, value reflect.Value, rule string) (ValidationErrors, error) {
	var validationErrors ValidationErrors

	limit, err := strconv.ParseInt(strings.TrimPrefix(rule, Min), 10, 64)
	if err != nil {
		return nil, fmt.Errorf(ErrConvertStringToInt.Error(), rule, err)
	}

	if value.Int() < limit {
		validationErrors = append(validationErrors, ValidationError{
			Field: fieldName,
			Value: fmt.Sprint(value.Interface()),
			Err:   fmt.Errorf(ErrMin.Error(), fieldName, limit),
		})
	}

	return validationErrors, nil
}

func validateMax(fieldName string, value reflect.Value, rule string) (ValidationErrors, error) {
	var validationErrors ValidationErrors

	limit, err := strconv.ParseInt(strings.TrimPrefix(rule, Max), 10, 64)
	if err != nil {
		return nil, fmt.Errorf(ErrConvertStringToInt.Error(), rule, err)
	}

	if value.Int() > limit {
		validationErrors = append(validationErrors, ValidationError{
			Field: fieldName,
			Value: fmt.Sprint(value.Interface()),
			Err:   fmt.Errorf(ErrMax.Error(), fieldName, limit),
		})
	}

	return validationErrors, nil
}
