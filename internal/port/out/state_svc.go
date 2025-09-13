package out

type StateService interface {
	GenerateState() (string, error)
	ValidateState(state string) bool
}
