package processor

import (
	"strings"
	"testing"

	"github.com/smartystreets/gunit"
)

func TestHandleFixture(t *testing.T) {
	gunit.Run(new(HandlerFixture), t)
}

type HandlerFixture struct {
	*gunit.Fixture

	input       chan *Envelope
	output      chan *Envelope
	application *FakeVerifier
	handler     *VerifyHandler
}

func (hf *HandlerFixture) Setup() {
	hf.input = make(chan *Envelope, 10)
	hf.output = make(chan *Envelope, 10)
	hf.application = NewFakeVerifier()
	hf.handler = NewVerifyHandler(hf.input, hf.output, hf.application)
}

func (hf *HandlerFixture) enqueueEnvelope(street1 string) *Envelope {
	envelope := &Envelope{
		Input: AddressInput{
			Street1: street1,
		},
	}
	hf.input <- envelope
	return envelope
}

func (hf *HandlerFixture) TestVerifierReceivesInput() {
	envelope := hf.enqueueEnvelope("street")
	close(hf.input)

	hf.handler.Handle()

	hf.AssertEqual(envelope, <-hf.output)
	hf.AssertEqual("STREET", envelope.Output.DeliveryLine1)
}

func (hf *HandlerFixture) TestInputQueueDrained() {
	envelope1 := hf.enqueueEnvelope("41")
	envelope2 := hf.enqueueEnvelope("42")
	envelope3 := hf.enqueueEnvelope("43")

	close(hf.input)
	hf.handler.Handle()
	hf.AssertEqual(envelope1, <-hf.output)
	hf.AssertEqual(envelope2, <-hf.output)
	hf.AssertEqual(envelope3, <-hf.output)
}

///////////////////////////////////////////////////////////
type FakeVerifier struct {
	input  AddressInput
	output AddressOutput
}

func NewFakeVerifier() *FakeVerifier { return &FakeVerifier{} }
func (fv *FakeVerifier) Verify(value AddressInput) AddressOutput {
	fv.input = value
	return AddressOutput{DeliveryLine1: strings.ToUpper(value.Street1)}
}
