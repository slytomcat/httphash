package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type test struct {
	url  string
	hash string
}

func TestHTTPHashWithMock(t *testing.T) {
	bodyContent := []byte("test body")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(bodyContent)
	}))
	defer ts.Close()
	client := ts.Client()
	received, err := urlHash(client, ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	expected := fmt.Sprintf("%x", md5.Sum(bodyContent))
	if expected != received {
		t.Fatalf("Expected %s, received %s", expected, received)
	}
}

func TestHTTPHash(t *testing.T) {
	// the boby hashes may bchange in future, so this test may fail
	testData := []test{
		{
			url:  "http://www.iana.org/",
			hash: "f3a75a754f40d80eca38b561bb45bfc3",
		},
		{
			url:  "example.com",
			hash: "84238dfc8092e5d9c0dac8ef93371a07",
		},
		{
			url:  "http://ipv6.he.net",
			hash: "00d16b63a5d7a6348aecaf0002407026",
		},
	}
	for _, testCase := range testData {
		res, err := urlHash(http.DefaultClient, testCase.url)
		if err != nil {
			t.Fatal(err)
		}
		if res != testCase.hash {
			t.Fatalf("hash of response from %s is %s while expected %s", testCase.url, res, testCase.hash)
		} else {
			t.Logf("url: %s passed", testCase.url)
		}
	}
}

func captureOutput(f func()) string {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = orig
	return string(out)
}

func TestRun(t *testing.T) {
	bodyContent := []byte("test body")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(bodyContent)
	}))
	defer ts.Close()
	out := captureOutput(func() {
		Run(1, []string{ts.URL})
	})
	expected := fmt.Sprintf("%s %s\n", ts.URL, fmt.Sprintf("%x", md5.Sum(bodyContent)))
	if out != expected {
		t.Fatalf("Expected: %s, received %s", expected, out)
	}
}
