package out

type ByteRandGen interface {
	Generate() ([]byte, error)
	GenerateN(n int) ([]byte, error)
}
