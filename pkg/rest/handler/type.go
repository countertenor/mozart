package handler

import "github.com/prashantgupta24/mozart/internal/bash"

type configurationRequest struct {
	Type  string `json:"type"`
	Nodes int    `json:"nodes"`
}

type task struct {
	TaskName       string              `json:"taskName"`
	FileExecStatus bash.FileExecStatus `json:"status"`
}

type step struct {
	Module string `json:"module"`
	Tasks  []task `json:"tasks"`
}

type stateJSON struct {
	State string `json:"state"`
	Steps []step `json:"steps"`
}
