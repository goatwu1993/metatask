package adapters

import (
	"encoding/json"
	"fmt"
	"os"

	"metatask/pkg/schema"

	"github.com/sirupsen/logrus"
)

type NpmAdapter struct {
	l       *logrus.Logger
	npmFile string
	dryRun  bool
}

func NewNpmAdapter(l *logrus.Logger, npmFile string, dryRun bool) *NpmAdapter {
	return &NpmAdapter{
		l:       l,
		npmFile: npmFile,
		dryRun:  dryRun,
	}
}

type NpmPackageJson struct {
	Scripts map[string]string `json:"tasks"`
	// preserve any other fields
}

func (n *NpmAdapter) GenerateFromMetaTaskFile(m *schema.TreeRoot, c *AdaptConfig) error {
	n.l.Debug("updating npm file: ", n.npmFile)
	// Check if the file exists
	if _, err := os.Stat(n.npmFile); os.IsNotExist(err) {
		if c.IgnoreNotFound {
			n.l.Debug("Ignoring not found")
			return nil
		}
		n.l.Error("File: ", n.npmFile, " does not exist")
		return err
	}

	// Read the existing package.json file
	file, err := os.ReadFile(n.npmFile)
	if err != nil {
		n.l.Error("Error reading file: ", n.npmFile)
		return err
	}

	// Decode JSON into a map
	var pkgMap map[string]interface{}
	if err := json.Unmarshal(file, &pkgMap); err != nil {
		n.l.Error("Error decoding JSON: ", err)
		return err
	}

	// Update the tasks field
	if tasks, ok := pkgMap["tasks"].(map[string]interface{}); ok {
		for _, v := range m.Tasks {
			tasks[v.Name] = v.Command
		}
		pkgMap["tasks"] = tasks
	} else {
		// If tasks field doesn't exist, create it
		pkgMap["tasks"] = m.Tasks
	}

	// Encode the updated map back to JSON
	updatedJSON, err := json.MarshalIndent(pkgMap, "", "    ")
	if err != nil {
		n.l.Error("Error encoding JSON: ", err)
		return err
	}

	// Write the updated JSON content back to the file
	if n.dryRun {
		n.l.Debug("Dry run, not writing to file")
		fmt.Println(string(updatedJSON))
		return nil
	}
	n.l.Debug("Updated tasks in package.json")
	if err := os.WriteFile(n.npmFile, updatedJSON, 0644); err != nil {
		n.l.Error("Error writing to file: ", err)
		return err
	}
	return nil
}
