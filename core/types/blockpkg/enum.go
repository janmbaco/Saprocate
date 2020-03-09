package blockpkg

type BlockType byte

const (
	Origin BlockType = iota
	Positive
	Negative
	Consumption
	Voucher
)

type PrevHashType byte

const (
	FirstPrevHash PrevHashType = iota
	SecondPrvHash
)
