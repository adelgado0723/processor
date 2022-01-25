package processor

import (
	"bytes"
	"errors"
	"fmt"
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

func (vf *VerifierFixture) TestMalformedJSONHandled() {

	const malformedRawJSONOutput = `alert('Hello, world!' DROP TABLE *);`
	vf.client.Configure(malformedRawJSONOutput, http.StatusOK, nil)
	result := vf.verifier.Verify(AddressInput{})
	vf.AssertEqual(result.Status, "Invalid API Response JSON")

}

func (vf *VerifierFixture) TestHTTPErrorHandled() {
	vf.client.Configure("", 0, errors.New("GOPHERS!"))
	result := vf.verifier.Verify(AddressInput{})
	vf.Assert(result.Status != "")
}

func (vf *VerifierFixture) TestHTTPResponseBodyClosed() {
	vf.client.Configure(rawJSONOutput, http.StatusOK, nil)
	vf.verifier.Verify(AddressInput{})
	vf.AssertEqual(vf.client.responseBody.closed, 1)
}

func (vf *VerifierFixture) TestAddressStatus() {
	var (
		deliverableJSON      string = buildAnalysisJSON("Y", "N", "Y")
		missingSecondaryJSON string = buildAnalysisJSON("D", "N", "Y")
		droppedSecondaryJSON string = buildAnalysisJSON("S", "N", "Y")
		vacantJSON           string = buildAnalysisJSON("Y", "Y", "Y")
		inactiveJSON         string = buildAnalysisJSON("Y", "N", "?")
		invalidJSON          string = buildAnalysisJSON("N", "?", "?")
	)
	vf.assertStatus(deliverableJSON, "Deliverable")
	vf.assertStatus(missingSecondaryJSON, "Deliverable")
	vf.assertStatus(droppedSecondaryJSON, "Deliverable")
	vf.assertStatus(vacantJSON, "Vacant")
	vf.assertStatus(inactiveJSON, "Inactive")
	vf.assertStatus(invalidJSON, "Invalid")

}
func buildAnalysisJSON(match, vacant, active string) string {
	template := `[
	{
		"delivery_line_1": "1 Santa Claus Ln",
		"last_line": "North Pole AK 99705-9901",
		"analysis": {
			"dpv_match_code": "%s",
			"dpv_vacant": "%s",
			"active": "%s"
		}
	}
]`
	return fmt.Sprintf(template, match, vacant, active)
}

func (vf *VerifierFixture) assertStatus(jsonResponse, expectedStatus string) {
	vf.client.Configure(jsonResponse, http.StatusOK, nil)
	output := vf.verifier.Verify(AddressInput{})
	vf.AssertEqual(output.Status, expectedStatus)
}

////////////////////////////////////////////////////////////////

type FakeHTTPClient struct {
	request      *http.Request
	response     *http.Response
	responseBody *SpyBuffer
	err          error
}

func (fhc *FakeHTTPClient) Configure(responsetext string, statusCode int, err error) {
	if err == nil {
		fhc.responseBody = NewSpyBuffer(responsetext)

		fhc.response = &http.Response{
			Body:       fhc.responseBody,
			StatusCode: statusCode,
		}
	}
	fhc.err = err
}

func (fhc *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {
	fhc.request = request
	return fhc.response, fhc.err
}

////////////////////////////////////////////////////////////////
// Create a Spy Buffer that counts how many times close() was called
type SpyBuffer struct {
	// this syntax allows SpyBuffer to have the same Read/Close/etc
	// functionality that bytes.Buffer has without having to implement those methods
	*bytes.Buffer
	closed int
}

func (sb *SpyBuffer) Close() error {
	sb.closed++
	sb.Buffer.Reset()
	return nil
}

func NewSpyBuffer(value string) *SpyBuffer { return &SpyBuffer{Buffer: bytes.NewBufferString(value)} }
