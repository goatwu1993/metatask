package pkg

import (
	"fmt"
	"os"
	"strings"

	"metatask/pkg/adapters"
	"metatask/pkg/schema"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
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
	g.l.Debug("Generating a new project from:", g.metatask)
	// check if the file exists
	// if it does, return an error
	// for all of the adapters, generate the project
	// init an empty tree root
	fr := schema.FileRoot{}
	tr := schema.TreeRoot{}
	fileReader, err := os.Open(g.metatask)
	if err != nil {
		g.l.Error("Error opening file: ", err)
		return err
	}
	defer fileReader.Close()
	g.parsor = NewV1YamlParsor(g.l)
	err = g.parsor.Parse(fileReader, &tr, &fr, &ParsorConfig{})
	if err != nil {
		return err
	}

	// dump the tree root to metatask-lock.yaml
	file, err := os.Create("metatask-lock.yaml")
	if err != nil {
		g.l.Error("Error creating file: ", err)
		return err
	}
	defer file.Close()
	err = yaml.NewEncoder(file).Encode(&tr)
	if err != nil {
		g.l.Error("Error encoding file: ", err)
		return err
	}

	if len(g.adapters) == 0 {
		err = g.AutoTargetWithGivenRoot(&tr, &schema.FileRoot{})
		if err != nil {
			return err
		}
	}

	for _, a := range g.adapters {
		err := a.GenerateFromMetaTaskFile(&tr, &adapters.AdaptConfig{
			IgnoreNotFound: false,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) AutoTargetWithGivenRoot(tr *schema.TreeRoot, fr *schema.FileRoot) error {
	// for each of the adapters, check if the adapter has a target
	// if it does, add it to the root
	for _, s := range fr.Syncs {
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