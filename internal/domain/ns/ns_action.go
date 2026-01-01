package ns

import "time"

func (n *Namespace) UpdateName(name NamespaceName) {
	n.name = name
	n.updatedAt = time.Now()
}

func (n *Namespace) UpdateDescription(desc NamespaceDescription) {
	n.desc = desc
	n.updatedAt = time.Now()
}
