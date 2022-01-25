package processor

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/gunit"
)

func TestAuthenticationClient(t *testing.T) {
	gunit.Run(new(AuthenticationClientFixture), t)
}

type AuthenticationClientFixture struct {
	*gunit.Fixture
	inner  *FakeHTTPClient
	client *AuthenticationClient
}

func (af *AuthenticationClientFixture) Setup() {
	af.inner = &FakeHTTPClient{}
	af.client = NewAuthenticationClient(af.inner, "http", "otherurl.com", "authid", "auth-token")
}

func (af *AuthenticationClientFixture) TestProvidedInformationAddedBeforeRequestIsSent() {
	// this request doesn't return an error
	// and just panics
	request := httptest.NewRequest("GET", "/path", nil)
	af.client.Do(request)
	af.AssertEqual(af.inner.request.Host, "otherurl.com")
	af.AssertEqual(af.inner.request.URL.Scheme, "http")
	af.AssertEqual(af.inner.request.URL.Host, "otherurl.com")
	af.AssertEqual(af.inner.request.URL.Query().Get("auth-id"), "authid")
	af.AssertEqual(af.inner.request.URL.Query().Get("auth-token"), "authtoken")
}

func (af *AuthenticationClientFixture) TestResponseFromInnerClientReturned() {
	af.inner.response = &http.Response{
		StatusCode: http.StatusTeapot,
	}
	af.inner.err = errors.New("HTTP ERROR")
	request := httptest.NewRequest("GET", "/path", nil)
	response, err := af.client.Do(request)
	af.AssertEqual(response.StatusCode, http.StatusTeapot)
	af.AssertEqual(err.Error(), "HTTP ERROR")
}
