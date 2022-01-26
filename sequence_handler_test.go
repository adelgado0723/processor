package processor

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestSequenceHandler(t *testing.T) {
	gunit.Run(new(SequenceHandlerFixture), t)
}

type SequenceHandlerFixture struct {
	*gunit.Fixture
	input   chan *Envelope
	output  chan *Envelope
	handler *SequenceHandler
}

func (shf *SequenceHandlerFixture) Setup() {
	shf.input = make(chan *Envelope, 10)
	shf.output = make(chan *Envelope, 10)
	shf.handler = NewSequenceHandler(shf.input, shf.output)
}

func (shf *SequenceHandlerFixture) TestExpectedEnvelopeSentToOutput() {
	envelope := &Envelope{Sequence: 0}
	shf.input <- envelope
	shf.handler.Handle()

	shf.AssertEqual(<-shf.output, envelope)
}

func (shf *SequenceHandlerFixture) TestEvelopeReceivedOutOfOrder_BufferedUntilContiguousBlock() {
	envelope0 := &Envelope{Sequence: 0}
	envelope1 := &Envelope{Sequence: 1}

	// Invalid order
	shf.input <- envelope1
	shf.input <- envelope0

	shf.handler.Handle()

	// Expecting envelope0 first despite invalid order
	shf.AssertEqual(<-shf.output, envelope0)
	shf.AssertEqual(<-shf.output, envelope1)

}
