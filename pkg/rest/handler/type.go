package handler

import (
	"strings"

	"github.com/prashantgupta24/mozart/internal/execution"
)

type moduleRequest struct {
	ModuleName string `json:"moduleName"`
}

func (m *moduleRequest) getModuleName() string {
	return strings.Join(strings.Split(m.ModuleName, " "), "/")
}

type task struct {
	TaskName       string                   `json:"taskName"`
	FileExecStatus execution.FileExecStatus `json:"status"`
}

type step struct {
	Directory string `json:"directory"`
	Module    string `json:"module"`
	Tasks     []task `json:"tasks"`
}

type stateJSON struct {
	State string `json:"state"`
	Steps []step `json:"steps"`
}
