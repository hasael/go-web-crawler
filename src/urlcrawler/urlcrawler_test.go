package urlcrawler

import (
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
)

type MockCall struct {
	InputUrl   string
	OutputBody string
}

type TestUrlClient struct {
	Mocks []MockCall
}

func (u TestUrlClient) Get(url string) (resp *http.Response, err error) {

	for _, mock := range u.Mocks {

		if mock.InputUrl == url {
			return &http.Response{
				StatusCode:       200,
				Status:           "",
				Proto:            "",
				ProtoMajor:       1,
				ProtoMinor:       0,
				Header:           make(http.Header),
				Body:             io.NopCloser(strings.NewReader(mock.OutputBody)),
				ContentLength:    10,
				Close:            true,
				TransferEncoding: nil,
				Uncompressed:     true,
				Trailer:          nil,
				Request:          nil,
				TLS:              nil,
			}, nil
		}
	}

	return &http.Response{
		StatusCode:       200,
		Status:           "",
		Proto:            "",
		ProtoMajor:       1,
		ProtoMinor:       0,
		Header:           make(http.Header),
		Body:             io.NopCloser(strings.NewReader("")),
		ContentLength:    10,
		Close:            true,
		TransferEncoding: nil,
		Uncompressed:     true,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}, nil

}

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func Test_CrawlUrls(t *testing.T) {

	mocks := make([]MockCall, 5)
	mocks = append(mocks, MockCall{"http://localhost:8089", "<response>http://localhost:8089/test1</response>"})
	mocks = append(mocks, MockCall{"http://localhost:8089/test1", "<response>http://localhost:8089/test2</response>" +
		"<u>http://localhost:8089/test3</u>"})
	mocks = append(mocks, MockCall{"http://localhost:8089/test2", "<response>http://localhost:8089/test4</response>"})
	mocks = append(mocks, MockCall{"http://localhost:8089/test3", "<response>OK</response>"})
	mocks = append(mocks, MockCall{"http://localhost:8089/test4", "<response>OK</response>"})
	testClient := TestUrlClient{mocks}

	u := new(UrlCrawlerT)
	u.Client = testClient

	result, err := u.CrawlUrls("http://localhost:8089")

	if err != nil {
		t.Fatal("Error calling url")
	}

	if len(result) <= 0 {
		t.Fatal("No result received")
	}
	if len(result) != 4 {
		t.Fatal("Wrong results received")
	}

	if !slices.Contains(result, "http://localhost:8089/test1") {
		t.Fatal("'http://localhost:8089/test1' not found")
	}
	if !slices.Contains(result, "http://localhost:8089/test2") {
		t.Fatal("'http://localhost:8089/test2' not found")
	}
	if !slices.Contains(result, "http://localhost:8089/test3") {
		t.Fatal("'http://localhost:8089/test3' not found")
	}
	if !slices.Contains(result, "http://localhost:8089/test4") {
		t.Fatal("'http://localhost:8089/test4' not found")
	}

	t.Log(result)
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
func TestHelloEmpty(t *testing.T) {

}
