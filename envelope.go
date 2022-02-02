package processor

const (
	initialSequenceValue = 0
)

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
		Sequence int
		EOF      bool
		Input    AddressInput
		Output   AddressOutput
	}
)
