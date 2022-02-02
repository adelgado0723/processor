package processor

const (
	initialSequenceValue = 0
	eofSequenceValue     = -1
)

var endOfFile = &Envelope{Sequence: eofSequenceValue}

type (
	AddressInput struct {
		Street1 string
		City    string
		State   string
		ZIPCode string
	}
	AddressOutput struct {
		Status        string // Is the result valid?
		DeliveryLine1 string
		LastLine      string
		City          string
		State         string
		ZIPCode       string
	}
	Envelope struct {
		Input    AddressInput
		Output   AddressOutput
		Sequence int
	}
)
