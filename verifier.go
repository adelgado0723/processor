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
	return AddressOutput{
		DeliveryLine1: candidate.DeliveryLine1,
		LastLine:      candidate.LastLine,
		City:          candidate.Components.City,
		State:         candidate.Components.State,
		ZIPCode:       candidate.Components.ZIPCode,
	}
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

type Candidate struct {
	DeliveryLine1 string `json:"delivery_line_1"`
	LastLine      string `json:"last_line"`
	Components    struct {
		City    string `json:"city_name"`
		State   string `json:"state_abbreviation"`
		ZIPCode string `json:"zipcode"`
	} `json:"components"`
}
