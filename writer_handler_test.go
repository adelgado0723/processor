package processor

import (
	"bytes"
	"encoding/csv"
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
	whf.input = make(chan *Envelope)
	whf.handler = NewWriterHandler(whf.input, whf.buffer)
}

func (whf *WriterHandlerFixture) TestHeaderWritten() {
	close(whf.input)

	whf.handler.Handle()
	whf.AssertEqual(whf.buffer.String(), "Status,DeliveryLine1,City,State,ZIPCode\n")

}

func (whf *WriterHandlerFixture) TestOutputClosed() {
	close(whf.input)
	whf.handler.Handle()

	whf.Assert(whf.buffer.closed == 1)

}

func (whf *WriterHandlerFixture) TestEnvelopeWritten() {
	whf.input <- &Envelope{
		Output: AddressOutput{
			Status:        "A",
			DeliveryLine1: "B",
			City:          "C",
			State:         "D",
			ZIPCode:       "E",
			LastLine:      "F",
		}}
	close(whf.input)
	whf.handler.Handle()

	outputFile := strings.TrimSpace(whf.buffer.String())
	lines := strings.Split(outputFile, "\n")

	if whf.AssertEqual(len(lines), 2) {
		whf.AssertEqual(lines[1], "A,B,C,D,E,F")
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
