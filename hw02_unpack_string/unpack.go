package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const slashRune rune = 92

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	resultBuilder := strings.Builder{}

	if len([]rune(input)) <= 1 {
		return input, nil
	}

	var prevRune rune
	var isPrevSlash bool
	var isNumberWithSlash bool

	for number, inputRune := range input {
		switch number {
		case 0:
			err := checkFirstRune(inputRune, &isPrevSlash, &prevRune)
			if err != nil {
				return "", err
			}
		case len(input) - 1:
			err := checkLastRunes(inputRune, &prevRune, &isNumberWithSlash, &resultBuilder, &isPrevSlash)
			if err != nil {
				return "", err
			}
		default:
			err := checkMiddleRune(inputRune, &prevRune, &isPrevSlash, &isNumberWithSlash, &resultBuilder)
			if err != nil {
				return "", err
			}
		}
	}
	return resultBuilder.String(), nil
}

func checkFirstRune(inputRune int32, isPrevSlash *bool, prevRune *rune) error {
	if unicode.IsDigit(inputRune) {
		return ErrInvalidString
	} else if inputRune == slashRune {
		*isPrevSlash = true
		rewriteFlags(isPrevSlash, nil, true, false)
	}
	*prevRune = inputRune

	return nil
}

func checkLastRunes(inputRune int32,
	prevRune *rune,
	isNumberWithSlash *bool,
	resultBuilder *strings.Builder,
	isPrevSlash *bool,
) error {
	switch unicode.IsDigit(inputRune) {
	case true:
		err := checkLastRuneInputRuneIsDigit(inputRune, prevRune, isNumberWithSlash, resultBuilder, isPrevSlash)
		if err != nil {
			return err
		}
	case false:
		err := checkLastRuneInputRuneIsNotDigit(inputRune, prevRune, isPrevSlash, resultBuilder)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkLastRuneInputRuneIsNotDigit(inputRune int32,
	prevRune *rune,
	isPrevSlash *bool,
	resultBuilder *strings.Builder,
) error {
	switch {
	case (*isPrevSlash && inputRune == slashRune) || unicode.IsDigit(*prevRune):
		resultBuilder.WriteRune(inputRune)
	case *isPrevSlash:
		return ErrInvalidString
	default:
		resultBuilder.WriteRune(*prevRune)
		resultBuilder.WriteRune(inputRune)
	}
	return nil
}

func checkLastRuneInputRuneIsDigit(inputRune int32,
	prevRune *rune,
	isNumberWithSlash *bool,
	resultBuilder *strings.Builder,
	isPrevSlash *bool,
) error {
	switch {
	case unicode.IsDigit(*prevRune) && *isNumberWithSlash:
		writeRepeatString(inputRune, prevRune, resultBuilder)
	case unicode.IsDigit(*prevRune):
		return ErrInvalidString
	case *isPrevSlash:
		resultBuilder.WriteRune(inputRune)
	default:
		writeRepeatString(inputRune, prevRune, resultBuilder)
	}
	return nil
}

func checkMiddleRune(inputRune int32,
	prevRune *rune,
	isPrevSlash *bool,
	isNumberWithSlash *bool,
	resultBuilder *strings.Builder,
) error {
	switch unicode.IsDigit(inputRune) {
	case true:
		err := checkMiddleRuneInputRuneIsDigit(inputRune, prevRune, isPrevSlash, isNumberWithSlash, resultBuilder)
		if err != nil {
			return err
		}
	case false:
		err := checkMiddleRuneInputRuneIsNotDigit(inputRune, prevRune, isPrevSlash, isNumberWithSlash, resultBuilder)
		if err != nil {
			return err
		}
	}
	*prevRune = inputRune
	return nil
}

func checkMiddleRuneInputRuneIsNotDigit(
	inputRune int32,
	prevRune *rune,
	isPrevSlash *bool,
	isNumberWithSlash *bool,
	resultBuilder *strings.Builder,
) error {
	switch {
	case *isPrevSlash && inputRune == slashRune:
		rewriteFlags(isPrevSlash, isNumberWithSlash, false, false)
	case *isPrevSlash:
		resultBuilder.WriteRune(*prevRune)
		rewriteFlags(isPrevSlash, isNumberWithSlash, inputRune == slashRune, false)
	case *isNumberWithSlash || (!unicode.IsDigit(*prevRune) && inputRune == slashRune):
		resultBuilder.WriteRune(*prevRune)
		rewriteFlags(isPrevSlash, isNumberWithSlash, inputRune == slashRune, false)
	case !unicode.IsDigit(*prevRune):
		resultBuilder.WriteRune(*prevRune)
		rewriteFlags(isPrevSlash, isNumberWithSlash, false, false)
	}
	return nil
}

func checkMiddleRuneInputRuneIsDigit(inputRune int32,
	prevRune *rune,
	isPrevSlash *bool,
	isNumberWithSlash *bool,
	resultBuilder *strings.Builder,
) error {
	switch {
	case unicode.IsDigit(*prevRune):
		return ErrInvalidString
	case *isPrevSlash:
		rewriteFlags(isPrevSlash, isNumberWithSlash, false, true)
	default:
		writeRepeatString(inputRune, prevRune, resultBuilder)
		rewriteFlags(isPrevSlash, isNumberWithSlash, false, false)
	}
	return nil
}

func writeRepeatString(inputRune int32, prevRune *rune, resultBuilder *strings.Builder) {
	inputRuneInt, _ := strconv.Atoi(string(inputRune))
	resultBuilder.WriteString(strings.Repeat(string(*prevRune), inputRuneInt))
}

func rewriteFlags(isPrevSlash, isNumberWithSlash *bool, newIsPrevSlashValue, newIsNumberValue bool) {
	*isPrevSlash = newIsPrevSlashValue
	if isNumberWithSlash != nil {
		*isNumberWithSlash = newIsNumberValue
	}
}
