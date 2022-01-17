package processor

type VerifyHandler struct {
	input       chan *Envelope
	output      chan *Envelope
	application Verifier
}

type Verifier interface {
	Verify(AddressInput)
}

func NewVerifyHandler(in, out chan *Envelope, application Verifier) *VerifyHandler {
	return &VerifyHandler{
		input:       in,
		output:      out,
		application: application,
	}
}
func (vh *VerifyHandler) Handle() {
	received := <-vh.input
	vh.application.Verify(received.Input)
	vh.output <- received

}
