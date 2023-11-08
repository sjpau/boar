package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var (
	TOR_PROXY     string = "socks5://127.0.0.1:9050"
	URL           string = "http://libraryfyuybp7oyidyya3ah5xvwgyx6weauoini7zyz555litmmumad.onion/"
	STRING_SEARCH string = "s?query="
	STRING_TYPE   string = "src="
)

var SOURCES = map[string]string{
	"bk":   "book",
	"art":  "article",
	"img":  "images",
	"vid":  "video",
	"adbk": "audiobook",
}

func main() {

	var client *http.Client
	var resp *http.Response
	var req *http.Request
	var err error

	log.Println(client, resp)

	InitCmdOptions()
	if FlagProxy == "" {
		client = NewTorProxyClient(TOR_PROXY)
	} else {
		client = NewTorProxyClient(FlagProxy)
	}
	log.Printf("Client is set.\n")
	if FlagSrc == "" {
		FlagSrc = "bk"
	}
	log.Printf("Setting request parameters...\n")
	searchUrl := JoinSearchUrl()
	req, err = http.NewRequest("GET", searchUrl, nil)
	Check(err)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	log.Printf("Sending request...\n")
	//resp, err = client.Do(req)
	//Check(err)
	//defer resp.Body.Close()
	log.Printf("Body received.\n")
	//body, err := io.ReadAll(resp.Body)
	body, err := os.ReadFile("index.html")
	Check(err)
	doc, err := html.Parse(strings.NewReader(string(body)))
	Check(err)
	log.Printf("Traversing HTML for class \"%s\"", SOURCES[FlagSrc])
	elems := FindElementsByClass(doc, SOURCES[FlagSrc])
	for i := 0; i < len(elems); i++ {
		log.Println((FindElementFromAttr(elems[i].FirstChild.FirstChild, "img", "title")))
	}
}

func FindElementsByClass(n *html.Node, targetClass string) []*html.Node {
	elements := make([]*html.Node, 0)

	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "class" {
				classes := strings.Fields(attr.Val)
				for _, class := range classes {
					if class == targetClass {
						elements = append(elements, n)
						break
					}
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		elements = append(elements, FindElementsByClass(c, targetClass)...)
	}

	return elements
}

func FindElementFromAttr(node *html.Node, data string, key string) string {
	var value string

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == data {
			for _, attr := range n.Attr {
				if attr.Key == key {
					value = attr.Val
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	if node != nil {
		extract(node)
	}

	return value
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

func JoinSearchUrl() string {
	var searchUrl string
	if FlagQuery == "" {
		log.Fatal("Query is empty!\n")
	}
	if FlagURL == "" {
		searchUrl = URL
	} else {
		searchUrl = FlagURL
	}
	searchUrl = fmt.Sprintf("%s%s%s%s%s%s", searchUrl, STRING_SEARCH, FlagQuery, "&", STRING_TYPE, FlagSrc)
	log.Printf("Joined URL: %s\n", searchUrl)
	return searchUrl
}

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
