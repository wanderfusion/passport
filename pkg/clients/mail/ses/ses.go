package ses

import (
	"net/http"

	"github.com/akxcix/passport/pkg/commons/httpclient"
)

type SesClient struct {
	httpClient *http.Client
}

func New() *SesClient {
	httpClient := httpclient.NewPooledHttpClient(10, 10, 10, 1000)

	sesClient := SesClient{
		httpClient: httpClient,
	}

	return &sesClient
}
