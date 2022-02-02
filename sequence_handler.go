package processor

type SequenceHandler struct {
	input   chan *Envelope
	output  chan *Envelope
	counter int
	buffer  map[int]*Envelope
}

func NewSequenceHandler(input, output chan *Envelope) *SequenceHandler {
	return &SequenceHandler{
		input:   input,
		output:  output,
		counter: initialSequenceValue,
		buffer:  make(map[int]*Envelope),
	}
}

func (sh *SequenceHandler) sendBufferedEnvelopesInOrder() {
	_, found := sh.buffer[sh.counter]
	for ; found; _, found = sh.buffer[sh.counter] {
		sh.output <- sh.buffer[sh.counter]
		delete(sh.buffer, sh.counter)
		sh.counter++
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
