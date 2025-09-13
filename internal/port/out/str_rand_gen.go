package out

type StrRandGen interface {
	Generate() (string, error)
	GenerateN(n int) (string, error)
}
