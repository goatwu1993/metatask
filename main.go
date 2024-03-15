package main

import (
	"runtime"

	pkg "metatask/pkg"

	"github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
)

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
			l := logrus.New()
			g := pkg.NewGenerator(l, "metataskfile.json", dryRun)

			err := g.Generate()
			if err != nil {
				return err
			}
			return nil
		},
	}
	// add a dry run flag to the generate command
	generateSubCmd.Flags().BoolP("dry-run", "d", false, "dry run")

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