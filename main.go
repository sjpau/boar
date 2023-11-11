package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var (
	TOR_PROXY     string = "socks5://127.0.0.1:9050"
	URL           string = "http://libraryfyuybp7oyidyya3ah5xvwgyx6weauoini7zyz555litmmumad.onion"
	STRING_SEARCH string = "s?query="
	STRING_TYPE   string = "src="
)

/* NOTE: Unimplemented features */
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
	var err error

	log.Println(client, resp)

	InitCmdOptions()
	if FlagURL == "" {
		FlagURL = URL
	}
	if FlagProxy == "" {
		FlagProxy = TOR_PROXY
	}

	/* TODO: Add multiple sources
	if FlagSrc == "" {
		FlagSrc = "bk"
	}
	*/
	FlagSrc = "bk"
	client = NewTorProxyClient(FlagProxy)
	log.Printf("Client is set.\n")
	log.Printf("Setting request parameters...\n")
	searchUrl := JoinSearchUrl()
	req, err := CreateRequest(searchUrl, "GET")
	Check(err)
	log.Printf("Sending request...\n")
	resp, err = client.Do(req)
	Check(err)
	defer resp.Body.Close()
	log.Printf("Body received.\n")
	log.Printf("Reading body...\n")
	body, err := io.ReadAll(resp.Body)
	Check(err)
	log.Printf("Parsing body...\n")
	doc, err := html.Parse(strings.NewReader(string(body)))
	Check(err)
	elems := MapElementsFormHTML(doc)
	tits := PrintTitlesf(elems)
	input := GetDownloadInput()
	downloadIndicies := DownloadInputStringToIntSlice(input)
	for i := 0; i < len(downloadIndicies); i++ {
		log.Printf("Downloading %s\n", elems[tits[downloadIndicies[i]]])
		DownloadWithClient(client, elems[tits[downloadIndicies[i]]], tits[downloadIndicies[i]])
		log.Println("Finished.")
	}
}

func DownloadInputStringToIntSlice(input string) []int {
	parts := strings.Fields(input)
	var numbers []int

	for _, part := range parts {
		number, err := strconv.Atoi(part)
		if err != nil {
			fmt.Printf("Error converting %s to an integer: %v\n", part, err)
			continue
		}
		numbers = append(numbers, number)
	}
	return numbers
}

func GetDownloadInput() string {
	fmt.Print("Indices to download (ex.: 1 4 12 2): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input.", err)
		return ""
	}

	input = strings.TrimSuffix(input, "\n")
	return input
}

func PrintTitlesf(elems map[string]string) []string {
	var titles []string
	for k := range elems {
		titles = append(titles, k)
	}
	for i := 0; i < len(titles); i++ {
		fmt.Printf("[%d]: %s\n", i, titles[i])
	}
	return titles
}

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

func MapElementsFormHTML(doc *html.Node) map[string]string {
	var titles []string

	hrefs := FindElementsByClass(doc, SOURCES[FlagSrc])
	l := len(hrefs)
	for i := 0; i < l; i++ {
		titles = append(titles, FindElementFromAttr(hrefs[i].FirstChild.FirstChild, "img", "title"))
	}

	elems := make(map[string]string, l)
	for i := 0; i < l; i++ {
		elems[titles[i]] = fmt.Sprintf("%s%s", FlagURL, FindElementFromAttr(hrefs[i], "a", "href"))
	}
	return elems
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

func JoinSearchUrl() string {
	var searchUrl string
	if FlagQuery == "" {
		log.Fatal("Query is empty!\n")
	}
	searchUrl = fmt.Sprintf("%s%s%s%s%s%s%s", FlagURL, "/", STRING_SEARCH, FlagQuery, "&", STRING_TYPE, FlagSrc)
	log.Printf("Joined URL: %s\n", searchUrl)
	return searchUrl
}

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
