package prefix_test

import (
	"testing"

	"moul.io/prefix"
)

// you can find integration tests in the cmd/prefix/main_test.go file

func TestAvailablePatterns(t *testing.T) {
	// this test is dummy, it just checks that everything runs without panicking
	for pattern := range prefix.AvailablePatterns {
		prefixer := prefix.New(pattern)
		prefixer.PrefixLine("first")
		prefixer.PrefixLine("second")
		prefixer.PrefixLine("third")
	}
}
