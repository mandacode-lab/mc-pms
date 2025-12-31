package ns

type NamespaceID string

func NewNamespaceID(id string) NamespaceID {
	return NamespaceID(id)
}

func (n NamespaceID) String() string {
	return string(n)
}

type NamespaceName string

func NewNamespaceName(name string) NamespaceName {
	return NamespaceName(name)
}

func (n NamespaceName) String() string {
	return string(n)
}

type NamespaceDescription string

func NewNamespaceDescription(description string) NamespaceDescription {
	return NamespaceDescription(description)
}

func (n NamespaceDescription) String() string {
	return string(n)
}
