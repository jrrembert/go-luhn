# go-luhn

A minimal, dependency-free Go library implementing the [Luhn algorithm](https://en.wikipedia.org/wiki/Luhn_algorithm) (ISO/IEC 7812-1) for generating and validating checksums. Commonly used for credit card validation, IMEI numbers, and other numeric identifiers.

## Installation

```bash
go get jrrembert/go-luhn
```

## Usage

```go
import "jrrembert/go-luhn"

// Generate a check digit and append it
result, _ := luhn.Generate("7992739871", false)
// result => "79927398713"

// Return only the check digit
digit, _ := luhn.Generate("7992739871", true)
// digit => "3"

// Validate a number with its check digit
valid, _ := luhn.Validate("79927398713")
// valid => true

invalid, _ := luhn.Validate("79927398710")
// invalid => false

// Generate a random 16-digit number with a valid Luhn checksum
card, _ := luhn.Random("16")
// card => e.g. "4532015112830366"

// Luhn mod-N: generate a check character from an expanded alphabet (0-9A-Z)
result, _ := luhn.GenerateModN("ABC123", 36, false)
// result => "ABC1230"

// Luhn mod-N: validate
valid, _ := luhn.ValidateModN("ABC1230", 36)
// valid => true

// Lower-level: get check digit as an integer
n, _ := luhn.ChecksumModN("ABC123", 36)
// n => 0
```

## API

### `Generate(value string, checksumOnly bool) (string, error)`

Computes a Luhn check digit for `value`.

- If `checksumOnly` is `false` (default): returns `value` with the check digit appended.
- If `checksumOnly` is `true`: returns only the check digit as a single-character string.

### `Validate(value string) (bool, error)`

Returns `true` if the last character of `value` is a valid Luhn check digit, `false` otherwise.

### `Random(length string) (string, error)`

Generates a random numeric string of exactly `length` characters with a valid Luhn checksum. The first digit is never zero.

Valid range: `"2"` to `"100"`.

### `GenerateModN(value string, n int, checksumOnly bool) (string, error)`

Luhn mod-N variant. Computes a check character from an alphabet of `n` code points (`0-9A-Z`, where `n` is 1–36).

### `ValidateModN(value string, n int) (bool, error)`

Returns `true` if the last character of `value` is a valid Luhn mod-N check character.

### `ChecksumModN(value string, n int) (int, error)`

Lower-level mod-N checksum calculation. Accepts alphanumeric input (`0-9`, `A-Z`, `a-z`) and returns the check digit as an integer.

## License

[MIT](LICENSE) © J. Ryan Rembert
