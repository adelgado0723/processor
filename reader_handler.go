package processor

import (
	"encoding/csv"
	"io"
)

type ReaderHandler struct {
	reader *csv.Reader
	closer io.Closer
	output chan *Envelope
}

func NewReaderHandler(reader io.ReadCloser, output chan *Envelope) *ReaderHandler {
	return &ReaderHandler{
		reader: csv.NewReader(reader),
		closer: reader,
		output: output,
	}
}

func (rh *ReaderHandler) skipHeader() {
	rh.reader.Read()
}
func (rh *ReaderHandler) Handle() {
	rh.skipHeader()

	for {
		record, err := rh.reader.Read()
		if err != nil {
			break
		}
		rh.output <- &Envelope{
			Input: createInput(record),
		}
	}
}

func createInput(record []string) AddressInput {
	return AddressInput{
		Street1: record[0],
		City:    record[1],
		State:   record[2],
		ZIPCode: record[3],
	}
}
