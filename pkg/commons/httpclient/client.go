package httpclient

import (
	"net/http"
	"time"
)

func NewPooledHttpClient(maxIdleConns, maxConnsPerHost, maxIdleConnsPerHost, timeoutMillis int) *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = maxIdleConns
	t.MaxConnsPerHost = maxConnsPerHost
	t.MaxIdleConnsPerHost = maxIdleConnsPerHost

	httpClient := &http.Client{
		Timeout:   time.Duration(timeoutMillis) * time.Millisecond,
		Transport: t,
	}

	return httpClient
}
