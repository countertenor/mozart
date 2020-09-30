package cmd

import (
	"strings"

	"github.com/prashantgupta24/mozart/internal/command"
	"github.com/spf13/cobra"
)

var cleanState bool

// stateCmd represents the stateCmd command
var stateCmd = &cobra.Command{
	Use:   "state [install|cleanup]",
	Short: "See state",
	Long: `See state. For example:

	mozart state [install|cleanup] [linbit|pacemaker]
	
It is implemented in such a way that you can query directories within directories.

	For example:

	mozart state <parent_dir> <child_dir>
	`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var err error
		state := strings.Join(args, "/")

		commandCenter := command.New(cmd.Flags())
		if cleanState {
			err = commandCenter.
				DeleteStateForDir(state).
				Error
		} else {
			err = commandCenter.
				PrintStateForDir(state).
				Error
		}
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	stateCmd.Flags().BoolVar(&cleanState, "clean", false, "Clean state for application")
}
