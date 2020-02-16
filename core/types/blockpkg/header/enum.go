package header

type Type byte

const(
	Origin Type = iota
	Numu
	Transfer
	Pay
)
