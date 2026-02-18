// Package luhn implements the Luhn algorithm (ISO/IEC 7812-1) for generating
// and validating checksums used in credit card numbers, IMEI numbers, and more.
package luhn

// Generate calculates and appends a Luhn check digit to value.
// If checksumOnly is true, only the check digit is returned.
// Returns an error if value fails input validation.
func Generate(value string, checksumOnly bool) (string, error) {
	return "", nil
}

// Validate determines whether value has a valid Luhn check digit as its last character.
// Returns an error if value fails input validation or has length 1.
func Validate(value string) (bool, error) {
	return false, nil
}

// Random generates a random numeric string of the given length (as a numeric string)
// with a valid Luhn check digit. The first digit is never zero.
// Returns an error if length fails input validation or is out of range [2, 100].
func Random(length string) (string, error) {
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
