package processor

import (
	"net/http"
	"testing"

	"github.com/smartystreets/gunit"
)

func TestWireupFixture(t *testing.T) {
	gunit.Run(new(WireupFixture), t)
}

type WireupFixture struct {
	*gunit.Fixture

	reader  *ReadWriteSpyBuffer
	writer  *ReadWriteSpyBuffer
	client  *FakeHTTPClient
	handler Handler
}

func (wf *WireupFixture) Setup() {
	wf.reader = NewReadWriteSpyBuffer("")
	wf.writer = NewReadWriteSpyBuffer("")
	wf.client = &FakeHTTPClient{}
	wf.handler = Configure(wf.reader, wf.writer, wf.client).Build()

}
func (wf *WireupFixture) LongTestPipeline() {
	wf.client.Configure(integrationJSONOutput, http.StatusOK, nil)
	wf.reader.WriteString("A,B,C,D")
	wf.reader.WriteString("A,B,C,D")
	wf.handler.Handle()
	expected := "Status,DeliveryLine1,LastLine,City,State,ZIPCode\n" +
		"Deliverable,AA,BB,CC,DD,EE\n" +
		"Deliverable,AA,BB,CC,DD,EE\n"
	wf.AssertEqual(expected, wf.writer.String())
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
