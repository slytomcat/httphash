package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	// HTTPprefix is prefix that have to be in the beggining of URL
	HTTPprefix     = "http://"
	requestTimeout = 5 * time.Second
	urlsChanSize   = 100
)

func urlHash(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	if !strings.HasPrefix(url, HTTPprefix) {
		url = fmt.Sprintf("%s%s", HTTPprefix, url)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(body)), nil
}

func main() {
	var routines int
	flag.IntVar(&routines, "parallel", 10, "number of parallel requests")
	flag.Parse()
	Run(routines, flag.Args())
}

// Run is a main() handler
func Run(routines int, args []string) {
	urls := make(chan string, urlsChanSize)
	wg := sync.WaitGroup{}
	wg.Add(routines)
	for i := 0; i < routines; i++ {
		go func() {
			defer wg.Done()
			for url := range urls {
				res, err := urlHash(url)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s %s\n", url, res)
			}
		}()
	}
	for _, url := range args {
		urls <- url
	}
	close(urls)
	wg.Wait()
}
