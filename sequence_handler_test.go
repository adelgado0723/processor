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
	shf.sendEnvelopesInSequence(0, 1, 2, 3, 4)
	shf.handler.Handle()

	shf.assertSequenceOrder(shf.sequenceOrder(), []int{0, 1, 2, 3})
	shf.AssertEqual(len(shf.handler.buffer), 0)
}

func maxInt(sequences []int) (max int) {
	max = 0
	for _, sequenceNum := range sequences {
		if sequenceNum > max {
			max = sequenceNum
		}
	}
	return max
}
func (shf *SequenceHandlerFixture) sendEnvelopesInSequence(sequences ...int) {
	max := maxInt(sequences)
	for _, sequence := range sequences {
		shf.input <- &Envelope{Sequence: sequence, EOF: sequence == max}
	}
	// close(shf.input)
}

func (shf *SequenceHandlerFixture) TestEvelopeReceivedOutOfOrder_BufferedUntilContiguousBlock() {
	// Invalid order
	shf.sendEnvelopesInSequence(4, 2, 0, 3, 1)
	shf.handler.Handle()

	shf.assertSequenceOrder(shf.sequenceOrder(), []int{0, 1, 2, 3})
	// Checking map deleted envelopes after processing
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
