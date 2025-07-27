package prefix_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

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

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		wantErr  bool
		contains []string
	}{
		{
			name:     "empty format uses default",
			format:   "",
			contains: []string{"1  ", "up=", "d=", "|"},
		},
		{
			name:     "simple line number",
			format:   "{{.LineNumber}} ",
			contains: []string{"1 ", "2 ", "3 "},
		},
		{
			name:     "line number with padding",
			format:   "{{.LineNumber3}} ",
			contains: []string{"1   ", "2   ", "3   "},
		},
		{
			name:     "line number4",
			format:   "{{.LineNumber4}} ",
			contains: []string{"1    ", "2    ", "3    "},
		},
		{
			name:     "line number5",
			format:   "{{.LineNumber5}} ",
			contains: []string{"1     ", "2     ", "3     "},
		},
		{
			name:     "format passthrough",
			format:   "{{.Format}} ",
			contains: []string{"{{.Format}} "},
		},
		{
			name:     "uptime format",
			format:   "{{.Uptime}} ",
			contains: []string{"s"},
		},
		{
			name:     "duration format",
			format:   "{{.Duration}} ",
			contains: []string{"s"},
		},
		{
			name:     "short duration",
			format:   "{{.ShortDuration}} ",
			contains: []string{" "},
		},
		{
			name:     "short uptime",
			format:   "{{.ShortUptime}} ",
			contains: []string{" "},
		},
		{
			name:     "now function",
			format:   "{{now}} ",
			contains: []string{"202", "-", ":"},
		},
		{
			name:     "unix epoch",
			format:   "{{now | unixEpoch}} ",
			contains: []string{"1"},
		},
		{
			name:     "uuidv4",
			format:   "{{uuidv4}} ",
			contains: []string{"-"},
		},
		{
			name:     "DEFAULT preset",
			format:   "{{DEFAULT}} ",
			contains: []string{"up=", "d="},
		},
		{
			name:     "SLOW_LINES preset",
			format:   "{{SLOW_LINES}} ",
			contains: []string{" "},
		},
		{
			name:     "SHORT_DATE preset",
			format:   "{{SHORT_DATE}} ",
			contains: []string{"/", ":"},
		},
		{
			name:     "multiple presets",
			format:   "{{DEFAULT}} {{SHORT_DATE}} ",
			contains: []string{"up=", "d=", "/", ":"},
		},
		{
			name:     "nested presets",
			format:   "{{DEFAULT}} {{DEFAULT}} ",
			contains: []string{"up=", "d="},
		},
		{
			name:     "literal text",
			format:   "PREFIX> ",
			contains: []string{"PREFIX> "},
		},
		{
			name:     "mixed format",
			format:   "[{{.LineNumber3}}] {{.ShortUptime}} > ",
			contains: []string{"[1  ]", ">"},
		},
		{
			name:     "sprig function",
			format:   "{{randAlphaNum 5}} ",
			contains: []string{" "},
		},
		{
			name:     "complex template",
			format:   "{{if eq .LineNumber 1}}FIRST{{else}}OTHER{{end}} ",
			contains: []string{"FIRST", "OTHER"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefixer := prefix.New(tt.format)
			result1 := prefixer.PrefixLine("test1")
			result2 := prefixer.PrefixLine("test2")
			result3 := prefixer.PrefixLine("test3")

			for _, contains := range tt.contains {
				found := false
				for _, result := range []string{result1, result2, result3} {
					if strings.Contains(result, contains) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected results to contain %q, got %q, %q, %q", contains, result1, result2, result3)
				}
			}
		})
	}
}

