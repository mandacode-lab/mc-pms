package iam

import "context"

// AuthorizationRequest represents a request to check authorization
type AuthorizationRequest struct {
	Token    string
	Action   ActionType
	Resource Resource
}

func NewAuthorizationRequest(token string, action ActionType, resource Resource) *AuthorizationRequest {
	return &AuthorizationRequest{
		Token:    token,
		Action:   action,
		Resource: resource,
	}
}

// AuthorizationResult represents the result of an authorization check
type AuthorizationResult struct {
	Allowed bool
	Reason  string // Reason for denial if not allowed
}

// IAMService is the port interface for IAM operations
// IAM service handles both token verification and authorization
type IAMService interface {
	// CheckAuthorization checks if the token holder is allowed to perform the action on the resource
	// Returns authorization result with allowed status and reason
	// Returns error only if token verification or service call fails
	CheckAuthorization(ctx context.Context, req *AuthorizationRequest) (*AuthorizationResult, error)
}
