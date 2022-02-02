package processor

import (
	"strconv"
	"testing"

	"github.com/smartystreets/gunit"
)

func TestReaderHandlerFixture(t *testing.T) {
	gunit.Run(new(ReaderHandlerFixture), t)

}

type ReaderHandlerFixture struct {
	*gunit.Fixture
	buffer *ReadWriteSpyBuffer
	output chan *Envelope
	reader *ReaderHandler
}

func (rhf *ReaderHandlerFixture) Setup() {
	rhf.buffer = NewReadWriteSpyBuffer("")
	rhf.output = make(chan *Envelope, 10)
	rhf.reader = NewReaderHandler(rhf.buffer, rhf.output)
	rhf.writeLine("Street1,City,State,ZIPCode")
}

func (rhf *ReaderHandlerFixture) writeLine(line string) {
	rhf.buffer.WriteString(line + "\n")
}

func (rhf *ReaderHandlerFixture) TestAllCSVRecordsSentToOutput() {
	rhf.writeLine("A1,B1,C1,D1")
	rhf.writeLine("A2,B2,C2,D2")
	rhf.reader.Handle()

	rhf.assertRecordsSent()
	rhf.assertCleanup()
}

func (rhf *ReaderHandlerFixture) assertRecordsSent() {
	rhf.AssertDeepEqual(<-rhf.output, buildEnvelope(initialSequenceValue))
	rhf.AssertDeepEqual(<-rhf.output, buildEnvelope(initialSequenceValue+1))

}
func (rhf *ReaderHandlerFixture) assertCleanup() {
	rhf.AssertEqual(<-rhf.output, endOfFile)
	rhf.Assert(<-rhf.output == nil)
	rhf.AssertEqual(rhf.buffer.closed, 1)
}

func buildEnvelope(index int) *Envelope {
	suffix := strconv.Itoa(index + 1)
	return &Envelope{
		Sequence: index,
		Input: AddressInput{
			Street1: "A" + suffix,
			City:    "B" + suffix,
			State:   "C" + suffix,
			ZIPCode: "D" + suffix,
		},
	}
}

func (rhf *ReaderHandlerFixture) TestMalformedInputReturnsError() {}
