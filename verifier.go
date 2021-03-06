package processor

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type SmartyVerifier struct {
	client HTTPClient
}

func NewSmartyVerifier(client HTTPClient) *SmartyVerifier { return &SmartyVerifier{client: client} }

func (sv *SmartyVerifier) Verify(input AddressInput) AddressOutput {
	request := sv.buildRequest(input)
	response, _ := sv.client.Do(request)
	candidates := sv.decodeResponse(response)

	return sv.prepareAddressOutput(candidates)
}

func (sv *SmartyVerifier) decodeResponse(response *http.Response) (output []Candidate) {
	if response != nil {
		defer response.Body.Close()
		body := response.Body
		json.NewDecoder(body).Decode(&output)
	}
	return output
}

func (sv *SmartyVerifier) prepareAddressOutput(candidates []Candidate) AddressOutput {
	if len(candidates) == 0 {
		return AddressOutput{Status: "Invalid API Response JSON"}
	}
	candidate := candidates[0]
	status := computeStatus(candidate)
	return AddressOutput{
		Status:        status,
		DeliveryLine1: candidate.DeliveryLine1,
		LastLine:      candidate.LastLine,
		City:          candidate.Components.City,
		State:         candidate.Components.State,
		ZIPCode:       candidate.Components.ZIPCode,
	}
}
func computeStatus(candidate Candidate) string {
	analysis := candidate.Analysis
	if !isDeliverable(analysis.Match) {
		return "Invalid"
	} else if analysis.Vacant == "Y" {
		return "Vacant"
	} else if analysis.Active != "Y" {
		return "Inactive"
	} else {
		return "Deliverable"
	}
}
func isDeliverable(value string) bool {
	return value == "Y" || value == "S" || value == "D"
}
func (verifier *SmartyVerifier) buildRequest(input AddressInput) *http.Request {
	var query url.Values = make(url.Values)
	query.Set("street", input.Street1)
	query.Set("city", input.City)
	query.Set("state", input.State)
	query.Set("zipcode", input.ZIPCode)

	request, _ := http.NewRequest("GET", "/street-address?"+query.Encode(), nil)
	return request

}
