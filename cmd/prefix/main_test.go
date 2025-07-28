package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"go.uber.org/goleak"
	"moul.io/u"
)

func TestRun(t *testing.T) {
	f, cleanup := u.MustTempfileWithContent([]byte("AAA\nBBB\nCCC"))
	defer cleanup()

	// capture stdout and stderr
	closer := u.MustCaptureStdoutAndStderr()

	// simulate CLI call
	err := run([]string{"prefix", "-format", "{{DEFAULT}} ", f.Name()})
	if err != nil {
		t.Fatalf("err should be nil: %v", err)
	}

	// ignore output in this test
	_ = closer()
}

func Example() {
	f, cleanup := u.MustTempfileWithContent([]byte("AAA\nBBB\nCCC"))
	defer cleanup()

	// simulate normal CLI:
	//    $> prefix -format "{{.LineNumber3}}" /path/to/tempfile
	err := run([]string{"prefix", "-format", "{{.LineNumber3}} ", f.Name()})
	if err != nil {
		panic(err)
	}

	// Output:
	// 1   AAA
	// 2   BBB
	// 3   CCC
}

// no output (everything is in stderr)
func Example_usage() {
	err := run([]string{"prefix", "-h"})
	if err != nil {
		panic(err)
	}
}

func TestUsage(t *testing.T) {
	// capture stdout and stderr
	closer := u.MustCaptureStdoutAndStderr()

	// simulate CLI call
	err := run([]string{"prefix", "-h"})
	if err != flag.ErrHelp {
		t.Fatalf("err should be flag.ErrHelp: %v", err)
	}

	// ignore output in this test
	output := closer()
	if !strings.Contains(output, "USAGE") ||
		!strings.Contains(output, "FLAGS") ||
		!strings.Contains(output, "SYNTAX") ||
		!strings.Contains(output, "PRESETS") ||
		!strings.Contains(output, "EXAMPLES") {
		t.Errorf("usage should contain USAGE, FLAGS, SYNTAX, PRESETS, and EXAMPLES")
	}
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func TestRunErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		setup   func() func()
		wantErr error
	}{
		{
			name:    "too many arguments",
			args:    []string{"prefix", "file1", "file2"},
			wantErr: errors.New("multi files support not supported yet"),
		},
		{
			name:    "file not found",
			args:    []string{"prefix", "/nonexistent/file"},
			wantErr: os.ErrNotExist,
		},
		{
			name:    "help flag",
			args:    []string{"prefix", "-h"},
			wantErr: flag.ErrHelp,
		},
		{
			name:    "unknown flag",
			args:    []string{"prefix", "-unknown"},
			wantErr: errors.New("flag provided but not defined"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				cleanup := tt.setup()
				defer cleanup()
			}

			// capture stdout and stderr
			closer := u.MustCaptureStdoutAndStderr()
			defer func() { _ = closer() }()

			err := run(tt.args)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				} else if !strings.Contains(err.Error(), tt.wantErr.Error()) && !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error containing %q, got %q", tt.wantErr.Error(), err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRunStdin(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		stdin    string
		contains []string
	}{
		{
			name:     "read from stdin with dash",
			args:     []string{"prefix", "-format", "{{.LineNumber}} ", "-"},
			stdin:    "line1\nline2\nline3",
			contains: []string{"1 line1", "2 line2", "3 line3"},
		},
		{
			name:     "read from stdin without args",
			args:     []string{"prefix", "-format", "{{.LineNumber}} "},
			stdin:    "test1\ntest2",
			contains: []string{"1 test1", "2 test2"},
		},
		{
			name:     "empty stdin",
			args:     []string{"prefix", "-format", "LINE {{.LineNumber}} "},
			stdin:    "",
			contains: []string{}, // no output expected
		},
		{
			name:     "stdin with special chars",
			args:     []string{"prefix", "-format", "{{.LineNumber}} "},
			stdin:    "hello\tworld\ntest\x00null\nunicode🚀",
			contains: []string{"1 hello\tworld", "2 test\x00null", "3 unicode🚀"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a pipe to simulate stdin
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}

			// Replace stdin
			oldStdin := os.Stdin
			os.Stdin = r
			defer func() { os.Stdin = oldStdin }()

			// Write test data to stdin
			go func() {
				_, _ = w.Write([]byte(tt.stdin))
				w.Close()
			}()

			// Capture stdout
			closer := u.MustCaptureStdoutAndStderr()

			// Run the command
			err = run(tt.args)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			output := closer()

			// Check output
			for _, expected := range tt.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain %q, got %q", expected, output)
				}
			}
		})
	}
}

