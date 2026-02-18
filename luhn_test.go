package luhn_test

import (
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
