package processor

import (
	"encoding/csv"
	"io"
)

type WriterHandler struct {
	input  chan *Envelope
	closer io.Closer
	writer *csv.Writer
}

func NewWriterHandler(input chan *Envelope, output io.WriteCloser) *WriterHandler {
	return &WriterHandler{
		input:  input,
		closer: output,
		writer: csv.NewWriter(output),
	}

}
func (wh *WriterHandler) Handle() {
	wh.writer.Write([]string{"Status", "DeliveryLine1", "City", "State", "ZIPCode"})

	envelope := <-wh.input
	if envelope != nil {
		output := envelope.Output
		wh.writer.Write([]string{
			output.Status,
			output.DeliveryLine1,
			output.City,
			output.State,
			output.ZIPCode,
			output.LastLine,
		})
	}
	wh.writer.Flush()
	wh.closer.Close()
}
