package processor

import (
	"log"
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
	log.SetFlags(log.Llongfile | log.Lmicroseconds)
	pf.reader = NewReadWriteSpyBuffer("")
	pf.writer = NewReadWriteSpyBuffer("")
	pf.client = &IntegrationHTTPClient{}
	pf.pipeline = Configure(pf.reader, pf.writer, pf.client, 2) // .Handle()?

}
func (pf *PipelineFixture) LongTestPipeline() {
	pf.reader.WriteString("A,B,C,D")
	pf.reader.WriteString("A,B,C,D")
	err := pf.pipeline.Process()
	expected := "Status,DeliveryLine1,LastLine,City,State,ZIPCode\n" +
		"Deliverable,AA,BB,CC,DD,EE\n" +
		"Deliverable,AA,BB,CC,DD,EE\n"
	pf.AssertEqual(expected, pf.writer.String())

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
