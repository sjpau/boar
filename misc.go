package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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
