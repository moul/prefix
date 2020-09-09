package main

import (
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

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
