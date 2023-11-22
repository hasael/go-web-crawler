package main

import (
	"fmt"
	"hasael/web-crawler/urlcrawler"
	"os"
)

func main() {
	inUrl := os.Args[1]
	fmt.Println("Crawling " + inUrl)
	u := new(urlcrawler.UrlCrawlerT)
	u.Client = new(urlcrawler.UrlHttpClient)
	urls, err := u.CrawlUrls(inUrl)

	if err != nil {
		fmt.Print(err.Error())
	}

	for _, u := range urls {
		fmt.Println(u)
	}
}
