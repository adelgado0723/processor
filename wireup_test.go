package processor

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestWireupFixture(t *testing.T) {
	gunit.Run(new(WireupFixture), t)
}

type WireupFixture struct {
	*gunit.Fixture

	reader  *ReadWriteSpyBuffer
	writer  *ReadWriteSpyBuffer
	handler Handler
}

func (wf *WireupFixture) Setup() {
	wf.reader = NewReadWriteSpyBuffer("")
	wf.writer = NewReadWriteSpyBuffer("")
	wf.handler = Configure(wf.reader, wf.writer, nil).Build()

}
func (wf *WireupFixture) Test() {}
