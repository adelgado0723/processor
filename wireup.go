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

func (w *Wireup) Build() {

}

func Configure(reader io.ReadCloser, writer io.WriteCloser, client HTTPClient) *Wireup {
	return &Wireup{
		reader: reader,
		writer: writer,
		client: client,
	}
}