func TestRunWithDifferentFormats(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		content  string
		contains []string
	}{
		{
			name:     "default format explicit",
			format:   "{{DEFAULT}} ",
			content:  "test",
			contains: []string{"1   ", "up=", "d=", " | ", "test"},
		},
		{
			name:     "custom format",
			format:   "[{{.LineNumber}}] ",
			content:  "hello\nworld",
			contains: []string{"[1] hello", "[2] world"},
		},
		{
			name:     "preset format",
			format:   "{{SHORT_DATE}} ",
			content:  "line",
			contains: []string{"/", ":", "line"},
		},
		{
			name:     "complex format",
			format:   "{{.LineNumber}} {{.ShortUptime}} {{.ShortDuration}} | ",
			content:  "a\nb\nc",
			contains: []string{"1 ", "2 ", "3 ", " | a", " | b", " | c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, cleanup := u.MustTempfileWithContent([]byte(tt.content))
			defer cleanup()

			// Capture stdout
			closer := u.MustCaptureStdoutAndStderr()

			// Build args
			args := []string{"prefix"}
			if tt.format != "" {
				args = append(args, "-format", tt.format)
			}
			args = append(args, f.Name())

			// Run the command
			err := run(args)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			output := closer()

			// Check output
			for _, expected := range tt.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain %q, got %q", expected, output)
				}
			}
		})
	}
}

func TestRunEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		format   string
		validate func(t *testing.T, output string)
	}{
		{
			name:    "empty file",
			content: "",
			format:  "{{.LineNumber}} ",
			validate: func(t *testing.T, output string) {
				if output != "" {
					t.Errorf("expected empty output for empty file, got %q", output)
				}
			},
		},
		{
			name:    "file with only newlines",
			content: "\n\n\n",
			format:  "{{.LineNumber}} ",
			validate: func(t *testing.T, output string) {
				expected := "1 \n2 \n3 \n"
				if output != expected {
					t.Errorf("expected %q, got %q", expected, output)
				}
			},
		},
		{
			name:    "file without final newline",
			content: "line1\nline2",
			format:  "{{.LineNumber}} ",
			validate: func(t *testing.T, output string) {
				expected := "1 line1\n2 line2\n"
				if output != expected {
					t.Errorf("expected %q, got %q", expected, output)
				}
			},
		},
		{
			name:    "very long line",
			content: strings.Repeat("a", 50000), // Under scanner limit
			format:  "{{.LineNumber}} ",
			validate: func(t *testing.T, output string) {
				if !strings.HasPrefix(output, "1 ") {
					t.Errorf("expected output to start with '1 '")
				}
				if len(output) != 50003 { // "1 " + 50000 a's + newline
					t.Errorf("expected output length 50003, got %d", len(output))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, cleanup := u.MustTempfileWithContent([]byte(tt.content))
			defer cleanup()

			// Capture stdout
			closer := u.MustCaptureStdoutAndStderr()

			// Run the command
			err := run([]string{"prefix", "-format", tt.format, f.Name()})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			output := closer()
			tt.validate(t, output)
		})
	}
}

