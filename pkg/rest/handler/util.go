package handler

import (
	"github.com/prashantgupta24/mozart/internal/flag"
	"github.com/spf13/pflag"
)

func getFlags(queryParams map[string][]string) *pflag.FlagSet {
	flags := pflag.NewFlagSet("mozart-rest", pflag.ContinueOnError)

	var confFile string
	var dryRun bool
	var reRun bool
	var doRunParallel bool

	flags.StringVarP(&confFile, flag.ConfigurationFile, "c", "mozart-sample.yaml", "configuration yaml file needed for application")
	flags.BoolVarP(&dryRun, flag.DryRun, "d", false, "(optional) shows what scripts will run, but does not run the scripts")
	flags.BoolVarP(&reRun, flag.ReRun, "r", false, "(optional) re-run script from initial state, ignoring previously saved state")
	flags.BoolVarP(&doRunParallel, flag.DoRunParallel, "p", false, "(optional) Run all scripts in parallel")

	for key, value := range queryParams {
		flags.Set(key, value[0])
	}
	return flags
}