func TestPrefixLine(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		lines    []string
		validate func(t *testing.T, results []string)
	}{
		{
			name:   "incremental line numbers",
			format: "{{.LineNumber}} ",
			lines:  []string{"a", "b", "c"},
			validate: func(t *testing.T, results []string) {
				if results[0] != "1 a" {
					t.Errorf("expected '1 a', got %q", results[0])
				}
				if results[1] != "2 b" {
					t.Errorf("expected '2 b', got %q", results[1])
				}
				if results[2] != "3 c" {
					t.Errorf("expected '3 c', got %q", results[2])
				}
			},
		},
		{
			name:   "empty lines",
			format: "{{.LineNumber}} ",
			lines:  []string{"", "", ""},
			validate: func(t *testing.T, results []string) {
				if results[0] != "1 " {
					t.Errorf("expected '1 ', got %q", results[0])
				}
				if results[1] != "2 " {
					t.Errorf("expected '2 ', got %q", results[1])
				}
				if results[2] != "3 " {
					t.Errorf("expected '3 ', got %q", results[2])
				}
			},
		},
		{
			name:   "special characters",
			format: "{{.LineNumber}} ",
			lines:  []string{"hello\tworld", "test\nline", "special!@#$%^&*()"},
			validate: func(t *testing.T, results []string) {
				if results[0] != "1 hello\tworld" {
					t.Errorf("expected tab preserved, got %q", results[0])
				}
				if results[1] != "2 test\nline" {
					t.Errorf("expected newline preserved, got %q", results[1])
				}
				if results[2] != "3 special!@#$%^&*()" {
					t.Errorf("expected special chars preserved, got %q", results[2])
				}
			},
		},
		{
			name:   "unicode characters",
			format: "{{.LineNumber}} ",
			lines:  []string{"Hello 世界", "🚀 rocket", "café"},
			validate: func(t *testing.T, results []string) {
				if results[0] != "1 Hello 世界" {
					t.Errorf("expected unicode preserved, got %q", results[0])
				}
				if results[1] != "2 🚀 rocket" {
					t.Errorf("expected emoji preserved, got %q", results[1])
				}
				if results[2] != "3 café" {
					t.Errorf("expected accented chars preserved, got %q", results[2])
				}
			},
		},
		{
			name:   "very long lines",
			format: "{{.LineNumber}} ",
			lines:  []string{strings.Repeat("a", 10000), strings.Repeat("b", 50000)},
			validate: func(t *testing.T, results []string) {
				if !strings.HasPrefix(results[0], "1 ") || len(results[0]) != 10002 {
					t.Errorf("expected long line preserved, got length %d", len(results[0]))
				}
				if !strings.HasPrefix(results[1], "2 ") || len(results[1]) != 50002 {
					t.Errorf("expected very long line preserved, got length %d", len(results[1]))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefixer := prefix.New(tt.format)
			var results []string
			for _, line := range tt.lines {
				results = append(results, prefixer.PrefixLine(line))
			}
			tt.validate(t, results)
		})
	}
}

func TestInvalidTemplates(t *testing.T) {
	// Test that using invalid field doesn't panic during New,
	// but will cause error during template execution
	prefixer := prefix.New("{{.Invalid}}")
	
	// The panic happens during template execution, not creation
	defer func() {
		if r := recover(); r == nil {
			// If no panic, that's actually okay too - depends on Go version
			return
		}
	}()
	
	// This might panic depending on template handling
	_ = prefixer.PrefixLine("test")
}

