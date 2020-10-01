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

	"github.com/prashantgupta24/mozart/internal/command"
	"github.com/prashantgupta24/mozart/internal/ws"
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
	commandCenter := command.New(flags).ParseDefault()

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &commandCenter.Config)

	err := commandCenter.PreCheck().Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	commandCenter.CreateFromConfig()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commandCenter.Config)

}

//ExecuteDir installs a module
func ExecuteDir(w http.ResponseWriter, r *http.Request) {

	commandCenter := command.New(getFlags(r.URL.Query()))
	err := commandCenter.ParseAll().Error
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
		http.Error(w, "invalid module", http.StatusBadRequest)
		return
	}

	go func() {
		commandCenter.RunBashScripts()
	}()

	fmt.Fprint(w, "Success!")
}

//GetState gets the state of a module
func GetState(w http.ResponseWriter, r *http.Request) {

	var moduleRequest moduleRequest
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &moduleRequest)

	commandCenter := command.New(getFlags(r.URL.Query())).ReturnStateForDir(moduleRequest.getModuleName())
	stateMap := commandCenter.ReturnStateMap

	regexFile := regexp.MustCompile("([0-9]+-)")
	//regex101 - /r/MFyEZw/2
	regexDir := regexp.MustCompile("([a-zA-Z0-9/]+/[0-9-]*)")
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
		stepInstance.Module = regexDir.ReplaceAllString(dir, "")
		taskInstance := task{}
		taskList := []task{}

		mapDir := stateMap[dir]
		var sortedFileKeys []string
		for key := range mapDir {
			sortedFileKeys = append(sortedFileKeys, key)
		}
		sort.Strings(sortedFileKeys)

		for _, file := range sortedFileKeys {
			taskInstance.TaskName = regexFile.ReplaceAllString(strings.TrimSuffix(file, commandCenter.Config.Metadata.Extension), "")
			taskInstance.FileExecStatus = mapDir[file]
			if taskInstance.FileExecStatus.State == "running" { //update running time
				taskInstance.FileExecStatus.TimeTaken = time.Since(taskInstance.FileExecStatus.StartTime).String()
			}
			if taskInstance.FileExecStatus.State != "" {
				stateOfExecution = string(taskInstance.FileExecStatus.State)
			}
			taskList = append(taskList, taskInstance)
		}
		stepInstance.Tasks = taskList
		stepList = append(stepList, stepInstance)
	}
	stateJSON.Steps = stepList
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
