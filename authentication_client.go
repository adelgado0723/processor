package processor

import "net/http"

type AuthenticationClient struct {
	inner     HTTPClient
	scheme    string
	hostname  string
	authId    string
	authToken string
}

func NewAuthenticationClient(inner HTTPClient, scheme string, hostname string, authId string, authToken string) *AuthenticationClient {

	return &AuthenticationClient{
		inner:     inner,
		scheme:    scheme,
		hostname:  hostname,
		authId:    authId,
		authToken: authToken,
	}
}

func (ac *AuthenticationClient) Do(request *http.Request) (*http.Response, error) {
	request.URL.Scheme = ac.scheme
	request.Host = ac.hostname
	request.URL.Host = ac.hostname
	return ac.inner.Do(request)
}
