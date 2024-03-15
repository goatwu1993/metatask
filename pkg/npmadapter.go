package pkg

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type NpmAdapter struct {
	l       *logrus.Logger
	npmFile string
	dryRun  bool
}

type AdaptConfig struct {
	IgnoreNotFound bool
}

type NpmPackageJson struct {
	Scripts map[string]string `json:"scripts"`
	// preserve any other fields
}

func (n *NpmAdapter) GenerateFromMetaTaskFile(m *MetaTaskFileStruct, c *AdaptConfig) error {
	n.l.Info("updating npm file: ", n.npmFile)
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

	// Update the scripts field
	if scripts, ok := pkgMap["scripts"].(map[string]interface{}); ok {
		for _, v := range m.Scripts {
			// Only update the command part of the script
			scripts[v.Name] = v.Command
		}
		pkgMap["scripts"] = scripts
	} else {
		// If scripts field doesn't exist, create it
		pkgMap["scripts"] = m.Scripts
	}

	// Encode the updated map back to JSON
	updatedJSON, err := json.MarshalIndent(pkgMap, "", "    ")
	if err != nil {
		n.l.Error("Error encoding JSON: ", err)
		return err
	}

	// Write the updated JSON content back to the file
	if n.dryRun {
		n.l.Info("Dry run, not writing to file")
		fmt.Println(string(updatedJSON))
		return nil
	}
	n.l.Info("Updated scripts in package.json")
	if err := os.WriteFile(n.npmFile, updatedJSON, 0644); err != nil {
		n.l.Error("Error writing to file: ", err)
		return err
	}
	return nil
}
