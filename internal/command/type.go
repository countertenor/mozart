package command

import (
	"time"

	"github.com/prashantgupta24/mozart/internal/bash"
	"github.com/prashantgupta24/mozart/internal/config"
	"github.com/spf13/pflag"
)

//Instance is the main struct for command configs
type Instance struct {
	Config    *config.Instance
	Error     error
	Flags     *pflag.FlagSet
	StartTime time.Time
	ConfigDir string
	bash.Instance
}

//SetAutoYesEnabled is a setter for autoYesEnabled field
func (i *Instance) SetAutoYesEnabled(autoYesEnabled bool) *Instance {
	i.AutoYesEnabled = autoYesEnabled
	return i
}

//SetRunOnLocal is a setter for runOnLocal field
func (i *Instance) SetRunOnLocal(runOnLocal bool) *Instance {
	i.RunOnLocal = runOnLocal
	return i
}
