package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	old := output
	output = &buf
	defer func() { output = old }()
	f()
	return buf.String()
}

type fakeRoundTripper struct{}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("ok")),
	}, nil
}

func TestHandleList(t *testing.T) {
	http.DefaultClient.Transport = &fakeRoundTripper{}

	out := captureOutput(func() { handleList("whitelist", "add", "") })
	if !strings.Contains(out, "CIDR is required") {
		t.Errorf("expected 'CIDR is required', got: %s", out)
	}

	out = captureOutput(func() { handleList("whitelist", "add", "192.168.1.0/24") })
	if !strings.Contains(out, "Status: 200 OK") {
		t.Errorf("expected 'Status: 200 OK', got: %s", out)
	}
}

func TestHandleBucket(t *testing.T) {
	http.DefaultClient.Transport = &fakeRoundTripper{}

	out := captureOutput(func() { handleBucket("wrong") })
	if !strings.Contains(out, "only reset supported") {
		t.Errorf("expected 'only reset supported', got: %s", out)
	}

	out = captureOutput(func() { handleBucket("reset") })
	if !strings.Contains(out, "Status: 200 OK") {
		t.Errorf("expected 'Status: 200 OK', got: %s", out)
	}
}

func TestUsage(t *testing.T) {
	out := captureOutput(func() { usage() })
	if !strings.Contains(out, "Usage:") {
		t.Errorf("expected 'Usage:' in output, got: %s", out)
	}
}

func TestMainCases(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	http.DefaultClient.Transport = &fakeRoundTripper{}

	os.Args = []string{"cli"}
	out := captureOutput(main)
	if !strings.Contains(out, "Usage:") {
		t.Errorf("expected usage output, got: %s", out)
	}

	os.Args = []string{"cli", "unknown"}
	out = captureOutput(main)
	if !strings.Contains(out, "Usage:") {
		t.Errorf("expected usage output, got: %s", out)
	}

	os.Args = []string{"cli", "whitelist", "add", "192.168.1.0/24"}
	out = captureOutput(main)
	if !strings.Contains(out, "Status: 200 OK") {
		t.Errorf("expected 'Status: 200 OK', got: %s", out)
	}

	os.Args = []string{"cli", "blacklist", "remove", "10.0.0.0/8"}
	out = captureOutput(main)
	if !strings.Contains(out, "Status: 200 OK") {
		t.Errorf("expected 'Status: 200 OK', got: %s", out)
	}

	os.Args = []string{"cli", "bucket", "reset"}
	out = captureOutput(main)
	if !strings.Contains(out, "Status: 200 OK") {
		t.Errorf("expected 'Status: 200 OK', got: %s", out)
	}
}
