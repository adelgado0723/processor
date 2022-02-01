package processor

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"strings"
	"testing"

	"github.com/smartystreets/gunit"
)

func TestWriterHandlerFixture(t *testing.T) {
	gunit.Run(new(WriterHandlerFixture), t)
}

type WriterHandlerFixture struct {
	*gunit.Fixture
	handler *WriterHandler
	input   chan *Envelope
	buffer  *WriterSpyBuffer
	writer  *csv.Writer
}

func (whf *WriterHandlerFixture) Setup() {
	whf.buffer = NewWriterSpyBuffer("")
	whf.input = make(chan *Envelope, 10)
	whf.handler = NewWriterHandler(whf.input, whf.buffer)
}

func (whf *WriterHandlerFixture) TestHeaderWritten() {
	close(whf.input)

	whf.handler.Handle()
	whf.AssertEqual(whf.buffer.String(), "Status,DeliveryLine1,LastLine,City,State,ZIPCode\n")

}

func (whf *WriterHandlerFixture) TestOutputClosed() {
	close(whf.input)
	whf.handler.Handle()

	whf.Assert(whf.buffer.closed == 1)

}

func (whf *WriterHandlerFixture) TestEnvelopeWritten() {
	whf.sendEnvelope(1)
	whf.handler.Handle()

	lines := whf.outputLines()
	if whf.AssertEqual(len(lines), 2) {
		whf.AssertEqual(lines[1], "A1,B1,C1,D1,E1,F1")
	}
}

func (whf *WriterHandlerFixture) TestAllEnvelopesWritten() {
	whf.sendEnvelope(2)
	whf.handler.Handle()

	lines := whf.outputLines()
	if whf.AssertEqual(len(lines), 3) {
		whf.AssertEqual(lines[1], "A1,B1,C1,D1,E1,F1")
		whf.AssertEqual(lines[2], "A2,B2,C2,D2,E2,F2")
	}
}

var recordMatchingHeader = AddressOutput{
	Status:        "Status",
	DeliveryLine1: "DeliveryLine1",
	LastLine:      "LastLine",
	City:          "City",
	State:         "State",
	ZIPCode:       "ZIPCode",
}

func (whf *WriterHandlerFixture) TestHeaderMatchesRecords() {
	whf.input <- &Envelope{
		Output: recordMatchingHeader,
	}
	close(whf.input)
	whf.handler.Handle()
	whf.assertHeaderMatchesRecord()
}

func (whf *WriterHandlerFixture) assertHeaderMatchesRecord() {
	lines := whf.outputLines()
	whf.AssertEqual(lines[0], "Status,DeliveryLine1,LastLine,City,State,ZIPCode")
	whf.AssertEqual(lines[0], lines[1])
}
func (whf *WriterHandlerFixture) outputLines() []string {

	outputFile := strings.TrimSpace(whf.buffer.String())
	return strings.Split(outputFile, "\n")
}

func (whf *WriterHandlerFixture) sendEnvelope(count int) {
	for i := 1; i < count+1; i++ {
		whf.input <- &Envelope{
			Output: createOutput(strconv.Itoa(i)),
		}
	}
	close(whf.input)
}
func createOutput(index string) AddressOutput {
	return AddressOutput{
		Status:        "A" + index,
		DeliveryLine1: "B" + index,
		LastLine:      "C" + index,
		City:          "D" + index,
		State:         "E" + index,
		ZIPCode:       "F" + index,
	}
}

////////////////////////////////////////////////////////////////
// Create a Spy Buffer that counts how many times close() was called
type WriterSpyBuffer struct {
	// this syntax allows SpyBuffer to have the same Read/Close/etc
	// functionality that bytes.Buffer has without having to implement those methods
	*bytes.Buffer
	closed int
}

func (sb *WriterSpyBuffer) Close() error {
	sb.closed++
	// sb.Buffer.Reset()
	return nil
}

func NewWriterSpyBuffer(value string) *WriterSpyBuffer {
	return &WriterSpyBuffer{Buffer: bytes.NewBufferString(value)}
}
