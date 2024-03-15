package pkg

import (
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
)

type AdapaterInterface interface {
	GenerateFromMetaTaskFile(*MetaTaskFileStruct, *AdaptConfig) error
}

type Generator struct {
	metataskfile string
	l            *logrus.Logger
	adapters     []AdapaterInterface
}

type ScriptStruct struct {
	Name        string `json:"name"`
	Command     string `json:"script"`
	Description string `json:"description"`
}

type MetaTaskFileStruct struct {
	// currently the script is any map string
	// probably not very extensive...
	//Scripts map[string]ScriptStruct `json:"scripts"`
	Scripts []ScriptStruct `json:"scripts"`
}

func NewGenerator(
	l *logrus.Logger,
	metataskfile string,
	dryRun bool,
) *Generator {
	return &Generator{
		metataskfile: metataskfile,
		l:            l,
		adapters: []AdapaterInterface{
			&NpmAdapter{
				npmFile: "package.json",
				l:       l,
				dryRun:  dryRun,
			},
			NewMakefileAdapter(l, "Makefile",
				"", "", dryRun),
		},
	}
}

func (g *Generator) Generate() error {
	g.l.Info("Generating a new project from:", g.metataskfile)
	// check if the file exists
	// if it does, return an error
	if _, err := os.Stat(g.metataskfile); os.IsNotExist(err) {
		g.l.Error("File: ", g.metataskfile, " does not exist")
		return err
	}
	// read the file
	// if it fails, return an error
	fp, err := os.Open(g.metataskfile)
	if err != nil {
		g.l.Error("Error opening file: ", g.metataskfile)
		return err
	}
	defer fp.Close()
	// decode the file
	// if it fails, return an error
	var m MetaTaskFileStruct
	err = json.NewDecoder(fp).Decode(&m)
	if err != nil {
		g.l.Error("Error decoding file: ", g.metataskfile)
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