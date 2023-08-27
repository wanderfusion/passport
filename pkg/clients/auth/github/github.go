package github

import (
	"net/http"

	"github.com/akxcix/passport/pkg/commons/httpclient"
)

type GithubOauthClient struct {
	host         string
	clientID     string
	clientSecret string
	client       *http.Client
}

func New(host, clientId, clientSecret string) *GithubOauthClient {
	pooledHttpClient := httpclient.NewPooledHttpClient(10, 10, 10, 1000)

	githubOauthClient := &GithubOauthClient{
		host:         host,
		clientID:     clientId,
		clientSecret: clientSecret,
		client:       pooledHttpClient,
	}

	return githubOauthClient
}
