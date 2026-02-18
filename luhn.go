// Package luhn implements the Luhn algorithm (ISO/IEC 7812-1) for generating
// and validating checksums used in credit card numbers, IMEI numbers, and more.
package luhn

import (
	"errors"
	"strconv"
	"strings"
)

var (
	errEmpty      = errors.New("string cannot be empty")
	errSpaces     = errors.New("string cannot contain spaces")
	errNegative   = errors.New("negative numbers are not allowed")
	errFloat      = errors.New("floating point numbers are not allowed")
	errNotNumeric = errors.New("string must be convertible to a number")
	errMinLength  = errors.New("string must be longer than 1 character")
	errRandomMax  = errors.New("string must be less than 100 characters")
	errRandomMin  = errors.New("string must be greater than 1")
)

// validateInput applies shared input validation in spec order.
func validateInput(value string) error {
	if value == "" {
		return errEmpty
	}
	if strings.Contains(value, " ") {
		return errSpaces
	}
	if strings.Contains(value, "-") {
		return errNegative
	}
	if strings.Contains(value, ".") {
		return errFloat
	}
	for _, c := range value {
		if c < '0' || c > '9' {
			return errNotNumeric
		}
	}
	return nil
}

// generateChecksum computes the Luhn check digit for a numeric string.
func generateChecksum(value string) byte {
	sum := 0
	shouldDouble := true

	for i := len(value) - 1; i >= 0; i-- {
		digit := int(value[i] - '0')

		if shouldDouble {
			doubled := digit * 2
			if doubled >= 10 {
				sum += doubled/10 + doubled%10
			} else {
				sum += doubled
			}
		} else {
			sum += digit
		}
		shouldDouble = !shouldDouble
	}

	checkDigit := (10 - (sum % 10)) % 10
	return byte('0' + checkDigit)
}

// Generate calculates and appends a Luhn check digit to value.
// If checksumOnly is true, only the check digit is returned.
// Returns an error if value fails input validation.
func Generate(value string, checksumOnly bool) (string, error) {
	if err := validateInput(value); err != nil {
		return "", err
	}

	check := generateChecksum(value)
	if checksumOnly {
		return string(check), nil
	}
	return value + string(check), nil
}

// Validate determines whether value has a valid Luhn check digit as its last character.
// Returns an error if value fails input validation or has length 1.
func Validate(value string) (bool, error) {
	if err := validateInput(value); err != nil {
		return false, err
	}
	if len(value) == 1 {
		return false, errMinLength
	}

	payload := value[:len(value)-1]
	generated, err := Generate(payload, false)
	if err != nil {
		return false, err
	}
	return generated == value, nil
}

// Random generates a random numeric string of the given length (as a numeric string)
// with a valid Luhn check digit. The first digit is never zero.
// Returns an error if length fails input validation or is out of range [2, 100].
func Random(length string) (string, error) {
	if err := validateInput(length); err != nil {
		return "", err
	}
	n, err := strconv.Atoi(length)
	if err != nil {
		// Overflow means the number is far greater than 100.
		return "", errRandomMax
	}
	if n > 100 {
		return "", errRandomMax
	}
	if n < 2 {
		return "", errRandomMin
	}
	return "", nil
}

// GenerateModN computes a Luhn mod-N check character for the given alphanumeric value.
// n must be between 1 and 36. If checksumOnly is true, only the check character is returned.
func GenerateModN(value string, n int, checksumOnly bool) (string, error) {
	return "", nil
}

// ValidateModN determines whether value has a valid Luhn mod-N check character.
// n must be between 1 and 36.
func ValidateModN(value string, n int) (bool, error) {
	return false, nil
}

// ChecksumModN returns the integer index of the Luhn mod-N check character for value.
// n must be between 1 and 36.
func ChecksumModN(value string, n int) (int, error) {
	return 0, nil
}
