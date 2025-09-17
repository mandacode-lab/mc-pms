package initialize

const (
	ErrInvalidRequestMsg           = "invalid request"
	ErrServiceNameRequiredMsg      = "service name is required"
	ErrClientNameRequiredMsg       = "client name is required"
	ErrInvalidServiceNameMsg       = "invalid service name format"
	ErrServiceCreationFailedMsg    = "failed to create service"
	ErrClientAppCreationFailedMsg  = "failed to create client application"
	ErrSecretGenerationFailedMsg   = "failed to generate client secret"
	ErrSecretHashingFailedMsg      = "failed to hash client secret"
	ErrInternalServerMsg           = "internal server error"
)