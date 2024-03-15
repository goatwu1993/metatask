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
			makefile, _ := cmd.Flags().GetString("makefile")
			if makefile != "" {
				g.AddAdapter(pkg.NewMakefileAdapter(l, makefile, "", "", dryRun))
			}
			npm, _ := cmd.Flags().GetString("package-json")
			if npm != "" {
				g.AddAdapter(pkg.NewNpmAdapter(
					l,
					npm,
					dryRun,
				))
			}

			err := g.Generate()
			if err != nil {
				return err
			}
			return nil
		},
	}
	// add a dry run flag to the generate command
	generateSubCmd.Flags().BoolP("dry-run", "d", false, "dry run")
	generateSubCmd.Flags().StringP("makefile", "m", "", "makefile")
	generateSubCmd.Flags().StringP("package-json", "n", "", "npm")

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