package blockpkg

type Type byte

const(
	Origin Type = iota
	Positive
	Negative
	Transfer
	Pay
)
