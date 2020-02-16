package memory

type State uint8
const (
	None = iota
	Stored
	Transferred
	Paid
)
