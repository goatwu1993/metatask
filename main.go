package main

import (
	"fmt"
	"os"
	"runtime"

	pkg "metatask/pkg"

	"github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
)

const (
	DEFAULT_YAML_FILE = "metatask.yml"
)

func checkIfFilePathExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return false
	}
	return true
}

func checkInputFileOrDefaultYamlFileIsRegularFileAndGetFileName(inputFileName string) (string, error) {
	var innerInputFileName string
	if inputFileName != "" {
		innerInputFileName = inputFileName
	} else {
		innerInputFileName = DEFAULT_YAML_FILE
	}
	if checkIfFilePathExists(innerInputFileName) {
		return "", fmt.Errorf("file path does not exist: %s", innerInputFileName)
	}
	if err := checkIsRegularFile(innerInputFileName); err != nil {
		return "", err
	}
	return innerInputFileName, nil
}

func checkIsRegularFile(filePath string) error {
	// if given file path exists and is not a directory
	//if fileInfo, err := os.Stat(filePath); err == nil {
	//    return !fileInfo.IsDir()
	//}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("file path does not exist: %s", filePath)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("file path is a directory: %s", filePath)
	}
	// check if the file is a textfile, not a binary file
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("file path is not a regular file: %s", filePath)
	}
	return nil
}

func main() {
	cmd := &cobra.Command{
		Use:   "metatask",
		Short: "metatask",
		Long:  `metatask`,
	}
	generateSubCmd := &cobra.Command{
		Use:   "generate",
		Short: "generate a new project",
		Long:  `generate a new project`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			infile, _ := cmd.Flags().GetString("in-file")
			l := logrus.New()
			infile, err := checkInputFileOrDefaultYamlFileIsRegularFileAndGetFileName(infile)
			if err != nil {
				return err
			}
			// if infile is not provided, use the default
			g := pkg.NewGenerator(l, infile, dryRun)
			//makefile, _ := cmd.Flags().GetString("makefile")
			outputMakefiles, _ := cmd.Flags().GetStringSlice("output-makefile")
			for _, makefile := range outputMakefiles {
				if makefile != "" {
					g.AddAdapter(pkg.NewMakefileAdapter(
						l,
						makefile,
						dryRun,
						"",
						"",
					))
				}
			}
			//packageJson, _ := cmd.Flags().GetString("package-json")
			outputPackageJsons, _ := cmd.Flags().GetStringSlice("output-package-json")
			for _, packageJson := range outputPackageJsons {
				if packageJson != "" {
					g.AddAdapter(pkg.NewNpmAdapter(
						l,
						packageJson,
						dryRun,
					))
				}
			}

			err = g.Generate()
			if err != nil {
				return err
			}
			return nil
		},
	}
	generateSubCmd.Flags().StringP("in-file", "i", "", "input file")
	generateSubCmd.Flags().BoolP("dry-run", "d", false, "dry run")
	generateSubCmd.Flags().StringSliceP("output-makefile", "m", []string{}, "makefile")
	generateSubCmd.Flags().StringSliceP("output-package-json", "n", []string{}, "npm")
	// sync command

	versionSubCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  `Print the version number`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("v1.0.0")
			cmd.Println(runtime.Version())
		},
	}

	cmd.AddCommand(generateSubCmd)
	cmd.AddCommand(versionSubCmd)
	err := cmd.Execute()
	if err != nil {
		cmd.Println(err)
	}
}