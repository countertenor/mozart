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

	flags.StringVarP(&confFile, flag.ConfigurationFile, "c", "mozart-sample.yaml", "configuration yaml file needed for application")
	flags.BoolVarP(&dryRun, flag.DryRun, "d", false, "(optional) shows what scripts will run, but does not run the scripts")
	flags.BoolVarP(&reRun, flag.ReRun, "r", false, "(optional) re-run script from initial state, ignoring previously saved state")

	for key, value := range queryParams {
		switch key {
		case "conf":
			flags.Set(flag.ConfigurationFile, value[0])
		case "dryRun":
			flags.Set(flag.DryRun, value[0])
		case "reRun":
			flags.Set(flag.ReRun, value[0])
		}
	}
	return flags
}
