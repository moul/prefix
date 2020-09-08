package main

import (
	"io/ioutil"
	"os"
	"testing"

	"go.uber.org/goleak"
	"moul.io/u"
)

func TestRun(t *testing.T) {
	// create temp file with some content
	var tmpfilePath string
	{
		content := []byte("AAA\nBBB\nCCC")
		tmpfile, err := ioutil.TempFile("", "example")
		checkErr(err)
		_, err = tmpfile.Write(content)
		checkErr(err)
		err = tmpfile.Close()
		checkErr(err)
		defer os.Remove(tmpfile.Name())
		tmpfilePath = tmpfile.Name()
	}

	// capture stdout and stderr
	closer := u.MustCaptureStdoutAndStderr()

	// simulate CLI call
	err := run([]string{"prefix", tmpfilePath})
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
	//    $> prefix /path/to/tempfile
	err := run([]string{"prefix", f.Name()})
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
