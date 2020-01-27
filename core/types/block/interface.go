package block

type Interface interface{
	GetType() Type
	SerializeKey() []byte
	SerializeValue() []byte
}
