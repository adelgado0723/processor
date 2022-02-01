package processor

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestReaderHandlerFixture(t *testing.T) {
	gunit.Run(new(ReaderHandlerFixture), t)

}

type ReaderHandlerFixture struct {
	*gunit.Fixture
}

func (rhf *ReaderHandlerFixture) Setup() {}
func (rhf *ReaderHandlerFixture) TestCSVRecordSentInEnvelope() {
	buffer := NewReadWriteSpyBuffer("Street1,City,State,ZIPCode")

	output := make(chan *Envelope, 10)

	reader := NewReaderHandler(buffer, output)
	reader.Handle()

	rhf.AssertEqual(Envelope(<-output).Output)
}
