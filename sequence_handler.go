package processor

type SequenceHandler struct {
	input  chan *Envelope
	output chan *Envelope
}

func NewSequenceHandler(input, output chan *Envelope) *SequenceHandler {
	return &SequenceHandler{
		input:  input,
		output: output,
	}
}

func (sh *SequenceHandler) Handle() {
	for i := 0; i < 2; i++ {
		input := <-sh.input
		sh.output <- input
	}
}
