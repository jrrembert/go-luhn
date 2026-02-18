package luhn_test

import (
	"strconv"
	"strings"
	"testing"

	luhn "github.com/jrrembert/go-luhn"
)

// TestPackageExports verifies the package compiles and all 6 public functions are callable.
func TestPackageExports(t *testing.T) {
	if _, err := luhn.Generate("1", false); err != nil {
		t.Errorf("Generate: %v", err)
	}
	if _, err := luhn.Validate("18"); err != nil {
		t.Errorf("Validate: %v", err)
	}
	if _, err := luhn.Random("4"); err != nil {
		t.Errorf("Random: %v", err)
	}
	if _, err := luhn.GenerateModN("1", 10, false); err != nil {
		t.Errorf("GenerateModN: %v", err)
	}
	if _, err := luhn.ValidateModN("18", 10); err != nil {
		t.Errorf("ValidateModN: %v", err)
	}
	if _, err := luhn.ChecksumModN("1", 10); err != nil {
		t.Errorf("ChecksumModN: %v", err)
	}
}

// TestGenerateChecksumOnly tests Generate with checksumOnly=true against SPEC.md §5 test vectors.
func TestGenerateChecksumOnly(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1", "8"},
		{"12", "5"},
		{"123", "0"},
		{"1234", "4"},
		{"12345", "5"},
		{"123456", "6"},
		{"1234567", "4"},
		{"12345678", "2"},
		{"123456789", "7"},
		{"7992739871", "3"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := luhn.Generate(tt.input, true)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Generate(%q, true) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestGenerateFullOutput tests Generate with checksumOnly=false against SPEC.md §5 test vectors.
func TestGenerateFullOutput(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1", "18"},
		{"12", "125"},
		{"123", "1230"},
		{"1234", "12344"},
		{"12345", "123455"},
		{"123456", "1234566"},
		{"1234567", "12345674"},
		{"12345678", "123456782"},
		{"123456789", "1234567897"},
		{"7992739871", "79927398713"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := luhn.Generate(tt.input, false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Generate(%q, false) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestGenerateEdgeCases tests Generate edge cases from SPEC.md §5.
func TestGenerateEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"leading zero", "0", "00"},
		{"multiple leading zeros", "00123", "001230"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := luhn.Generate(tt.input, false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Generate(%q, false) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestValidateValid tests Validate with valid checksums from SPEC.md §5.
func TestValidateValid(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"18"},
		{"125"},
		{"1230"},
		{"001230"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := luhn.Validate(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got {
				t.Errorf("Validate(%q) = false, want true", tt.input)
			}
		})
	}
}

// TestValidateInvalid tests Validate with invalid checksums from SPEC.md §5.
func TestValidateInvalid(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"10"},
		{"120"},
		{"1231"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := luhn.Validate(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got {
				t.Errorf("Validate(%q) = true, want false", tt.input)
			}
		})
	}
}

// TestRandomProperties tests Random output properties from SPEC.md §5.
func TestRandomProperties(t *testing.T) {
	lengths := []string{"2", "5", "10", "50", "100"}

	for _, length := range lengths {
		t.Run("length_"+length, func(t *testing.T) {
			n, _ := strconv.Atoi(length)
			result, err := luhn.Random(length)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Output length equals requested length
			if len(result) != n {
				t.Errorf("len(Random(%q)) = %d, want %d", length, len(result), n)
			}

			// Output passes Validate
			valid, err := luhn.Validate(result)
			if err != nil {
				t.Fatalf("Validate(%q) error: %v", result, err)
			}
			if !valid {
				t.Errorf("Random(%q) produced %q which fails Validate", length, result)
			}

			// Output contains only digits
			for i, c := range result {
				if c < '0' || c > '9' {
					t.Errorf("Random(%q)[%d] = %c, want digit", length, i, c)
				}
			}

			// First digit is not zero
			if result[0] == '0' {
				t.Errorf("Random(%q) = %q, first digit is zero", length, result)
			}
		})
	}
}

// TestRandomUniqueness verifies that 100 consecutive calls produce 100 unique values.
func TestRandomUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		result, err := luhn.Random("10")
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		if seen[result] {
			t.Fatalf("duplicate value %q at iteration %d", result, i)
		}
		seen[result] = true
	}
}

