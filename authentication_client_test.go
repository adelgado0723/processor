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
	af.client = NewAuthenticationClient(af.inner, "http", "otherurl.com", "authid", "authtoken")
}

func (af *AuthenticationClientFixture) assertQueryStringValue(key, expectedValue string) {

	af.AssertEqual(af.inner.request.URL.Query().Get(key), expectedValue)
}
func (af *AuthenticationClientFixture) TestProvidedInformationAddedBeforeRequestIsSent() {
	// this request doesn't return an error
	// and just panics
	request := httptest.NewRequest("GET", "/path?existingKey=existingValue", nil)
	af.client.Do(request)

	af.assertReequestConnectionInformation()
	af.assertQueryStringIncludesAuthentication()
}

func (af *AuthenticationClientFixture) assertReequestConnectionInformation() {
	af.AssertEqual(af.inner.request.URL.Scheme, "http")
	af.AssertEqual(af.inner.request.Host, "otherurl.com")
	af.AssertEqual(af.inner.request.URL.Host, "otherurl.com")
}
func (af *AuthenticationClientFixture) assertQueryStringIncludesAuthentication() {
	af.assertQueryStringValue("auth-id", "authid")
	af.assertQueryStringValue("auth-token", "authtoken")
	af.assertQueryStringValue("existingKey", "existingValue")

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
