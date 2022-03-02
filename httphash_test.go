package main

import (
	"testing"
)

type test struct {
	url  string
	hash string
}

func TestHTTPHash(t *testing.T) {
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
		res, err := urlHash(testCase.url)
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

func ExampleRun() {
	Run(1, []string{"example.com"})
	// Output: example.com 84238dfc8092e5d9c0dac8ef93371a07
}
