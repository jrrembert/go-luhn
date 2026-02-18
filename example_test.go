package luhn_test

import (
	"fmt"

	luhn "github.com/jrrembert/go-luhn"
)

func ExampleGenerate() {
	// Compute the check digit for an IMEI body (14 digits) and append it.
	imei, _ := luhn.Generate("35209900176148", false)
	fmt.Println(imei)

	// Get only the check digit.
	checkDigit, _ := luhn.Generate("35209900176148", true)
	fmt.Println(checkDigit)

	// Output:
	// 352099001761481
	// 1
}

func ExampleValidate() {
	// Validate a complete IMEI number (15 digits).
	valid, _ := luhn.Validate("352099001761481")
	fmt.Println(valid)

	// An IMEI with a wrong check digit is rejected.
	invalid, _ := luhn.Validate("352099001761480")
	fmt.Println(invalid)

	// Output:
	// true
	// false
}

func ExampleRandom() {
	// Generate a random 16-digit number with a valid Luhn check digit.
	result, _ := luhn.Random("16")
	fmt.Println(len(result))

	// The result always passes validation.
	valid, _ := luhn.Validate(result)
	fmt.Println(valid)

	// Output:
	// 16
	// true
}

func ExampleGenerateModN() {
	// Append a check character using a base-16 (hexadecimal) alphabet.
	result, _ := luhn.GenerateModN("FF", 16, false)
	fmt.Println(result)

	// Use base-36 (0-9, A-Z) for alphanumeric payloads.
	result36, _ := luhn.GenerateModN("HELLO", 36, false)
	fmt.Println(result36)

	// Output:
	// FF2
	// HELLOJ
}

func ExampleValidateModN() {
	// Validate an alphanumeric string with a Luhn mod-36 check character.
	valid, _ := luhn.ValidateModN("HELLOJ", 36)
	fmt.Println(valid)

	// An incorrect check character is rejected.
	invalid, _ := luhn.ValidateModN("HELLOA", 36)
	fmt.Println(invalid)

	// Output:
	// true
	// false
}

func ExampleChecksumModN() {
	// Get the check character index for a base-36 payload.
	idx, _ := luhn.ChecksumModN("HELLO", 36)
	fmt.Println(idx)

	// Index 19 corresponds to 'J' in the 0-9A-Z alphabet.
	codePoints := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fmt.Println(string(codePoints[idx]))

	// Output:
	// 19
	// J
}
