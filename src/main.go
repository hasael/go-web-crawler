package main

import (
	"fmt"
	"hasael/web-crawler/urlcrawler"
	"os"
)

func main() {
	inUrl := os.Args[1]
	fmt.Println("Crawling " + inUrl)
	urls, err := urlcrawler.GetUrls(inUrl)

	if err != nil {
		fmt.Print(err.Error())
	}

	for _, u := range urls {
		fmt.Println(u)
	}

}
