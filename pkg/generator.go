package pkg

import (
	"github.com/sirupsen/logrus"
)

type AdapaterInterface interface {
	GenerateFromMetaTaskFile(*MetaTaskRoot, *AdaptConfig) error
}

type Generator struct {
	metatask string
	parsor   *YamlParsor
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
	g.parsor = NewYamlParsor(g.metatask, g.l)
	g.parsor.Parse(&m, &ParsorConfig{})
	// read the file

	for _, a := range g.adapters {
		a.GenerateFromMetaTaskFile(&m, &AdaptConfig{
			IgnoreNotFound: false,
		})
	}

	return nil
}