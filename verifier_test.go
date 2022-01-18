package processor

import (
	"bytes"
	"io/ioutil"
	"net/http"
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
		City:    "City",
		State:   "State",
		ZIPCode: "ZIPCode",
	}
	vf.client.Configure("[{}]", http.StatusOK, nil)
	vf.verifier.Verify(input)

	vf.AssertEqual("GET", vf.client.request.Method)
	vf.AssertQueryStringValue("street", input.Street1)
	vf.AssertQueryStringValue("city", input.City)
	vf.AssertQueryStringValue("state", input.State)
	vf.AssertQueryStringValue("zipcode", input.ZIPCode)
	vf.AssertEqual("/street-address", vf.client.request.URL.Path)
}

func (vf *VerifierFixture) AssertQueryStringValue(key, expected string) {
	query := vf.client.request.URL.Query()
	vf.AssertEqual(expected, query.Get(key))
}

func (vf *VerifierFixture) TestResponseParsed() {
	vf.client.Configure(rawJSONOutput, http.StatusOK, nil)
	result := vf.verifier.Verify(AddressInput{})
	vf.AssertEqual(result.DeliveryLine1, "1 Santa Claus Ln")
	vf.AssertEqual(result.LastLine, "North Pole AK 99705-9901")
	vf.AssertEqual(result.City, "North Pole")
	vf.AssertEqual(result.State, "AK")
	vf.AssertEqual(result.ZIPCode, "99705")
}

const rawJSONOutput = `
[
	{
		"delivery_line_1": "1 Santa Claus Ln",
		"last_line": "North Pole AK 99705-9901",
		"components": {
			"city_name": "North Pole",
			"state_abbreviation": "AK",
			"zipcode": "99705"
		}
	}
]`

func (sv *VerifierFixture) TestMalformedJSONHandled()

const malformedRawJSONOutput = `alert('Hello, world!' DROP TABLE *);`

////////////////////////////////////////////////////////////////

type FakeHTTPClient struct {
	request  *http.Request
	response *http.Response
	err      error
}

func (fhc *FakeHTTPClient) Configure(responsetext string, statusCode int, err error) {
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
