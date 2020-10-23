package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/countertenor/mozart/internal/command"
	"github.com/countertenor/mozart/internal/execution"
	"github.com/countertenor/mozart/internal/ws"
)

//Index page of application
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

//WebSocket starts the websocket
func WebSocket(w http.ResponseWriter, r *http.Request) {
	ws.ServeWs(w, r)
}

//GetModules gets the possible modules to execute
func GetModules(w http.ResponseWriter, r *http.Request) {
	modules, err := command.GetAllDirsInsideTmpl()
	if err != nil {
		http.Error(w, "could not fetch modules, err : "+err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modules)
}

//Configuration endpoint
func Configuration(w http.ResponseWriter, r *http.Request) {

	flags := getFlags(r.URL.Query())
	commandCenter := command.New(flags)

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &commandCenter.Config)

	commandCenter.CreateFromConfig()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commandCenter.Config)

}

//ExecuteDir executes a module
func ExecuteDir(w http.ResponseWriter, r *http.Request) {

	commandCenter := command.New(getFlags(r.URL.Query()))
	err := commandCenter.ParseConfig().Error
	if err != nil {
		http.Error(w, "error with configuration, err: "+err.Error(), http.StatusBadRequest)
		return
	}

	var moduleRequest moduleRequest
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &moduleRequest)
	if moduleRequest.getModuleName() == "" {
		http.Error(w, "invalid module", http.StatusBadRequest)
		return
	}

	err = commandCenter.GenerateConfigFilesFromDir(moduleRequest.getModuleName()).Error
	if err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusBadRequest)
		return
	}

	errChan := make(chan error)
	go func(errChan chan error) {
		errChan <- commandCenter.RunScripts().Error
	}(errChan)

	select {
	case err := <-errChan:
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		break
	case <-time.After(time.Second * 2):
		break
	}
	fmt.Fprint(w, "Success!")
	return
}

//GetState gets the state of a module
func GetState(w http.ResponseWriter, r *http.Request) {

	var moduleRequest moduleRequest
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &moduleRequest)

	commandCenter := command.New(getFlags(r.URL.Query())).ReturnStateForDir(moduleRequest.getModuleName())
	stateMap := commandCenter.ReturnStateMap

	regexFile := regexp.MustCompile("([0-9]+-)")
	//regex101 - /r/MFyEZw/3
	regexModule := regexp.MustCompile("([a-zA-Z0-9/-]+/[0-9-]*)")
	//regex101 - /r/EZ9ZQ6/1
	regexDir := regexp.MustCompile("^([a-zA-Z]+/)")
	stateJSON := stateJSON{}
	stepInstance := step{}
	stepList := []step{}
	var stateOfExecution string

	var sortedDirkeys []string
	for key := range stateMap {
		sortedDirkeys = append(sortedDirkeys, key)
	}
	sort.Strings(sortedDirkeys)

	for _, dir := range sortedDirkeys {
		stepInstance.Directory = regexDir.ReplaceAllString(dir, "")
		stepInstance.Module = regexModule.ReplaceAllString(dir, "")
		taskInstance := task{}
		taskList := []task{}

		mapDir := stateMap[dir]
		var sortedFileKeys []string
		for key := range mapDir {
			sortedFileKeys = append(sortedFileKeys, key)
		}
		sort.Strings(sortedFileKeys)

		for _, file := range sortedFileKeys {
			taskInstance.TaskName = regexFile.ReplaceAllString(strings.TrimSuffix(file, commandCenter.ExecFileExtension), "")
			taskInstance.FileExecStatus = mapDir[file]
			if taskInstance.FileExecStatus.State == execution.RunningState { //update running time
				taskInstance.FileExecStatus.TimeTaken = time.Since(taskInstance.FileExecStatus.StartTime).String()
				stateOfExecution = string(execution.RunningState)
			}
			taskList = append(taskList, taskInstance)
		}
		stepInstance.Tasks = taskList
		stepList = append(stepList, stepInstance)
	}
	stateJSON.Steps = stepList
	if stateOfExecution == "" {
		stateOfExecution = "completed"
	}
	stateJSON.State = stateOfExecution

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stateJSON)
}

//Cancel cancels running command
func Cancel(w http.ResponseWriter, r *http.Request) {
	err := command.New(getFlags(r.URL.Query())).
		StopRunningCommand().
		Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
