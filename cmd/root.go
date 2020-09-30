package cmd

import (
	"fmt"
	"os"

	"github.com/prashantgupta24/mozart/internal/flag"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool
var confFile string

var dryRun bool
var rerun bool
var noGenerate bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mozart",
	Short: "This is an orchestrator application for executing scripts.",
	Long:  `This is an orchestrator application for executing scripts.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mozart.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, flag.Verbose, "v", false, "print verbosely")
	rootCmd.PersistentFlags().StringVarP(&confFile, flag.ConfigurationFile, "c", "mozart-sample.yaml", "configuration yaml file needed for application")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, flag.DryRun, "d", false, "(optional) shows what scripts will run, but does not run the scripts")
	rootCmd.PersistentFlags().BoolVarP(&rerun, flag.ReRun, "r", false, "(optional) re-run script from initial state, ignoring previously saved state")
	rootCmd.PersistentFlags().BoolVarP(&noGenerate, flag.NoGenerate, "n", false, "(optional) do not generate bash scripts as part of install, instead use the ones in generated folder. Useful for running local change to the scripts")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".mozart" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mozart")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
