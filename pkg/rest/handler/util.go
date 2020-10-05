package handler

import (
	"github.com/prashantgupta24/mozart/internal/flag"
	"github.com/spf13/pflag"
)

func getFlags(queryParams map[string][]string) *pflag.FlagSet {
	flags := pflag.NewFlagSet("mozart-rest", pflag.ContinueOnError)
	flag.Init(flags)
	for key, value := range queryParams {
		flags.Set(key, value[0])
	}
	return flags
}
