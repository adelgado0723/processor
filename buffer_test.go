package processor

import (
	"encoding/csv"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/smartystreets/gunit"
)

func TestSpyBuffer(t *testing.T) {
	gunit.Run(new(SpyBufferFixture), t)
}

type SpyBufferFixture struct {
	*gunit.Fixture
}

func (sbf *SpyBufferFixture) Setup() {}

func (sbf *SpyBufferFixture) Test() {
	buffer := NewReadWriteSpyBuffer("")
	buffer.WriteString("Hello, World!")

	sbf.AssertEqual(buffer.String(), "Hello, World!")

	raw, err := ioutil.ReadAll(buffer)
	sbf.Assert(err == nil)
	sbf.AssertEqual(string(raw), "Hello, World!")

	reader := csv.NewReader(strings.NewReader("Hello, World!"))
	record, err2 := reader.Read()
	sbf.Assert(err2 == nil)

	sbf.AssertDeepEqual(record, []string{"Hello", " World!"})
}
