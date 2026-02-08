package services

import (
	"io"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Silence stdout/stderr during tests to reduce noisy flog/println output.
	oldOut := os.Stdout
	oldErr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	code := m.Run()

	_ = w.Close()
	io.Copy(io.Discard, r)
	os.Stdout = oldOut
	os.Stderr = oldErr

	os.Exit(code)
}