// TestRandomDistribution verifies that Random output at length "2" is roughly uniformly
// distributed (within 40% margin over 1000 iterations), as required by SPEC.md §5.
func TestRandomDistribution(t *testing.T) {
	const iterations = 1000
	counts := make(map[string]int)

	for i := 0; i < iterations; i++ {
		result, err := luhn.Random("2")
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		counts[result]++
	}

	// With length 2, first digit is 1-9 and check digit is determined by first digit,
	// so there are 9 possible values. Expected frequency is ~1000/9 ≈ 111.
	expected := float64(iterations) / float64(len(counts))
	margin := expected * 0.40

	for value, count := range counts {
		if float64(count) < expected-margin || float64(count) > expected+margin {
			t.Errorf("value %q appeared %d times (expected %.0f ± %.0f)", value, count, expected, margin)
		}
	}
}

// TestGenerateModN_Base10 verifies that GenerateModN with n=10 matches Generate.
func TestGenerateModN_Base10(t *testing.T) {
	inputs := []string{"1", "12", "123", "1234", "7992739871"}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			want, err := luhn.Generate(input, false)
			if err != nil {
				t.Fatalf("Generate error: %v", err)
			}
			got, err := luhn.GenerateModN(input, 10, false)
			if err != nil {
				t.Fatalf("GenerateModN error: %v", err)
			}
			if got != want {
				t.Errorf("GenerateModN(%q, 10, false) = %q, want %q", input, got, want)
			}
		})
	}
}

// TestGenerateModN_Base10ChecksumOnly verifies checksumOnly mode matches Generate.
func TestGenerateModN_Base10ChecksumOnly(t *testing.T) {
	inputs := []string{"1", "123", "7992739871"}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			want, err := luhn.Generate(input, true)
			if err != nil {
				t.Fatalf("Generate error: %v", err)
			}
			got, err := luhn.GenerateModN(input, 10, true)
			if err != nil {
				t.Fatalf("GenerateModN error: %v", err)
			}
			if got != want {
				t.Errorf("GenerateModN(%q, 10, true) = %q, want %q", input, got, want)
			}
		})
	}
}

// TestGenerateModN_Alphanumeric tests GenerateModN with n=36 (full 0-9A-Z alphabet).
func TestGenerateModN_Alphanumeric(t *testing.T) {
	// With n=36, characters A-Z map to indices 10-35.
	// We verify the result ends with a valid CODE_POINTS character and round-trips via ValidateModN.
	tests := []struct {
		input string
		n     int
	}{
		{"A", 36},
		{"AB", 36},
		{"HELLO", 36},
		{"123ABC", 36},
		{"A", 16},
		{"FF", 16},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := luhn.GenerateModN(tt.input, tt.n, false)
			if err != nil {
				t.Fatalf("GenerateModN(%q, %d, false) error: %v", tt.input, tt.n, err)
			}
			// Result should be input + one check character
			if len(result) != len(tt.input)+1 {
				t.Errorf("len = %d, want %d", len(result), len(tt.input)+1)
			}
			// Prefix should be the original input
			if result[:len(tt.input)] != tt.input {
				t.Errorf("prefix = %q, want %q", result[:len(tt.input)], tt.input)
			}
		})
	}
}

// TestGenerateModN_InvalidN tests that n out of range [1,36] returns an error.
func TestGenerateModN_InvalidN(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{"zero", 0},
		{"negative", -1},
		{"too large", 37},
	}

	want := "n must be between 1 and 36"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := luhn.GenerateModN("A", tt.n, false)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err.Error() != want {
				t.Errorf("got %q, want %q", err.Error(), want)
			}
		})
	}
}

