package pkg

import (
	// yaml
	"os"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type ParsorConfig struct{}

type ParsorInterface interface {
	ParsefromMetaTaskFile(*ParsorConfig) error
}

type YamlParsor struct {
	metataskyaml string
	l            *logrus.Logger
}

func NewYamlParsor(
	metataskyaml string,
	l *logrus.Logger,
) *YamlParsor {
	return &YamlParsor{
		metataskyaml: metataskyaml,
		l:            l,
	}
}

type MetaTaskScript struct {
	Name        string `yaml:"name"`
	Command     string `yaml:"script"`
	Description string `yaml:"description"`
}

type MetaTaskRoot struct {
	// currently the script is any map string
	// probably not very extensive...
	//Scripts map[string]ScriptStruct `yaml:"scripts"`
	Scripts []MetaTaskScript `yaml:"scripts"`
}

func (p *YamlParsor) Parse(r *MetaTaskRoot, c *ParsorConfig) error {
	p.l.Info("Generating a new project from:", p.metataskyaml)
	// check if the file exists
	// if it does, return an error
	if _, err := os.Stat(p.metataskyaml); os.IsNotExist(err) {
		p.l.Error("File: ", p.metataskyaml, " does not exist")
		return err
	}
	// read the file
	// if it fails, return an error
	fp, err := os.Open(p.metataskyaml)
	if err != nil {
		p.l.Error("Error opening file: ", p.metataskyaml)
		return err
	}
	defer fp.Close()
	// decode the file
	// if it fails, return an error
	//err = json.NewDecoder(fp).Decode(&m)
	err = yaml.NewDecoder(fp).Decode(&r)
	if err != nil {
		p.l.Error("Error decoding file: ", p.metataskyaml)
		return err
	}
	return nil
}