package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"metatask/pkg"

	"github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
)

const (
	DEFAULT_YAML_FILE = "metatask.yml"
)

var debug bool

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
	//if fileDebug, err := os.Stat(filePath); err == nil {
	//    return !fileDebug.IsDir()
	//}
	fileDebug, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("file path does not exist: %s", filePath)
	}
	if fileDebug.IsDir() {
		return fmt.Errorf("file path is a directory: %s", filePath)
	}
	// check if the file is a textfile, not a binary file
	if !fileDebug.Mode().IsRegular() {
		return fmt.Errorf("file path is not a regular file: %s", filePath)
	}
	return nil
}

func main() {
	l := logrus.New()
	var verbosity string

	rootCmd := &cobra.Command{
		Use:   "myapp",
		Short: "My application",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Adjust logging level based on verbosity flag
			switch strings.ToLower(verbosity) {
			case "":
				l.SetLevel(logrus.WarnLevel)
			case "v":
				l.SetLevel(logrus.InfoLevel)
			case "vv":
				l.SetLevel(logrus.DebugLevel)
			case "vvv":
				l.SetLevel(logrus.TraceLevel)
			default:
				l.SetLevel(logrus.WarnLevel)
			}
			fmt.Println("verbosity: ", verbosity)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", "", "verbosity level (v, vv, vvv)")
	generateSubCmd := &cobra.Command{
		Use:   "generate",
		Short: "generate a new project",
		Long:  `generate a new project`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			infile, _ := cmd.Flags().GetString("in-file")
			infile, err := checkInputFileOrDefaultYamlFileIsRegularFileAndGetFileName(infile)
			if err != nil {
				return err
			}
			// if infile is not provided, use the default
			g := pkg.NewGenerator(l, infile, dryRun)
			//makefile, _ := cmd.Flags().GetString("makefile")
			err = g.Generate()
			if err != nil {
				return err
			}
			return nil
		},
	}
	generateSubCmd.Flags().StringP("in-file", "i", "", "input file")
	generateSubCmd.Flags().BoolP("dry-run", "d", false, "dry run")
	// debug
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

	rootCmd.AddCommand(generateSubCmd)
	rootCmd.AddCommand(versionSubCmd)
	// debug
	rootCmd.Flags().BoolP("debug", "D", false, "debug")

	err := rootCmd.Execute()
	if err != nil {
		rootCmd.Println(err)
	}
}