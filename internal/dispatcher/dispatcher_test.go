package dispatcher

import (
	"net/http"
	url2 "net/url"
	"testing"
)

func TestPrepareRequestPath(t *testing.T) {
	request, _ := http.NewRequest("GET", "https://test.com/search1?q=12312&122", nil)
	url, _ := url2.Parse("https://google.com/search2?a=123&z=00aaaa")

	const expectedUrl = "https://google.com/search2?q=12312&122&a=123&z=00aaaa"
	prepareRequest(request, url, &cutPathPrefix{cut: true, value: "/search1"})

	if request.URL.String() != expectedUrl {
		t.Logf("Expected %s, got %s", expectedUrl, request.URL.String())
		t.FailNow()
	}

	t.Log("Request matched expected result!")
}

func TestPrepareRequestPathWithMalformedUrl(t *testing.T) {
	request, _ := http.NewRequest("GET", "https://test.com//search1?q=12312&122", nil)
	url, _ := url2.Parse("https://google.com/search2?a=123&z=00aaaa")

	const expectedUrl = "https://google.com/search2?q=12312&122&a=123&z=00aaaa"
	prepareRequest(request, url, &cutPathPrefix{cut: true, value: "/search1"})

	if request.URL.String() != expectedUrl {
		t.Logf("Expected %s, got %s", expectedUrl, request.URL.String())
		t.FailNow()
	}

	t.Log("Request matched expected result!")
}

func TestPrepareRequestSubdomain(t *testing.T) {
	request, _ := http.NewRequest("GET", "https://t.test.com/search?q=12312&122", nil)
	url, _ := url2.Parse("https://google.com/")

	const expectedUrl = "https://google.com/search?q=12312&122"
	prepareRequest(request, url, &cutPathPrefix{cut: false})

	if request.URL.String() != expectedUrl {
		t.Logf("Expected %s, got %s", expectedUrl, request.URL.String())
		t.FailNow()
	}

	t.Log("Request matched expected result!")
}

func TestPrepareRequestSubdomainWithBadPath(t *testing.T) {
	request, _ := http.NewRequest("GET", "https://t.test.com//search?q=12312&122", nil)
	url, _ := url2.Parse("https://google.com//")

	const expectedUrl = "https://google.com/search?q=12312&122"
	prepareRequest(request, url, &cutPathPrefix{cut: false})

	if request.URL.String() != expectedUrl {
		t.Logf("Expected %s, got %s", expectedUrl, request.URL.String())
		t.FailNow()
	}

	t.Log("Request matched expected result!")
}
