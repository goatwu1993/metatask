package schema

type FileTask struct {
	Name        string   `yaml:"name"`
	Command     string   `yaml:"script"`
	Description string   `yaml:"description"`
	DependsOn   []string `yaml:"dependsOn"`
}

type FileSyncTarget struct {
	FileType string `yaml:"fileType"`
	FilePath string `yaml:"filePath"`
}

type FileRoot struct {
	Tasks []FileTask       `yaml:"tasks"`
	Syncs []FileSyncTarget `yaml:"syncs"`
}
