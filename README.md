# go-luhn

[![CI](https://github.com/jrrembert/go-luhn/actions/workflows/ci.yml/badge.svg)](https://github.com/jrrembert/go-luhn/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/jrrembert/go-luhn/graph/badge.svg)](https://codecov.io/gh/jrrembert/go-luhn)
[![Go Reference](https://pkg.go.dev/badge/github.com/jrrembert/go-luhn.svg)](https://pkg.go.dev/github.com/jrrembert/go-luhn)
[![Go Report Card](https://goreportcard.com/badge/github.com/jrrembert/go-luhn)](https://goreportcard.com/report/github.com/jrrembert/go-luhn)

A Go implementation of the Luhn algorithm for generating and validating checksums.

Published as [`github.com/jrrembert/go-luhn`](https://pkg.go.dev/github.com/jrrembert/go-luhn).

## Getting Started

### Prerequisites

Install [Go](https://go.dev/) (>=1.21).

### Installation

```bash
$ go get github.com/jrrembert/go-luhn
```

### Usage

```go
import luhn "github.com/jrrembert/go-luhn"

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

## Commands

```bash
# Build the library
$ go build ./...

# Run tests
$ go test ./...

# Run Go vet
$ go vet ./...

# Lint
$ golangci-lint run
```

## Documentation

- [Specification](docs/SPEC.md) - Algorithm specification and test vectors
- [Contributing](CONTRIBUTING.md) - How to contribute
- [Code of Conduct](CODE_OF_CONDUCT.md) - Community guidelines
- [Security](SECURITY.md) - Reporting vulnerabilities

## Contact

Email: [J. Ryan Rembert](mailto:j.ryan.rembert@gmail.com)

## License

[MIT](LICENSE)

Copyright © 2022-2026 J. Ryan Rembert
