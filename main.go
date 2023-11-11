package main

import (
	"io"
	"log"
	"net/http"
	"strings"

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
