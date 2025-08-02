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
			if unicode.IsDigit(inputRune) {
				return "", ErrInvalidString
			} else if inputRune == slashRune {
				isPrevSlash = true
			} else {
				prevRune = inputRune
			}
		case len(input) - 1:
			if unicode.IsDigit(inputRune) {
				if unicode.IsDigit(prevRune) && isNumberWithSlash {
					inputRuneInt, err := strconv.Atoi(string(inputRune))
					if err != nil {
						return "", err
					}
					resultBuilder.WriteString(strings.Repeat(string(prevRune), inputRuneInt))
				} else if unicode.IsDigit(prevRune) {
					return "", ErrInvalidString
				} else if isPrevSlash {
					resultBuilder.WriteRune(inputRune)
				} else if isNumberWithSlash {
					inputRuneInt, err := strconv.Atoi(string(inputRune))
					if err != nil {
						return "", err
					}
					resultBuilder.WriteString(strings.Repeat(string(prevRune), inputRuneInt))
				} else if unicode.IsDigit(inputRune) && !unicode.IsDigit(prevRune) {
					inputRuneInt, err := strconv.Atoi(string(inputRune))
					if err != nil {
						return "", err
					}
					resultBuilder.WriteString(strings.Repeat(string(prevRune), inputRuneInt))
				}
			} else {
				if isPrevSlash && inputRune == slashRune {
					resultBuilder.WriteRune(inputRune)
				} else if isPrevSlash {
					return "", ErrInvalidString
				} else if isNumberWithSlash {
					resultBuilder.WriteRune(prevRune)
					resultBuilder.WriteRune(inputRune)
				} else if unicode.IsDigit(prevRune) {
					resultBuilder.WriteRune(inputRune)
				} else {
					resultBuilder.WriteRune(prevRune)
					resultBuilder.WriteRune(inputRune)
				}
			}
		default:
			if unicode.IsDigit(inputRune) {
				if prevRune != 0 && unicode.IsDigit(prevRune) {
					return "", ErrInvalidString
				} else if isPrevSlash {
					isNumberWithSlash = true
					isPrevSlash = false
				} else if isNumberWithSlash ||
					!unicode.IsDigit(prevRune) ||
					(prevRune != 0 && unicode.IsDigit(prevRune) && isNumberWithSlash) {
					inputRuneInt, err := strconv.Atoi(string(inputRune))
					if err != nil {
						return "", err
					}
					resultBuilder.WriteString(strings.Repeat(string(prevRune), inputRuneInt))
					isPrevSlash = false
					isNumberWithSlash = false
				}
			} else {
				if isPrevSlash && inputRune == slashRune {
					isPrevSlash = false
					isNumberWithSlash = false
				} else if isPrevSlash {
					return "", ErrInvalidString
				} else if isNumberWithSlash || (!unicode.IsDigit(prevRune) && inputRune == slashRune) {
					resultBuilder.WriteRune(prevRune)
					if inputRune == slashRune {
						isPrevSlash = true
					}
					isNumberWithSlash = false
				} else if !unicode.IsDigit(prevRune) {
					resultBuilder.WriteRune(prevRune)
					isPrevSlash = false
					isNumberWithSlash = false
				}
			}
			prevRune = inputRune
		}
	}
	return resultBuilder.String(), nil
}
