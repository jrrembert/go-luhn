package luhn_test

import (
	"strings"
	"testing"

	luhn "github.com/jrrembert/go-luhn"
)

// TestPackageExports verifies the package compiles and all 6 public functions are exported.
func TestPackageExports(t *testing.T) {
	// Generate
	if _, err := luhn.Generate("1", false); err != nil {
		t.Logf("Generate returned error (stub not yet implemented): %v", err)
	}

	// Validate
	if _, err := luhn.Validate("18"); err != nil {
		t.Logf("Validate returned error (stub not yet implemented): %v", err)
	}

	// Random
	if _, err := luhn.Random("4"); err != nil {
		t.Logf("Random returned error (stub not yet implemented): %v", err)
	}

	// GenerateModN
	if _, err := luhn.GenerateModN("1", 10, false); err != nil {
		t.Logf("GenerateModN returned error (stub not yet implemented): %v", err)
	}

	// ValidateModN
	if _, err := luhn.ValidateModN("18", 10); err != nil {
		t.Logf("ValidateModN returned error (stub not yet implemented): %v", err)
	}

	// ChecksumModN
	if _, err := luhn.ChecksumModN("1", 10); err != nil {
		t.Logf("ChecksumModN returned error (stub not yet implemented): %v", err)
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
