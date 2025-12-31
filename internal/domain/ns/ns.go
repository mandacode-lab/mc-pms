package ns

import "time"

type Namespace struct {
	id        NamespaceID
	name      NamespaceName
	desc      NamespaceDescription
	createdAt time.Time
	updatedAt time.Time
}

func NewNamespace(
	id NamespaceID,
	name NamespaceName,
	desc NamespaceDescription,
	createdAt time.Time,
	updatedAt time.Time,
) *Namespace {
	return &Namespace{
		id:        id,
		name:      name,
		desc:      desc,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (n *Namespace) ID() NamespaceID                   { return n.id }
func (n *Namespace) Name() NamespaceName               { return n.name }
func (n *Namespace) Description() NamespaceDescription { return n.desc }
func (n *Namespace) CreatedAt() time.Time              { return n.createdAt }
func (n *Namespace) UpdatedAt() time.Time              { return n.updatedAt }
