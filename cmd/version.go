package cmd

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var versionInJSON bool

var gitCommitHash string
var buildTime string
var gitBranch string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long: `Print version. For example:

	mozart version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		versionMap := make(map[string]string)
		versionMap["Git commit hash"] = gitCommitHash
		versionMap["Build time"] = buildTime
		versionMap["Git branch"] = gitBranch

		if versionInJSON {
			jsonData, err := json.MarshalIndent(versionMap, "", "  ")
			if err != nil {
				return fmt.Errorf("unable to parse version, err : %v", err)
			}
			fmt.Printf("\nVersion information: %s\n", jsonData)
		} else {
			var keysInMap []string
			for key := range versionMap {
				keysInMap = append(keysInMap, key)
			}
			sort.Strings(keysInMap)
			fmt.Println("")
			for _, key := range keysInMap {
				fmt.Printf("%-20v : %v\n", key, versionMap[key])
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	versionCmd.Flags().BoolVar(&versionInJSON, "json", false, "Get JSON output")
}
