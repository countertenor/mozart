package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const longDescription = `

To load completions:

Bash:

$ source <(mozart completion bash)

# To load completions for each session, execute once:
Linux:
  $ mozart completion bash > /etc/bash_completion.d/mozart
MacOS:
  $ mozart completion bash > /usr/local/etc/bash_completion.d/mozart

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ mozart completion zsh > "${fpath[1]}/_mozart"

# You will need to start a new shell for this setup to take effect.

Fish:

$ mozart completion fish | source

# To load completions for each session, execute once:
$ mozart completion fish > ~/.config/fish/completions/mozart.fish
`

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "Generate completion script",
	Long:                  longDescription,
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			fmt.Println(longDescription)
			return nil
		}
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		default:
			fmt.Printf("\nerror: not a valid command")
			fmt.Println(longDescription)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
