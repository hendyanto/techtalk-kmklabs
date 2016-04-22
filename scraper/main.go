package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

func crawl(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("Error: Failed to crawl \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			ok, url := getHref(t)
			if !ok {
				continue
			}

			hasProto := strings.Index(url, "http") == 0
			hasJpg := strings.Contains(url, ".jpg") == true

			if hasProto && hasJpg {
				ch <- url
			}
		}
	}
}

func downloadImage(url string, fileNumber int, chFileNumber chan int, chFinished chan bool) {
	fmt.Println(" - " + url)
	filename := fmt.Sprintf("download_folder/%d-downloaded.jpg", fileNumber)
	resp, _ := http.Get(url)
	data, err2 := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err2 == nil {
		ioutil.WriteFile(filename, data, 0777)
		chFileNumber <- fileNumber
	}

	defer func() {
		chFinished <- true
	}()
}

func main() {

	foundUrls := make(map[string]bool)
	finishedDownloaded := make(map[int]bool)
	seedUrls := os.Args[1:]

	chUrls := make(chan string)
	chFinished := make(chan bool)
	chFileNumber := make(chan int)

	for _, url := range seedUrls {
		go crawl(url, chUrls, chFinished)
	}

	for c := 0; c < len(seedUrls); {
		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFinished:
			c++
		}
	}

	fmt.Println("\nFound", len(foundUrls), "unique urls:\n")
	os.Mkdir("download_folder", 0777)
	fmt.Println("\nDownloading")

	c := 0
	for url, _ := range foundUrls {
		go downloadImage(url, c, chFileNumber, chFinished)
		c++
	}

	for c := 0; c < len(foundUrls); {
		select {
		case fileNumber := <-chFileNumber:
			finishedDownloaded[fileNumber] = true
			fmt.Println(fileNumber)
		case <-chFinished:
			c++
		}
	}

	close(chUrls)
}
