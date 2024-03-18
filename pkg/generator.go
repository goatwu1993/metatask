package pkg

import (
	"os"

	"github.com/sirupsen/logrus"
)

type AdapaterInterface interface {
	GenerateFromMetaTaskFile(*MetaTaskRoot, *AdaptConfig) error
}

type Generator struct {
	metatask string
	parsor   *V1YamlParsor
	l        *logrus.Logger
	adapters []AdapaterInterface
}

func NewGenerator(
	l *logrus.Logger,
	metatask string,
	dryRun bool,
) *Generator {
	return &Generator{
		metatask: metatask,
		l:        l,
		adapters: []AdapaterInterface{},
	}
}

func (g *Generator) AddAdapter(a AdapaterInterface) {
	g.adapters = append(g.adapters, a)
}

func (g *Generator) Generate() error {
	g.l.Info("Generating a new project from:", g.metatask)
	// check if the file exists
	// if it does, return an error
	// for all of the adapters, generate the project
	var m MetaTaskRoot
	fileReader, err := os.Open(g.metatask)
	if err != nil {
		g.l.Error("Error opening file: ", err)
		return err
	}
	defer fileReader.Close()
	g.parsor = NewV1YamlParsorm(g.l)
	g.parsor.Parse(fileReader, &m, &ParsorConfig{})
	// read the file

	for _, a := range g.adapters {
		a.GenerateFromMetaTaskFile(&m, &AdaptConfig{
			IgnoreNotFound: false,
		})
	}

	return nil
}