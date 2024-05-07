package cmd

import "github.com/spf13/cobra"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the application",
	Long:  `All software has versions. This is the version number of the application`,

	Run: VersionCmd,
}

func VersionCmd(cmd *cobra.Command, args []string) {
	cmd.Printf("GeeLoadBalance version %s\n", "v0.0.1")
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
