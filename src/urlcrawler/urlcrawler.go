package urlcrawler

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func GetUrls(url string) ([]string, error) {

	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status error %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("read body: %v", err)
	}

	return parseForUrls(data, "monzo.com")
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
