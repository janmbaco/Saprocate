package block

type Interface interface{
	SerializeKey() []byte
	SerializeValue() []byte
}
