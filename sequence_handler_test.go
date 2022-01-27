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
	close(shf.input)
	shf.handler.Handle()

	shf.AssertEqual(<-shf.output, envelope)
}

func (shf *SequenceHandlerFixture) TestEvelopeReceivedOutOfOrder_BufferedUntilContiguousBlock() {
	// Invalid order
	shf.input <- &Envelope{Sequence: 4}
	shf.input <- &Envelope{Sequence: 1}
	shf.input <- &Envelope{Sequence: 2}
	shf.input <- &Envelope{Sequence: 3}
	shf.input <- &Envelope{Sequence: 0}

	close(shf.input)
	shf.handler.Handle()

	close(shf.output)

	shf.assertSequenceOrder(shf.sequenceOrder(), []int{0, 1, 2, 3, 4})
	// Expecting envelope0 first despite invalid order
	// shf.AssertEqual((<-shf.output).Sequence, 0)
	// shf.AssertEqual((<-shf.output).Sequence, 1)
	// shf.AssertEqual((<-shf.output).Sequence, 2)
	// shf.AssertEqual((<-shf.output).Sequence, 3)
	// shf.AssertEqual((<-shf.output).Sequence, 4)
	shf.AssertEqual(len(shf.handler.buffer), 0)
}

func (shf *SequenceHandlerFixture) sequenceOrder() (order []int) {
	for envelope := range shf.output {
		order = append(order, envelope.Sequence)
	}
	return order
}
func (shf *SequenceHandlerFixture) assertSequenceOrder(actual, expected []int) {

	shf.AssertEqual(len(actual), len(expected))
	for index, seqNo := range actual {
		shf.AssertEqual(expected[index], seqNo)
	}
}
