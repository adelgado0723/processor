package processor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/smartystreets/gunit"
)

func TestPipelineFixture(t *testing.T) {
	gunit.Run(new(PipelineFixture), t)
}

type PipelineFixture struct {
	*gunit.Fixture

	reader   *ReadWriteSpyBuffer
	writer   *ReadWriteSpyBuffer
	client   *IntegrationHTTPClient
	pipeline *Pipeline
}

func (pf *PipelineFixture) Setup() {
	pf.reader = NewReadWriteSpyBuffer("")
	pf.writer = NewReadWriteSpyBuffer("")
	pf.client = &IntegrationHTTPClient{}
	pf.pipeline = Configure(ioutil.NopCloser(pf.reader), pf.writer, pf.client, 2)
}
func (pf *PipelineFixture) LongTestPipeline() {
	fmt.Fprintln(pf.reader, "Street1,City,State,ZIPCode")
	fmt.Fprintln(pf.reader, "A,B,C,D")
	fmt.Fprintln(pf.reader, "A,B,C,D")

	err := pf.pipeline.Process()

	expected := "Status,DeliveryLine1,LastLine,City,State,ZIPCode\n" +
		"Deliverable,AA,BB,CC,DD,EE\n" +
		"Deliverable,AA,BB,CC,DD,EE\n"
	pf.AssertEqual(expected, pf.writer.String())
	fmt.Println("expected: ", expected, "\nactual: ", pf.writer.String())

	pf.Assert(err == nil)
}

type IntegrationHTTPClient struct{}

func (c *IntegrationHTTPClient) Do(request *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:       NewReadWriteSpyBuffer(integrationJSONOutput),
		StatusCode: http.StatusOK,
	}, nil
}

const integrationJSONOutput = `
[
	{
		"delivery_line_1": "AA",
		"last_line": "BB",
		"components": {
			"city_name": "CC",
			"state_abbreviation": "DD",
			"zipcode": "EE"
		},
		"analysis": {
			"dpv_match_code": "Y",
			"dpv_vacant": "N",
			"active": "Y"
		}
	}
]`
