package pkg

import (
	// yaml
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/sirupsen/logrus"
)

type AdapaterInterface interface {
	GenerateFromMetaTaskFile(*MetaTaskFileStruct, *AdaptConfig) error
}

type Generator struct {
	metatask string
	l            *logrus.Logger
	adapters     []AdapaterInterface
}

type ScriptStruct struct {
	Name        string `yaml:"name"`
	Command     string `yaml:"script"`
	Description string `yaml:"description"`
}

type MetaTaskFileStruct struct {
	// currently the script is any map string
	// probably not very extensive...
	//Scripts map[string]ScriptStruct `yaml:"scripts"`
	Scripts []ScriptStruct `yaml:"scripts"`
}

func NewGenerator(
	l *logrus.Logger,
	metatask string,
	dryRun bool,
) *Generator {
	return &Generator{
		metatask: metatask,
		l:            l,
		adapters:     []AdapaterInterface{},
	}
}

func (g *Generator) AddAdapter(a AdapaterInterface) {
	g.adapters = append(g.adapters, a)
}

func (g *Generator) Generate() error {
	g.l.Info("Generating a new project from:", g.metatask)
	// check if the file exists
	// if it does, return an error
	if _, err := os.Stat(g.metatask); os.IsNotExist(err) {
		g.l.Error("File: ", g.metatask, " does not exist")
		return err
	}
	// read the file
	// if it fails, return an error
	fp, err := os.Open(g.metatask)
	if err != nil {
		g.l.Error("Error opening file: ", g.metatask)
		return err
	}
	defer fp.Close()
	// decode the file
	// if it fails, return an error
	var m MetaTaskFileStruct
	//err = json.NewDecoder(fp).Decode(&m)
	err = yaml.NewDecoder(fp).Decode(&m)
	if err != nil {
		g.l.Error("Error decoding file: ", g.metatask)
		return err
	}
	// print the structure
	g.l.Info("Scripts: ", m.Scripts)
	// for all of the adapters, generate the project
	for _, a := range g.adapters {
		a.GenerateFromMetaTaskFile(&m, &AdaptConfig{
			IgnoreNotFound: false,
		})
	}

	return nil
}