package schema

type TreeTask struct {
	Name        string
	Command     string
	Description string
	DependsOn   []*TreeTask
	Visited     bool
	InStack     bool // Used for cycle detection
}

type TreeRoot struct {
	Tasks []*TreeTask `yaml:"tasks"`
}
