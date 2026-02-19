// Package luhn implements the Luhn algorithm (ISO/IEC 7812-1) for generating
// and validating checksums used in credit card numbers, IMEI numbers, and more.
package luhn

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const codePoints = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	errEmpty         = errors.New("string cannot be empty")
	errSpaces        = errors.New("string cannot contain spaces")
	errNegative      = errors.New("negative numbers are not allowed")
	errFloat         = errors.New("floating point numbers are not allowed")
	errNotNumeric    = errors.New("string must be convertible to a number")
	errMinLength     = errors.New("string must be longer than 1 character")
	errRandomMax     = errors.New("string must be less than 100 characters")
	errRandomMin     = errors.New("string must be greater than 1")
	errInvalidN      = errors.New("n must be between 1 and 36")
	errModNMaxLength = errors.New("string must be less than 10000 characters")
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
	// Use constant-time comparison to prevent timing side-channel attacks
	// that could reveal information about valid check digits.
	return subtle.ConstantTimeCompare([]byte(generated), []byte(value)) == 1, nil
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

	// Generate n-1 random digits (first digit 1-9, rest 0-9)
	buf := make([]byte, n-1)
	first, err := rand.Int(rand.Reader, big.NewInt(9))
	if err != nil {
		return "", err
	}
	buf[0] = byte('1' + first.Int64())
	for i := 1; i < n-1; i++ {
		d, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		buf[i] = byte('0' + d.Int64())
	}

	// Append check digit via Generate
	result, err := Generate(string(buf), false)
	if err != nil {
		return "", err
	}
	return result, nil
}

// validateModNInput validates input for mod-N functions.
// Checks empty and spaces (shared with validateInput), then validates each
// character against the CODE_POINTS alphabet for the given n.
func validateModNInput(value string, n int) error {
	if value == "" {
		return errEmpty
	}
	if strings.Contains(value, " ") {
		return errSpaces
	}
	for i := 0; i < len(value); i++ {
		if charIndex(value[i], n) < 0 {
			return fmt.Errorf("invalid character: %q", value[i])
		}
	}
	return nil
}

// charIndex returns the index of c in the CODE_POINTS alphabet, or -1 if not found.
func charIndex(c byte, n int) int {
	var idx int
	if c >= '0' && c <= '9' {
		idx = int(c - '0')
	} else if c >= 'A' && c <= 'Z' {
		idx = int(c-'A') + 10
	} else if c >= 'a' && c <= 'z' {
		idx = int(c-'a') + 10
	} else {
		return -1
	}
	if idx >= n {
		return -1
	}
	return idx
}

// generateChecksumModN computes the Luhn mod-N check character index.
func generateChecksumModN(value string, n int) (int, error) {
	sum := 0
	shouldDouble := true

	for i := len(value) - 1; i >= 0; i-- {
		idx := charIndex(value[i], n)
		if idx < 0 {
			return 0, fmt.Errorf("invalid character: %q", value[i])
		}

		if shouldDouble {
			doubled := idx * 2
			if doubled >= n {
				sum += doubled/n + doubled%n
			} else {
				sum += doubled
			}
		} else {
			sum += idx
		}
		shouldDouble = !shouldDouble
	}

	checkIdx := (n - (sum % n)) % n
	return checkIdx, nil
}

// GenerateModN computes a Luhn mod-N check character for the given alphanumeric value.
// n must be between 1 and 36. If checksumOnly is true, only the check character is returned.
func GenerateModN(value string, n int, checksumOnly bool) (string, error) {
	if n < 1 || n > 36 {
		return "", errInvalidN
	}
	if err := validateModNInput(value, n); err != nil {
		return "", err
	}
	if len(value) >= 10000 {
		return "", errModNMaxLength
	}

	checkIdx, err := generateChecksumModN(value, n)
	if err != nil {
		return "", err
	}

	checkChar := codePoints[checkIdx]
	if checksumOnly {
		return string(checkChar), nil
	}
	return value + string(checkChar), nil
}

// ValidateModN determines whether value has a valid Luhn mod-N check character.
// n must be between 1 and 36.
func ValidateModN(value string, n int) (bool, error) {
	if n < 1 || n > 36 {
		return false, errInvalidN
	}
	if err := validateModNInput(value, n); err != nil {
		return false, err
	}
	if len(value) == 1 {
		return false, errMinLength
	}
	if len(value) >= 10000 {
		return false, errModNMaxLength
	}

	// Normalize to uppercase so that lowercase input matches the uppercase
	// CODE_POINTS alphabet used by GenerateModN.
	upper := strings.ToUpper(value)

	payload := upper[:len(upper)-1]
	generated, err := GenerateModN(payload, n, false)
	if err != nil {
		return false, err
	}
	// Use constant-time comparison to prevent timing side-channel attacks
	// that could reveal information about valid check digits.
	return subtle.ConstantTimeCompare([]byte(generated), []byte(upper)) == 1, nil
}

// ChecksumModN returns the integer index of the Luhn mod-N check character for value.
// n must be between 1 and 36.
func ChecksumModN(value string, n int) (int, error) {
	if n < 1 || n > 36 {
		return 0, errInvalidN
	}
	if err := validateModNInput(value, n); err != nil {
		return 0, err
	}
	if len(value) >= 10000 {
		return 0, errModNMaxLength
	}
	return generateChecksumModN(value, n)
}
