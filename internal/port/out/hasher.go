package out

type Hasher interface {
	Hash(data []byte) ([]byte, error)
	Compare(hash, data []byte) error
}
