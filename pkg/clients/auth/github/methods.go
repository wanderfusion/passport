package github

import (
	"github.com/wanderfusion/passport/pkg/commons/httpclient"
)

const (
	endpointOauth string = "/login/oauth/access_token"
)

func (c *GithubOauthClient) GetToken(code string) (string, error) {
	payload := OAuthRequest{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		Code:         code,
	}
	bytes, err := httpclient.MakePOSTRequest(c.client, c.host, endpointOauth, payload, nil, nil)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
