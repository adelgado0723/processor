package processor

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/smartystreets/gunit"
)

func TestVerifierFixture(t *testing.T) {
	gunit.Run(new(VerifierFixture), t)
}

type VerifierFixture struct {
	*gunit.Fixture
	client   *FakeHTTPClient
	verifier *SmartyVerifier
}

func (vf *VerifierFixture) Setup() {
	vf.client = &FakeHTTPClient{}
	vf.verifier = NewSmartyVerifier(vf.client)

}
func NewSmartyVerifier(client HTTPClient) *SmartyVerifier { return &SmartyVerifier{client: client} }

func (vf *VerifierFixture) rawQuery() string {
	return vf.client.request.URL.RawQuery
}

func (vf *VerifierFixture) TestRequestComposedProperly() {
	input := AddressInput{
		Street1: "Street1",
		City:    "City & City",
		State:   "State",
		ZIPCode: "ZIPCode",
	}
	vf.verifier.Verify(input)

	vf.AssertEqual("GET", vf.client.request.Method)
	vf.AssertQueryStringValue("street", input.Street1)
	vf.AssertQueryStringValue("city", input.City)
	vf.AssertQueryStringValue("state", input.State)
	vf.AssertQueryStringValue("zipcode", input.ZIPCode)
	vf.AssertQueryStringValue("street", "Street1")
	vf.AssertEqual("/street-address", vf.client.request.URL.Path)
	vf.Assert(strings.Contains(vf.client.request.URL.RawQuery, "%26"))

}

func (vf *VerifierFixture) AssertQueryStringValue(key, expected string) {
	query := vf.client.request.URL.Query()
	vf.AssertEqual(expected, query.Get(key))
}

func (vf *VerifierFixture) TestResponseParsed() {
	vf.client.response = &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString(`[{}]`)),
		StatusCode: http.StatusOK,
	}
	result := vf.verifier.Verify(AddressInput{})
	vf.AssertEqual(result.DeliveryLine1, "Hello World")
}

////////////////////////////////////////////////////////////////

type FakeHTTPClient struct {
	request  *http.Request
	response *http.Response
	err      error
}

func (fhc *FakeHTTPClient) NewFakeHTTPClient(responsetext string, statusCode int, err error) {
	fhc.response = &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString(responsetext)),
		StatusCode: statusCode,
	}
	fhc.err = err
}

func (fhc *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {
	fhc.request = request
	return fhc.response, fhc.err
}
