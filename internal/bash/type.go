package bash

import (
	"sync"
	"time"

	"github.com/prashantgupta24/mozart/internal/config"
)

type stateType string

//States of execution
const (
	SuccessState stateType = "success"
	ErrorState   stateType = "error"
	RunningState stateType = "running"
	skipped      stateType = "skipped"
	timeout      stateType = "timeout"
	canceled     stateType = "canceled"
	notStarted   stateType = ""
)

//DirExecStatusMap is a map of:
// directory -> file exec status of each file within that directory
type DirExecStatusMap map[string]map[string]FileExecStatus

//FileExecStatus captures the status of a file
type FileExecStatus struct {
	StartTime       time.Time `json:"startTime"`
	TimeTaken       string    `json:"timeTaken"`
	LastSuccessTime string    `json:"lastSuccessTime"`
	LastErrorTime   string    `json:"lastErrorTime"`
	State           stateType `json:"state"`
	LogFilePath     string    `json:"logFilePath"`
}

//Instance is the main struct for bash configs
type Instance struct {
	DryRunEnabled   bool
	ReRun           bool
	LogDir          string
	GeneratedDir    string
	TemplateDir     string
	Error           error
	DoRunParallel   bool
	WaitGroup       *sync.WaitGroup
	Config          *config.Instance
	TimeoutInterval time.Duration
	State
}

func makeStatusMap() DirExecStatusMap {
	statusMap := make(DirExecStatusMap)
	return statusMap
}