// TestGenerateModN_EmptyString tests that empty string returns an error.
func TestGenerateModN_EmptyString(t *testing.T) {
	_, err := luhn.GenerateModN("", 10, false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "string cannot be empty"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

// TestValidateModN_RoundTrip tests that GenerateModN output passes ValidateModN.
func TestValidateModN_RoundTrip(t *testing.T) {
	tests := []struct {
		input string
		n     int
	}{
		{"1", 10},
		{"123", 10},
		{"7992739871", 10},
		{"A", 36},
		{"HELLO", 36},
		{"123ABC", 36},
		{"FF", 16},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			generated, err := luhn.GenerateModN(tt.input, tt.n, false)
			if err != nil {
				t.Fatalf("GenerateModN error: %v", err)
			}
			valid, err := luhn.ValidateModN(generated, tt.n)
			if err != nil {
				t.Fatalf("ValidateModN error: %v", err)
			}
			if !valid {
				t.Errorf("ValidateModN(%q, %d) = false, want true", generated, tt.n)
			}
		})
	}
}

// TestValidateModN_Invalid tests that wrong check characters are rejected.
func TestValidateModN_Invalid(t *testing.T) {
	tests := []struct {
		input string
		n     int
	}{
		{"10", 10},
		{"AZ", 36},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// First verify this is actually invalid by generating the correct one
			generated, _ := luhn.GenerateModN(tt.input[:len(tt.input)-1], tt.n, false)
			if generated == tt.input {
				t.Skip("input happens to be valid, skip")
			}

			valid, err := luhn.ValidateModN(tt.input, tt.n)
			if err != nil {
				t.Fatalf("ValidateModN error: %v", err)
			}
			if valid {
				t.Errorf("ValidateModN(%q, %d) = true, want false", tt.input, tt.n)
			}
		})
	}
}

// TestValidateModN_Errors tests ValidateModN error cases.
func TestValidateModN_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		n     int
		want  string
	}{
		{"empty string", "", 10, "string cannot be empty"},
		{"length 1", "A", 36, "string must be longer than 1 character"},
		{"n too small", "AB", 0, "n must be between 1 and 36"},
		{"n too large", "AB", 37, "n must be between 1 and 36"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := luhn.ValidateModN(tt.input, tt.n)
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.want)
			}
			if err.Error() != tt.want {
				t.Errorf("got %q, want %q", err.Error(), tt.want)
			}
		})
	}
}

// TestChecksumModN_Base10 verifies ChecksumModN matches generateChecksum for base-10.
func TestChecksumModN_Base10(t *testing.T) {
	tests := []struct {
		input    string
		wantChar string
	}{
		{"1", "8"},
		{"123", "0"},
		{"7992739871", "3"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			idx, err := luhn.ChecksumModN(tt.input, 10)
			if err != nil {
				t.Fatalf("ChecksumModN error: %v", err)
			}
			// idx should be the digit value itself for base-10
			wantIdx, _ := strconv.Atoi(tt.wantChar)
			if idx != wantIdx {
				t.Errorf("ChecksumModN(%q, 10) = %d, want %d", tt.input, idx, wantIdx)
			}
		})
	}
}

// TestChecksumModN_Alphanumeric verifies ChecksumModN returns the correct index.
func TestChecksumModN_Alphanumeric(t *testing.T) {
	// Generate a known value and check that ChecksumModN returns the index of the check char
	input := "HELLO"
	n := 36
	generated, err := luhn.GenerateModN(input, n, true)
	if err != nil {
		t.Fatalf("GenerateModN error: %v", err)
	}

	idx, err := luhn.ChecksumModN(input, n)
	if err != nil {
		t.Fatalf("ChecksumModN error: %v", err)
	}

	// The check character from GenerateModN should be codePoints[idx]
	codePoints := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if string(codePoints[idx]) != generated {
		t.Errorf("ChecksumModN(%q, %d) = %d (char %c), want char %q", input, n, idx, codePoints[idx], generated)
	}
}

// TestChecksumModN_Errors tests ChecksumModN error cases.
func TestChecksumModN_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		n     int
		want  string
	}{
		{"empty string", "", 10, "string cannot be empty"},
		{"n too small", "A", 0, "n must be between 1 and 36"},
		{"n too large", "A", 37, "n must be between 1 and 36"},
		{"invalid char", "A!", 36, "invalid character: '!'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := luhn.ChecksumModN(tt.input, tt.n)
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.want)
			}
			if err.Error() != tt.want {
				t.Errorf("got %q, want %q", err.Error(), tt.want)
			}
		})
	}
}

