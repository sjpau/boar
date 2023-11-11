package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

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
