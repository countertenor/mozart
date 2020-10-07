package command

import (
	"time"

	"github.com/prashantgupta24/mozart/internal/execution"
	"github.com/spf13/pflag"
)

//Instance is the main struct for command configs
type Instance struct {
	Config    map[string]interface{}
	Error     error
	Flags     *pflag.FlagSet
	StartTime time.Time
	ConfigDir string
	execution.Instance
}
