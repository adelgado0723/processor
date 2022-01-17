package processor

type VerifyHandler struct {
	input       chan *Envelope
	output      chan *Envelope
	application Verifier
}

type Verifier interface {
	Verify(AddressInput) AddressOutput
}

func NewVerifyHandler(in, out chan *Envelope, verifier Verifier) *VerifyHandler {
	return &VerifyHandler{
		input:       in,
		output:      out,
		application: verifier,
	}
}
func (vh *VerifyHandler) Handle() {
	for envelope := range vh.input {
		envelope.Output = vh.application.Verify(envelope.Input)
		vh.output <- envelope
	}

}
