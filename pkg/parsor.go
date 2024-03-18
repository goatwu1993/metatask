package pkg

import (
	// yaml
	"io"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type ParsorConfig struct{}

type ParsorInterface interface {
	ParsefromMetaTaskFile(r *MetaTaskRoot, c *ParsorConfig) error
}

type V1YamlParsor struct {
	l *logrus.Logger
}

func NewV1YamlParsorm(
	l *logrus.Logger,
) *V1YamlParsor {
	return &V1YamlParsor{
		l: l,
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

func (p *V1YamlParsor) Parse(reader io.Reader, r *MetaTaskRoot, c *ParsorConfig) error {
	// check if the file exists
	// if it does, return an error
	// if it fails, return an error
	//err = json.NewDecoder(fp).Decode(&m)
	err := yaml.NewDecoder(reader).Decode(&r)
	if err != nil {
		p.l.Error("Error decoding file: ", err)
		return err
	}
	return nil
}