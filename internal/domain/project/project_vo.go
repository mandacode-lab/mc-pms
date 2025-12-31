package project

type ProjectID string

func NewProjectID(id string) ProjectID {
	return ProjectID(id)
}

func (p ProjectID) String() string {
	return string(p)
}

type ProjectName string

func NewProjectName(name string) ProjectName {
	return ProjectName(name)
}

func (p ProjectName) String() string {
	return string(p)
}

type ProjectDescription string

func NewProjectDescription(description string) ProjectDescription {
	return ProjectDescription(description)
}

func (p ProjectDescription) String() string {
	return string(p)
}
