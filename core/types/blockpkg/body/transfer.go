package body

import "github.com/janmbaco/Saprocate/core/types/blockpkg/header"

type Transfer struct{
	From header.Key
	To header.Key
	Sign []byte
}
