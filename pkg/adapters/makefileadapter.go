package adapters

import (
	"fmt"
	"os"
	"strings"

	"metatask/pkg/schema"

	"github.com/sirupsen/logrus"
)

// install the go makefile parser by go get github.com/rogpeppe/go-internal/modfile

type MakefileAdapter struct {
	l            *logrus.Logger
	dryRun       bool
	makefilePath string
	startSection string
	endSection   string
}

func NewMakefileAdapter(l *logrus.Logger, makefilePath string, dryRun bool, startSection string, endSection string) *MakefileAdapter {
	if startSection == "" {
		startSection = "## metatask-start"
	}
	if endSection == "" {
		endSection = "## metatask-end"
	}
	return &MakefileAdapter{
		l:            l,
		makefilePath: makefilePath,
		startSection: startSection,
		endSection:   endSection,
		dryRun:       dryRun,
	}
}

func (ma *MakefileAdapter) GenerateFromMetaTaskFile(m *schema.FileRoot, c *AdaptConfig) error {
	ma.l.Info("updating makefile: ", ma.makefilePath)
	// Check if the file exists
	if _, err := os.Stat(ma.makefilePath); os.IsNotExist(err) {
		//if c.IgnoreNotFound {
		ma.l.Debug("Ignoring not found")
		ma.l.Info("Creating file: ", ma.makefilePath)
		_, err = os.Create(ma.makefilePath)
		if err != nil {
			ma.l.Error("Error creating file: ", ma.makefilePath)
			return err
		}
		// check if the file was created
		if _, err := os.Stat(ma.makefilePath); os.IsNotExist(err) {
			ma.l.Error("File was not created: ", ma.makefilePath)
			return err
		}
		// check is the file is a directory
		if fi, err := os.Stat(ma.makefilePath); err == nil && fi.IsDir() {
			ma.l.Error("File is a directory: ", ma.makefilePath)
			return err
		}
	}

	// Read the existing package.json file
	file, err := os.ReadFile(ma.makefilePath)
	if err != nil {
		ma.l.Error("Error reading file: ", ma.makefilePath)
		return err
	}
	// I cannot find a go library to parse makefile, so I will use a regex to see if there is a section for the metatask
	// if not append the metatask section to the end of the file
	// Look for ## metatask-start
	// Look for ## metatask-end
	// If both are found, replace the section in between with the new
	originalFile := string(file)
	start := strings.Index(originalFile, ma.startSection)
	end := strings.Index(originalFile, ma.endSection)
	section := ma.GenerateSection(m)
	if start == -1 || end == -1 {
		// append the metatask section to the end of the file
		ma.l.Debug("Appending metatask section to the end of the file")
		originalFile += section

	} else {
		// replace the section in between with the new
		ma.l.Debug("Replacing the section in between with the new")
		originalFile = originalFile[:start] + section + originalFile[end+len(ma.endSection):]
	}
	if ma.dryRun {
		ma.l.Info("Dry run, not writing to file")
		fmt.Println(originalFile)
	} else {
		err = os.WriteFile(ma.makefilePath, []byte(originalFile), 0644)
		if err != nil {
			ma.l.Error("Error writing file: ", ma.makefilePath)
			return err
		}
	}

	return nil
}

func (n *MakefileAdapter) GenerateSection(m *schema.FileRoot) string {
	// Generate the section
	allPhonies := ""
	// m.Scripts is a map[string]string
	for _, v := range m.Tasks {
		// ignore Makefile itself
		// user should probably avoid the circular dependency
		allPhonies += ".PHONY: " + v.Name + "\n" + v.Name + ": ## " + v.Description + "\n" + "\t" + v.Command + "\n\n"
	}
	return n.startSection + "\n" + allPhonies + n.endSection
}
