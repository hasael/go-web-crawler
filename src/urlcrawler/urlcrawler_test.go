package urlcrawler

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/wiremock/go-wiremock"
)

type TestUrlClient struct {
}

func (u TestUrlClient) Get(url string) (resp *http.Response, err error) {
	//if url == "http://localhost:8089" {
	return &http.Response{
		StatusCode:       200,
		Status:           "",
		Proto:            "",
		ProtoMajor:       1,
		ProtoMinor:       0,
		Header:           make(http.Header),
		Body:             io.NopCloser(strings.NewReader("<response>http://localhost:8089/test1</response>")),
		ContentLength:    10,
		Close:            true,
		TransferEncoding: nil,
		Uncompressed:     true,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}, nil
	//}
}

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func Test_CrawlUrls(t *testing.T) {
	wiremockClient := wiremock.NewClient("http://localhost:8089")

	defer wiremockClient.Reset()

	wiremockClient.StubFor(wiremock.Get(wiremock.URLPathEqualTo("/")).
		WillReturnResponse(
			wiremock.NewResponse().
				WithBody("<response>http://localhost:8089/test1</response>").
				WithHeader("Content-Type", "text/xml").
				WithStatus(http.StatusOK),
		).AtPriority(1))

	wiremockClient.StubFor(wiremock.Get(wiremock.URLEqualTo("/test1")).
		WillReturnResponse(
			wiremock.NewResponse().
				WithBody("<response>http://localhost:8089/test2</response>"+
					"<u>http://localhost:8089/test3</u>").
				WithHeader("Content-Type", "text/xml").
				WithStatus(http.StatusOK),
		))

	wiremockClient.StubFor(wiremock.Get(wiremock.URLEqualTo("/test2")).
		WillReturnResponse(
			wiremock.NewResponse().
				WithBody("<response>http://localhost:8089/test4</response>").
				WithHeader("Content-Type", "text/html").
				WithStatus(http.StatusOK),
		))

	wiremockClient.StubFor(wiremock.Get(wiremock.URLEqualTo("/test3")).
		WillReturnResponse(
			wiremock.NewResponse().
				WithBody("<response>OK</response>").
				WithHeader("Content-Type", "text/xml").
				WithStatus(http.StatusOK),
		))

	wiremockClient.StubFor(wiremock.Get(wiremock.URLEqualTo("/test4")).
		WillReturnResponse(
			wiremock.NewResponse().
				WithBody("<response>OK</response>").
				WithHeader("Content-Type", "text/xml").
				WithStatus(http.StatusOK),
		))

	u := new(UrlCrawlerT)
	u.Client = new(TestUrlClient)
	result, err := u.CrawlUrls("http://localhost:8089")

	if err != nil {
		t.Fatal("Error calling url")
	}

	if len(result) <= 0 {
		t.Fatal("No result received")
	}
	t.Log(result)
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
func TestHelloEmpty(t *testing.T) {

}
