package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func DownloadWithClient(client *http.Client, URL, fileName string) error {
	req, err := CreateRequest(URL, "GET")
	Check(err)
	response, err := client.Do(req)
	Check(err)
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("Received %d response code when downloading %s.", response.StatusCode, URL)
		return errors.New("Download failed!")
	}
	file, err := os.Create(filepath.Join(FlagDest, fileName))
	Check(err)
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	Check(err)

	return nil
}

func CreateRequest(URL, r string) (*http.Request, error) {
	req, err := http.NewRequest(r, URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	return req, nil
}

func NewTorProxyClient(proxyPath string) *http.Client {
	log.Printf("Proxy path: %s\n", proxyPath)
	torProxyUrl, err := url.Parse(proxyPath)
	if err != nil {
		log.Fatal("Error parsing Tor proxy URL:", proxyPath, ".", err)
	}

	log.Printf("Creating tor transport...\n")
	torTransport := &http.Transport{Proxy: http.ProxyURL(torProxyUrl)}
	log.Printf("Setting the client...\n")
	client := &http.Client{Transport: torTransport, Timeout: time.Second * 60 * 5}
	return client
}
