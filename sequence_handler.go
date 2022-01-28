package processor

type SequenceHandler struct {
	input       chan *Envelope
	output      chan *Envelope
	envelopeIdx int
	buffer      map[int]*Envelope
}

func NewSequenceHandler(input, output chan *Envelope) *SequenceHandler {
	return &SequenceHandler{
		input:  input,
		output: output,
		buffer: make(map[int]*Envelope),
	}
}

func (sh *SequenceHandler) sendBufferedEnvelopesInOrder() {
	_, found := sh.buffer[sh.envelopeIdx]
	for ; found; _, found = sh.buffer[sh.envelopeIdx] {
		sh.output <- sh.buffer[sh.envelopeIdx]
		delete(sh.buffer, sh.envelopeIdx)
		sh.envelopeIdx++
	}
}
func (sh *SequenceHandler) processEnvelope(envelope *Envelope) {
	sh.buffer[envelope.Sequence] = envelope
	sh.sendBufferedEnvelopesInOrder()
}

func (sh *SequenceHandler) Handle() {
	for envelope := range sh.input {
		sh.processEnvelope(envelope)
	}
	close(sh.output)
}