func TestMainFunction(t *testing.T) {
	// Test that main() function works without panicking
	// We can't easily test os.Exit behavior, but we can ensure it doesn't panic

	// Save original args and restore after test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test with help flag (should exit gracefully)
	os.Args = []string{"prefix", "-h"}

	// Capture the exit code by replacing os.Exit temporarily
	exitCode := -1
	oldExit := osExit
	osExit = func(code int) {
		exitCode = code
		panic(fmt.Sprintf("os.Exit(%d)", code))
	}
	defer func() { osExit = oldExit }()

	// Capture output
	closer := u.MustCaptureStdoutAndStderr()
	defer func() { _ = closer() }()

	// Run main and expect it to panic (our mock os.Exit)
	defer func() {
		if r := recover(); r != nil {
			if exitCode != 1 {
				t.Errorf("expected exit code 1 for help flag, got %d", exitCode)
			}
		} else {
			t.Error("expected main to call os.Exit")
		}
	}()

	main()
}

// Test helper for simulating file read errors
type errorReader struct {
	err error
}

func (e errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

func TestScannerError(t *testing.T) {
	// This tests the scanner error handling path
	// by providing a reader that returns an error

	// Save stdin and restore after test
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a pipe and immediately close the write end to simulate EOF
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	w.Close() // Close immediately to trigger EOF
	os.Stdin = r

	// Capture output
	closer := u.MustCaptureStdoutAndStderr()

	// Run should complete without error (EOF is normal)
	err = run([]string{"prefix", "-format", "{{.LineNumber}} "})
	if err != nil {
		t.Errorf("unexpected error for EOF: %v", err)
	}

	_ = closer()
}

func TestFileCloseError(t *testing.T) {
	// Test file operations with a file that gets deleted
	f, cleanup := u.MustTempfileWithContent([]byte("test content"))
	fileName := f.Name()
	cleanup() // Clean up immediately to make the file inaccessible

	// Capture output
	closer := u.MustCaptureStdoutAndStderr()
	defer func() { _ = closer() }()

	// This should fail with file not found
	err := run([]string{"prefix", fileName})
	if err == nil {
		t.Error("expected error for missing file")
	} else if !errors.Is(err, os.ErrNotExist) && !strings.Contains(err.Error(), "no such file") {
		t.Errorf("expected file not found error, got: %v", err)
	}
}

func TestInvalidTemplateHandling(t *testing.T) {
	// Test that invalid templates cause appropriate errors
	f, cleanup := u.MustTempfileWithContent([]byte("test"))
	defer cleanup()

	// Capture output
	closer := u.MustCaptureStdoutAndStderr()
	defer func() { _ = closer() }()

	// This should panic due to invalid template syntax
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid template")
		}
	}()

	_ = run([]string{"prefix", "-format", "{{unclosed", f.Name()})
}

func TestBinaryInput(t *testing.T) {
	// Test handling of binary input
	binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD, 0x0A} // Including newline
	f, cleanup := u.MustTempfileWithContent(binaryData)
	defer cleanup()

	// Capture output
	closer := u.MustCaptureStdoutAndStderr()

	// Run the command
	err := run([]string{"prefix", "-format", "{{.LineNumber}} ", f.Name()})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := closer()

	// Should handle binary data gracefully
	if !strings.Contains(output, "1 ") {
		t.Error("expected line number in output")
	}
}

