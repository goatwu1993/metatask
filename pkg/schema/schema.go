package schema

type MetaTaskScript struct {
	Name        string `yaml:"name"`
	Command     string `yaml:"script"`
	Description string `yaml:"description"`
}

type SyncTarget struct {
	FileType string `yaml:"fileType"`
	FilePath string `yaml:"filePath"`
}

type MetaTaskRoot struct {
	// currently the script is any map string
	// probably not very extensive...
	//Scripts map[string]ScriptStruct `yaml:"scripts"`
	Scripts []MetaTaskScript `yaml:"scripts"`
	Syncs   []SyncTarget     `yaml:"syncs"`
}
