package dispatcher

import (
	"net/http"
	url2 "net/url"
	"strings"
	"testing"
)

func TestPrepareRequest(t *testing.T) {
	request, _ := http.NewRequest("GET", "https://test.com/search1?q=12312&122", nil)
	url, _ := url2.Parse("https://google.com/search2?a=123&z=00aaaa")

	prepareRequest(request, url)
	t.Logf("%+v\n", request)
	n, _ := strings.CutPrefix(request.URL.Path, "/search1")
	t.Log(n)
}

func TestPrepareRequestOptimized(t *testing.T) {
	request, _ := http.NewRequest("GET", "https://test.com/search?q=12312&122", nil)
	url, _ := url2.Parse("https://google.com/?a=123&z=00aaaa")

	prepareRequest(request, url)
	t.Logf("%+v\n", request)
}