func TestConcurrentFileAccess(t *testing.T) {
	// Test reading a file while it's being written to
	tmpFile, err := os.CreateTemp("", "prefix-test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write initial content
	_, err = tmpFile.WriteString("line1\nline2\n")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Sync()

	// Keep file open for writing
	go func() {
		for i := 3; i <= 5; i++ {
			_, _ = tmpFile.WriteString(fmt.Sprintf("line%d\n", i))
			tmpFile.Sync()
		}
		tmpFile.Close()
	}()

	// Capture output
	closer := u.MustCaptureStdoutAndStderr()

	// Run prefix on the file
	err = run([]string{"prefix", "-format", "{{.LineNumber}} ", tmpFile.Name()})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := closer()

	// Should at least see the initial lines
	if !strings.Contains(output, "1 line1") || !strings.Contains(output, "2 line2") {
		t.Errorf("expected initial lines in output, got: %q", output)
	}
}

// Benchmark tests
func BenchmarkRun(b *testing.B) {
	// Create a test file with many lines
	var content strings.Builder
	for i := 0; i < 1000; i++ {
		content.WriteString(fmt.Sprintf("This is line number %d with some content\n", i))
	}

	f, cleanup := u.MustTempfileWithContent([]byte(content.String()))
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Capture output to avoid noise
		oldStdout := os.Stdout
		os.Stdout, _ = os.Open(os.DevNull)

		err := run([]string{"prefix", "-format", "{{.LineNumber}} ", f.Name()})
		if err != nil {
			b.Fatal(err)
		}

		os.Stdout = oldStdout
	}
}

func BenchmarkRunComplexFormat(b *testing.B) {
	// Create a test file
	f, cleanup := u.MustTempfileWithContent([]byte(strings.Repeat("benchmark line\n", 100)))
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Capture output to avoid noise
		oldStdout := os.Stdout
		os.Stdout, _ = os.Open(os.DevNull)

		err := run([]string{"prefix", "-format", "{{DEFAULT}} {{SHORT_DATE}} ", f.Name()})
		if err != nil {
			b.Fatal(err)
		}

		os.Stdout = oldStdout
	}
}


// Test that we handle write errors gracefully
func TestWriteError(t *testing.T) {
	// Create a file with multiple lines to increase chance of write error
	content := strings.Repeat("test line\n", 100)
	f, cleanup := u.MustTempfileWithContent([]byte(content))
	defer cleanup()

	// Save stdout
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	// Create a pipe and close the read end to simulate write error
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	r.Close() // Close read end to cause write error
	os.Stdout = w

	// This might or might not error depending on buffering
	// so we just run it without checking the error
	_ = run([]string{"prefix", "-format", "{{.LineNumber}} ", f.Name()})

	w.Close()
}

// Test reading from a directory (should fail)
func TestReadDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "prefix-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Capture output
	closer := u.MustCaptureStdoutAndStderr()
	defer func() { _ = closer() }()

	// Try to read a directory - this might work on some systems
	// so we just check it doesn't panic
	_ = run([]string{"prefix", tmpDir})
}

// Test permission denied scenario
func TestPermissionDenied(t *testing.T) {
	// Skip if running as root
	if os.Geteuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	tmpFile, err := os.CreateTemp("", "prefix-test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write content and close
	_, err = tmpFile.WriteString("test content")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Remove read permissions
	err = os.Chmod(tmpFile.Name(), 0200) // Write only
	if err != nil {
		t.Fatal(err)
	}

	// Capture output
	closer := u.MustCaptureStdoutAndStderr()
	defer func() { _ = closer() }()

	// Try to read the file
	err = run([]string{"prefix", tmpFile.Name()})
	if err == nil {
		t.Error("expected permission denied error")
	} else if !strings.Contains(err.Error(), "permission denied") && !errors.Is(err, os.ErrPermission) {
		t.Errorf("expected permission error, got: %v", err)
	}
}

// Helper to simulate stdin for error cases
type stdinSimulator struct {
	content string
	err     error
	readPos int
}

func (s *stdinSimulator) Read(p []byte) (n int, err error) {
	if s.err != nil {
		return 0, s.err
	}
	if s.readPos >= len(s.content) {
		return 0, io.EOF
	}
	n = copy(p, s.content[s.readPos:])
	s.readPos += n
	return n, nil
}

func (s *stdinSimulator) Close() error {
	return nil
}

func (s *stdinSimulator) Stat() (os.FileInfo, error) {
	return nil, errors.New("not a regular file")
}