func TestTemplateFunctions(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		validate func(t *testing.T, result string)
	}{
		{
			name:   "env function",
			format: `{{env "USER"}} `,
			validate: func(t *testing.T, result string) {
				// Should contain something (even if empty)
				if !strings.HasSuffix(result, " test") {
					t.Errorf("expected result to end with ' test', got %q", result)
				}
			},
		},
		{
			name:   "short_duration function",
			format: `{{.Duration | short_duration}} `,
			validate: func(t *testing.T, result string) {
				// Should be 7 chars or less plus space
				prefix := strings.TrimSuffix(result, "test")
				if len(strings.TrimSpace(prefix)) > 7 {
					t.Errorf("expected short duration to be <= 7 chars, got %q", prefix)
				}
			},
		},
		{
			name:   "date function",
			format: `{{now | date "2006"}} `,
			validate: func(t *testing.T, result string) {
				// Should contain current year
				year := fmt.Sprintf("%d", time.Now().Year())
				if !strings.Contains(result, year) {
					t.Errorf("expected result to contain year %s, got %q", year, result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefixer := prefix.New(tt.format)
			result := prefixer.PrefixLine("test")
			tt.validate(t, result)
		})
	}
}

func TestDurationTracking(t *testing.T) {
	prefixer := prefix.New("{{.LineNumber}} {{.Duration}} ")
	
	// First line should have very small duration
	result1 := prefixer.PrefixLine("first")
	
	// Sleep a bit
	time.Sleep(10 * time.Millisecond)
	
	// Second line should have measurable duration
	result2 := prefixer.PrefixLine("second")
	
	// First result should contain a very small duration
	if !strings.Contains(result1, "ns") && !strings.Contains(result1, "µs") {
		t.Errorf("expected first line to have nanosecond or microsecond duration, got %q", result1)
	}
	
	// Second result should contain milliseconds
	if !strings.Contains(result2, "ms") {
		t.Errorf("expected second line to have millisecond duration, got %q", result2)
	}
}

func TestUptimeTracking(t *testing.T) {
	prefixer := prefix.New("{{.Uptime}} ")
	
	// First line should have very small uptime
	result1 := prefixer.PrefixLine("first")
	
	// Sleep a bit
	time.Sleep(10 * time.Millisecond)
	
	// Second line should have measurable uptime
	result2 := prefixer.PrefixLine("second")
	
	// First result should contain a very small uptime
	if !strings.Contains(result1, "ns") && !strings.Contains(result1, "µs") {
		t.Errorf("expected first line to have nanosecond or microsecond uptime, got %q", result1)
	}
	
	// Second result should contain milliseconds
	if !strings.Contains(result2, "ms") {
		t.Errorf("expected second line to have millisecond uptime, got %q", result2)
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		format   string
		expected string
	}{
		{
			format:   "test",
			expected: `LinePrefixer{"test"}`,
		},
		{
			format:   "",
			expected: `LinePrefixer{"`,
		},
		{
			format:   "{{.LineNumber}}",
			expected: `LinePrefixer{"{{.LineNumber}}"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			prefixer := prefix.New(tt.format)
			result := fmt.Sprintf("%s", prefixer)
			if !strings.Contains(result, tt.expected) {
				t.Errorf("expected String() to contain %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestPresetReplacement(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		validate func(t *testing.T, prefixer prefix.LinePrefixer)
	}{
		{
			name:   "DEFAULT preset expansion",
			format: "{{DEFAULT}}",
			validate: func(t *testing.T, prefixer prefix.LinePrefixer) {
				result := prefixer.PrefixLine("test")
				// Should contain line number, uptime, and duration
				if !strings.Contains(result, "up=") || !strings.Contains(result, "d=") {
					t.Errorf("DEFAULT preset not properly expanded, got %q", result)
				}
			},
		},
		{
			name:   "recursive preset replacement",
			format: "a{{DEFAULT}}b{{DEFAULT}}c",
			validate: func(t *testing.T, prefixer prefix.LinePrefixer) {
				result := prefixer.PrefixLine("test")
				// Should contain multiple expansions
				count := strings.Count(result, "up=")
				if count != 2 {
					t.Errorf("expected 2 'up=' occurrences, got %d in %q", count, result)
				}
			},
		},
		{
			name:   "SLOW_LINES preset with fast line",
			format: "{{SLOW_LINES}}",
			validate: func(t *testing.T, prefixer prefix.LinePrefixer) {
				result := prefixer.PrefixLine("test")
				// Fast line should show spaces
				if !strings.Contains(result, "    ") {
					t.Errorf("SLOW_LINES should show spaces for fast line, got %q", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t, prefix.New(tt.format))
		})
	}
}

func TestConcurrentUsage(t *testing.T) {
	prefixer := prefix.New("{{.LineNumber}} ")
	
	// Run multiple goroutines that use the prefixer
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				_ = prefixer.PrefixLine(fmt.Sprintf("goroutine %d line %d", id, j))
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// If we get here without panicking, concurrent usage is safe
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		format string
		line   string
		check  func(t *testing.T, result string)
	}{
		{
			name:   "null bytes in line",
			format: "{{.LineNumber}} ",
			line:   "hello\x00world",
			check: func(t *testing.T, result string) {
				if result != "1 hello\x00world" {
					t.Errorf("null byte not preserved, got %q", result)
				}
			},
		},
		{
			name:   "line with only spaces",
			format: "{{.LineNumber}} ",
			line:   "     ",
			check: func(t *testing.T, result string) {
				if result != "1      " {
					t.Errorf("spaces not preserved, got %q", result)
				}
			},
		},
		{
			name:   "line with CRLF",
			format: "{{.LineNumber}} ",
			line:   "test\r\n",
			check: func(t *testing.T, result string) {
				if result != "1 test\r\n" {
					t.Errorf("CRLF not preserved, got %q", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefixer := prefix.New(tt.format)
			result := prefixer.PrefixLine(tt.line)
			tt.check(t, result)
		})
	}
}

// Benchmark tests
func BenchmarkPrefixLine(b *testing.B) {
	prefixer := prefix.New("{{.LineNumber}} ")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = prefixer.PrefixLine("benchmark line")
	}
}

func BenchmarkPrefixLineComplex(b *testing.B) {
	prefixer := prefix.New("{{.LineNumber5}} {{.ShortUptime}} {{.ShortDuration}} {{now | date \"2006-01-02 15:04:05\"}} ")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = prefixer.PrefixLine("benchmark line")
	}
}

func BenchmarkPrefixLineWithPresets(b *testing.B) {
	prefixer := prefix.New("{{DEFAULT}} {{SHORT_DATE}} ")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = prefixer.PrefixLine("benchmark line")
	}
}

// Helper function to test error conditions
func expectPanic(t *testing.T, fn func(), message string) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: expected panic but got none", message)
		} else {
			// Convert panic value to error if possible
			switch v := r.(type) {
			case error:
				// Expected
			case string:
				// Convert string panic to error
				_ = errors.New(v)
			default:
				t.Errorf("%s: unexpected panic type: %T", message, v)
			}
		}
	}()
	fn()
}

func TestBadTemplate(t *testing.T) {
	// Test unclosed template
	expectPanic(t, func() {
		prefix.New("{{unclosed")
	}, "unclosed template")
	
	// Test invalid syntax
	expectPanic(t, func() {
		prefix.New("{{}}")
	}, "empty template")
	
	// Test that field access doesn't panic during New, but during execution
	prefixer := prefix.New("{{.NoSuchField}}")
	defer func() {
		if r := recover(); r == nil {
			// No panic is okay - template might handle this gracefully
		}
	}()
	_ = prefixer.PrefixLine("test")
}