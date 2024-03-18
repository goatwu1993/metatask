package pkg

import (
	"fmt"
	"os"
	"strings"

	"metatask/pkg/adapters"
	"metatask/pkg/schema"

	"github.com/sirupsen/logrus"
)

type Generator struct {
	metatask string
	parsor   *V1YamlParsor
	l        *logrus.Logger
	adapters []adapters.AdapaterInterface
}

func NewGenerator(
	l *logrus.Logger,
	metatask string,
	dryRun bool,
) *Generator {
	return &Generator{
		metatask: metatask,
		l:        l,
		adapters: []adapters.AdapaterInterface{},
	}
}

func (g *Generator) AddAdapter(a adapters.AdapaterInterface) {
	g.adapters = append(g.adapters, a)
}

func (g *Generator) Generate() error {
	g.l.Info("Generating a new project from:", g.metatask)
	// check if the file exists
	// if it does, return an error
	// for all of the adapters, generate the project
	var m schema.FileRoot
	fileReader, err := os.Open(g.metatask)
	if err != nil {
		g.l.Error("Error opening file: ", err)
		return err
	}
	defer fileReader.Close()
	g.parsor = NewV1YamlParsorm(g.l)
	g.parsor.Parse(fileReader, &m, &ParsorConfig{})
	for k, _ := range m.Tasks {
		g.l.Info("Adding task: ", k)
	}
	if len(g.adapters) == 0 {
		g.l.Info("No adapters found, auto generating targets")
		err = g.AutoTargetWithGivenRoot(&m)
		if err != nil {
			return err
		}
	}

	for _, a := range g.adapters {
		err := a.GenerateFromMetaTaskFile(&m, &adapters.AdaptConfig{
			IgnoreNotFound: false,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) AutoTargetWithGivenRoot(root *schema.FileRoot) error {
	// for each of the adapters, check if the adapter has a target
	// if it does, add it to the root
	for _, s := range root.Syncs {
		switch strings.ToLower(s.FileType) {
		case "makefile":
			g.AddAdapter(adapters.NewMakefileAdapter(
				g.l,
				s.FilePath,
				false,
				"",
				"",
			))
		case "npm":
			g.AddAdapter(adapters.NewNpmAdapter(
				g.l,
				s.FilePath,
				false,
			))
		default:
			return fmt.Errorf("unknown file type: %s", s.FileType)
		}
	}
	return nil
}