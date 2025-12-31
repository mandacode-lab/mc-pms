package project

func (p *Project) UpdateName(name ProjectName) {
	p.name = name
	p.updatedAt = time.Now()
}

func (p *Project) UpdateDescription(desc ProjectDescription) {
	p.description = desc
	p.updatedAt = time.Now()
}

func (p *Project) UpdateNamespaceID(namespaceID ns.NamespaceID) {
	p.namespaceID = namespaceID
	p.updatedAt = time.Now()
}
