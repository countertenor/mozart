package cmd

import (
	"github.com/prashantgupta24/mozart/internal/command"
	"github.com/spf13/cobra"
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse and validate configuration from yaml file",
	Long: `
This command parses and validates configuration from yaml file. This is to check if the conf file is valid.
You can also print the configuration parsed using the '-v' flag.

For example:

	mozart parse -c mozart-sample.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := command.New(cmd.Flags()).
			ParseAll(confFile).
			TimeTaken().
			Error
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// parseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
