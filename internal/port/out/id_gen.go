package out

type IDGenerator interface {
	Generate() (string, error)
}
