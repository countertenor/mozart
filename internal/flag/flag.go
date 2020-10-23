package flag

import (
	"github.com/spf13/pflag"
)

//Constant flags used in CLI
const (
	OS                = "os"
	DoRunParallel     = "parallel"
	NoGenerate        = "no-generate"
	ReRun             = "re-run"
	DryRun            = "dry-run"
	Verbose           = "verbose"
	ConfigurationFile = "conf"
)

//Init initializes all flags
func Init(flags *pflag.FlagSet) {
	AddVerboseFlag(flags)
	AddConfFileFlag(flags)
	AddDryRunFlag(flags)
	AddReRunFlag(flags)
	AddNoGenFlag(flags)
	AddRunParFlag(flags)
	AddOSFlag(flags)
}

//AddVerboseFlag flag
func AddVerboseFlag(flags *pflag.FlagSet) {
	flags.BoolP(Verbose, "v", false, "print verbosely")
}

//AddConfFileFlag flag
func AddConfFileFlag(flags *pflag.FlagSet) {
	flags.StringP(ConfigurationFile, "c", "mozart-sample.yaml", "configuration yaml file needed for application")
}

//AddDryRunFlag flag
func AddDryRunFlag(flags *pflag.FlagSet) {
	flags.BoolP(DryRun, "d", false, "(optional) shows what scripts will run, but does not run the scripts")
}

//AddReRunFlag flag
func AddReRunFlag(flags *pflag.FlagSet) {
	flags.BoolP(ReRun, "r", false, "(optional) re-run script from initial state, ignoring previously saved state")
}

//AddNoGenFlag flag
func AddNoGenFlag(flags *pflag.FlagSet) {
	flags.BoolP(NoGenerate, "n", false, "(optional) do not generate scripts as part of execute, instead use the ones in generated folder. Useful for running local change to the scripts")
}

//AddRunParFlag flag
func AddRunParFlag(flags *pflag.FlagSet) {
	flags.BoolP(DoRunParallel, "p", false, "(optional) Run all scripts in parallel")
}

//AddOSFlag flag
func AddOSFlag(flags *pflag.FlagSet) {
	flags.StringP(OS, "o", "darwin", "(optional) OS on which scripts are allowed to run")
}
