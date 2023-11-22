package urlcrawler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func (u UrlCrawlerT) CrawlUrls(inputUrl string) ([]string, error) {

	crawledUrls := make(map[string]bool)

	u.innerCrawlUrls(inputUrl, crawledUrls)
	result := make([]string, len(crawledUrls))
	for k := range crawledUrls {
		result = append(result, k)
	}
	return result, nil

}

func (u UrlCrawlerT) innerCrawlUrls(inputUrl string, crawled map[string]bool) {

	if inputUrl != "" {
		foundUrls, err := u.GetUrls(inputUrl)

		if err != nil {
			fmt.Println(err)
			return
		}

		for _, foundUrl := range foundUrls {
			_, ok := crawled[foundUrl]

			if !ok {
				crawled[foundUrl] = true
				u.innerCrawlUrls(foundUrl, crawled)
			}
		}
	}

}

type UrlClient interface {
	Get(url string) (resp *http.Response, err error)
}

type UrlHttpClient struct {
}

func (u UrlHttpClient) Get(url string) (resp *http.Response, err error) {
	return http.Get(url)
}

type UrlCrawlerT struct {
	Client UrlClient
}

func (h UrlCrawlerT) GetUrls(inputUrl string) ([]string, error) {
	url, err := url.Parse(inputUrl)

	if err != nil {
		return nil, err
	}

	hostname := strings.TrimPrefix(url.Hostname(), "www.")

	resp, err := h.Client.Get(inputUrl)

	if err != nil {
		return nil, fmt.Errorf("GET error: %v %v", inputUrl, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status error %v %v", inputUrl, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("read body: %v", err)
	}

	return parseForUrls(data, hostname)
}

func parseForUrls(data []byte, domain string) ([]string, error) {

	reg, err := regexp.Compile(fmt.Sprintf("(https?:\\/\\/%v(\\/[A-Za-z0-9\\-\\._~:\\/\\?\\[\\]@!$'\\(\\*\\+,;]*)?)", domain))
	if err != nil {
		return nil, fmt.Errorf("error while compiling regexp: %v", err)
	}

	matches := reg.FindAll(data, -1)

	var result = make([]string, len(matches))

	for _, byteArr := range matches {
		result = append(result, string(byteArr))
	}

	return result, nil

}
