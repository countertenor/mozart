package cmd

import (
	"fmt"
	"strings"

	"github.com/countertenor/mozart/internal/command"
	"github.com/countertenor/mozart/internal/flag"
	"github.com/spf13/cobra"
)

// execute represents the executeCmd command
var execute = &cobra.Command{
	Use:   "execute folder-name [folder-name ...]",
	Short: "Execute scripts inside any folder",
	Long: `Execute scripts inside any folder. For example:

	mozart execute folder-name [folder-name ...]

	After you setup auto-completion, pressing [tab][tab] will show all possible options.`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return getCommandsToRun(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("")
			fmt.Println("*****************************************************************")
			fmt.Println("Available commands:")
			fmt.Println("")
			commands := getCommandsToRun("")
			for _, command := range commands {
				fmt.Println("mozart execute", command)
			}
			fmt.Println("*****************************************************************")
			fmt.Println("")
			return fmt.Errorf("Execute needs at least one argument")
		}
		cmdToExecute := strings.Join(args, "/")
		err := command.New(cmd.Flags()).
			ParseConfig().
			GenerateConfigFilesFromDir(cmdToExecute).
			RunScripts().
			TimeTaken().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func getCommandsToRun(complete string) []string {
	var commands []string
	commands, _ = command.GetAllDirsInsideTmpl()
	return commands
}

func init() {
	rootCmd.AddCommand(execute)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execute.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execute.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	flag.Init(execute.Flags())
}
