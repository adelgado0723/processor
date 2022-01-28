package processor

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestWriterHandlerFixture(t *testing.T) {
	gunit.Run(new(WriterHandlerFixture), t)
}

type WriterHandlerFixture struct {
	*gunit.Fixture
	handler *WriterHandler
	file *
}

func (wh *WriterHandlerFixture) Setup() {}
func (wh *WriterHandlerFixture) Test() {
	wh.handler = &WriterHandler{}
}
