package processor

import (
	"io"
)

type Handler interface {
	Handle()
}
type Pipeline struct {
	reader  io.ReadCloser
	writer  io.WriteCloser
	client  HTTPClient
	workers int
	err     error
}

func Configure(reader io.ReadCloser, writer io.WriteCloser, client HTTPClient, workers int) *Pipeline {
	return &Pipeline{
		reader:  reader,
		writer:  writer,
		client:  client,
		workers: workers,
	}
}
func (p *Pipeline) Process() (err error) {
	verifyInput := make(chan *Envelope, 1024)
	sequenceInput := make(chan *Envelope, 1024)
	writerInput := make(chan *Envelope, 1024)
	verifier := NewSmartyVerifier(p.client)
	for i := 0; i < p.workers; i++ {
		go NewVerifyHandler(verifyInput, sequenceInput, verifier).Handle()
	}

	go func() { err = NewReaderHandler(p.reader, verifyInput).Handle() }()
	go NewSequenceHandler(sequenceInput, writerInput).Handle()
	NewWriterHandler(writerInput, p.writer).Handle()
	return err
}