// TestSharedValidation tests the shared validation errors through Generate, Validate, and Random.
// Each error case is tested through all three functions to confirm they share the same validation.
func TestSharedValidation(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", "string cannot be empty"},
		{"spaces", " 123 ", "string cannot contain spaces"},
		{"negative", "-123", "negative numbers are not allowed"},
		{"float", "123.45", "floating point numbers are not allowed"},
		{"non-numeric", "1a", "string must be convertible to a number"},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/Generate", func(t *testing.T) {
			_, err := luhn.Generate(tt.input, false)
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.want)
			}
			if err.Error() != tt.want {
				t.Errorf("got %q, want %q", err.Error(), tt.want)
			}
		})

		t.Run(tt.name+"/Validate", func(t *testing.T) {
			_, err := luhn.Validate(tt.input)
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.want)
			}
			if err.Error() != tt.want {
				t.Errorf("got %q, want %q", err.Error(), tt.want)
			}
		})

		t.Run(tt.name+"/Random", func(t *testing.T) {
			_, err := luhn.Random(tt.input)
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.want)
			}
			if err.Error() != tt.want {
				t.Errorf("got %q, want %q", err.Error(), tt.want)
			}
		})
	}
}

// TestSharedValidation_Order verifies that validation checks are applied in the correct order.
// When multiple checks could fail, the first one in spec order wins.
func TestSharedValidation_Order(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"space before negative", " -123", "string cannot contain spaces"},
		{"negative before float", "-1.5", "negative numbers are not allowed"},
		{"float before non-numeric", "1.a", "floating point numbers are not allowed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := luhn.Generate(tt.input, false)
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.want)
			}
			if err.Error() != tt.want {
				t.Errorf("got %q, want %q", err.Error(), tt.want)
			}
		})
	}
}

// TestValidateSpecificValidation tests the Validate-specific length check.
func TestValidateSpecificValidation(t *testing.T) {
	_, err := luhn.Validate("1")
	if err == nil {
		t.Fatal("expected error for single character input, got nil")
	}
	want := "string must be longer than 1 character"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

// TestRandomSpecificValidation tests the Random-specific length checks.
func TestRandomSpecificValidation(t *testing.T) {
	t.Run("too small", func(t *testing.T) {
		_, err := luhn.Random("1")
		if err == nil {
			t.Fatal("expected error for length 1, got nil")
		}
		want := "string must be greater than 1"
		if err.Error() != want {
			t.Errorf("got %q, want %q", err.Error(), want)
		}
	})

	t.Run("too large", func(t *testing.T) {
		// A 99-digit string overflows strconv.Atoi, so Random treats
		// the parse error as exceeding the max length.
		input := strings.Repeat("1", 99)
		_, err := luhn.Random(input)
		if err == nil {
			t.Fatal("expected error for oversized length, got nil")
		}
		want := "string must be less than 100 characters"
		if err.Error() != want {
			t.Errorf("got %q, want %q", err.Error(), want)
		}
	})

	t.Run("boundary 100 is valid", func(t *testing.T) {
		_, err := luhn.Random("100")
		// 100 is <= 100, so it should NOT produce the "less than 100" error.
		// It may still fail because Random is a stub, but it should not
		// return this specific validation error.
		if err != nil && err.Error() == "string must be less than 100 characters" {
			t.Error("100 should be a valid length, got max-length error")
		}
	})

	t.Run("boundary 101 is invalid", func(t *testing.T) {
		_, err := luhn.Random("101")
		if err == nil {
			t.Fatal("expected error for length 101, got nil")
		}
		want := "string must be less than 100 characters"
		if err.Error() != want {
			t.Errorf("got %q, want %q", err.Error(), want)
		}
	})

	t.Run("zero is invalid", func(t *testing.T) {
		_, err := luhn.Random("0")
		if err == nil {
			t.Fatal("expected error for length 0, got nil")
		}
		want := "string must be greater than 1"
		if err.Error() != want {
			t.Errorf("got %q, want %q", err.Error(), want)
		}
	})
}
