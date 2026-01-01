package iam

// ActionType represents the type of action being performed
type ActionType string

const (
	// Namespace actions
	ActionNamespaceCreate ActionType = "namespace:Create"
	ActionNamespaceRead   ActionType = "namespace:Read"
	ActionNamespaceUpdate ActionType = "namespace:Update"
	ActionNamespaceDelete ActionType = "namespace:Delete"
	ActionNamespaceList   ActionType = "namespace:List"

	// Project actions
	ActionProjectCreate ActionType = "project:Create"
	ActionProjectRead   ActionType = "project:Read"
	ActionProjectUpdate ActionType = "project:Update"
	ActionProjectDelete ActionType = "project:Delete"
	ActionProjectList   ActionType = "project:List"
)

func (a ActionType) String() string {
	return string(a)
}
