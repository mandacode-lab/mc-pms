package out

type Encoder interface {
	Encode(data []byte) ([]byte, error)
	Decode(encoded []byte) ([]byte, error)
}
