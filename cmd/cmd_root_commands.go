package cmd

import "github.com/spf13/cobra"

func assembleRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "go-arch-lint",
		Short: "Golang architecture linter",
		Long: `
Check all project imports and compare to arch rules defined in yaml file.
Read full documentation in: https://github.com/fe3dback/go-arch-lint`,
		Run:                func(cmd *cobra.Command, args []string) {},
		DisableSuggestions: true,
		SilenceErrors:      true,
	}

	// version
	root.AddCommand(&cobra.Command{
		Use:   cmdNameVersion,
		Short: "Print go arch linter version",
		Run:   cmdVersion,
	})

	// check
	root.AddCommand(&cobra.Command{
		Use:   cmdNameCheck,
		Short: "check project architecture by yaml file",
		Run:   cmdCheck,
	})

	// assemble
	return root
}
