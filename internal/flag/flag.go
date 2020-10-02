package flag

import (
	"github.com/spf13/pflag"
)

//Constant flags used in CLI
const (
	OS                = "os"
	ExecutionSource   = "source"
	ExecFileExtension = "ext"
	DoRunParallel     = "parallel"
	NoGenerate        = "no-generate"
	ReRun             = "re-run"
	DryRun            = "dry-run"
	Verbose           = "verbose"
	ConfigurationFile = "conf"
)

//Init initializes all flags
func Init(flags *pflag.FlagSet) {
	flags.BoolP(Verbose, "v", false, "print verbosely")
	flags.StringP(ConfigurationFile, "c", "mozart-sample.yaml", "configuration yaml file needed for application")
	flags.BoolP(DryRun, "d", false, "(optional) shows what scripts will run, but does not run the scripts")
	flags.BoolP(ReRun, "r", false, "(optional) re-run script from initial state, ignoring previously saved state")
	flags.BoolP(NoGenerate, "n", false, "(optional) do not generate scripts as part of execute, instead use the ones in generated folder. Useful for running local change to the scripts")
	flags.BoolP(DoRunParallel, "p", false, "(optional) Run all scripts in parallel")
	flags.StringP(OS, "o", "darwin", "(optional) OS on which scripts are allowed to run")
	flags.StringP(ExecutionSource, "s", "bash", "(optional) Execution source to use [Bash|Python|...]")
	flags.StringP(ExecFileExtension, "x", "sh", "(optional) Extension for execution files [sh|py|...]")
}
