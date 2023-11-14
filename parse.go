package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

func MapElementsFromHTML(doc *html.Node, tag, attr string) map[string]string {
	var titles []string

	hrefs := FindElementsByClass(doc, SOURCES[FlagSrc])
	l := len(hrefs)
	for i := 0; i < l; i++ {
		titles = append(titles, FindElementFromAttr(hrefs[i].FirstChild.FirstChild, tag, attr))
	}

	elems := make(map[string]string, l)
	for i := 0; i < l; i++ {
		s := fmt.Sprintf("%s%s", FlagURL, FindElementFromAttr(hrefs[i], "a", "href"))
		log.Printf("Checking: %s\n", s)
		if strings.HasSuffix(s, ".pdf") || strings.HasSuffix(s, ".epub") {
			log.Printf("Found a pdf link while parsing!\n")
			elems[titles[i]] = s
		} else if strings.HasSuffix(s, ".html") {
			log.Printf("Found html link, don't know what to do with it, sorry! (!UNIMPLEMENTED!)\n")
		} else if strings.HasSuffix(s, ".jpg") || strings.HasSuffix(s, ".png") || strings.HasSuffix(s, ".jpeg") {
			log.Printf("Found image link!\n")
			s = RemoveFromStr(s, "view?img=")
			elems[titles[i]] = s
		} else if strings.HasSuffix(s, ".mp4") || strings.HasSuffix(s, ".webm") || strings.HasSuffix(s, ".mkv") {
			log.Printf("Found video link!\n")
			s = RemoveFromStr(s, "watch?vid=")
			elems[titles[i]] = s
		} else if strings.HasSuffix(s, ".mp3") {
			log.Printf("Found audio link!\n")
			s = RemoveFromStr(s, "admeta?fadbk=")
			elems[titles[i]] = s
		}
	}
	return elems
}

func RemoveFromStr(s, strip string) string {
	var news string
	if strings.Contains(s, strip) {
		news = strings.Replace(s, strip, "", -1)
	}
	return news
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
