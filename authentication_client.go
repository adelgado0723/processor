package processor

import (
	"fmt"
	"net/http"
)

// AuthenticationClient Decorates HTTP Client
type AuthenticationClient struct {
	inner     HTTPClient
	scheme    string
	hostname  string
	authID    string
	authToken string
}

func NewAuthenticationClient(inner HTTPClient, scheme string, hostname string, authId string, authToken string) *AuthenticationClient {

	return &AuthenticationClient{
		inner:     inner,
		scheme:    scheme,
		hostname:  hostname,
		authID:    authId,
		authToken: authToken,
	}
}

func (ac *AuthenticationClient) Do(request *http.Request) (*http.Response, error) {
	request.URL.Scheme = ac.scheme
	request.URL.Host = ac.hostname
	query := request.URL.Query()
	query.Set("auth-id", ac.authID)
	query.Set("auth-token", ac.authToken)
	request.URL.RawQuery = query.Encode()
	fmt.Println(request.URL.String())

	request.Host = ac.hostname
	return ac.inner.Do(request)
}
