package main

import (
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
	if err != nil {
		t.Fatalf("err should be nil: %v", err)
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
