package out

type Encoder interface {
	Encode(data []byte) []byte
}
