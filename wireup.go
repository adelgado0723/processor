package processor

import (
	"io"
)

type Handler interface {
	Handle()
}
type Wireup struct {
	reader io.ReadCloser
	writer io.WriteCloser
	client HTTPClient
}

func (w *Wireup) Build() Handler {
	verifyInput := make(chan *Envelope, 1024)
	reorderInput := make(chan *Envelope, 1024)
	writerInput := make(chan *Envelope, 1024)
}

func Configure(reader io.ReadCloser, writer io.WriteCloser, client HTTPClient) *Wireup {
	return &Wireup{
		reader: reader,
		writer: writer,
		client: client,
	}
}
