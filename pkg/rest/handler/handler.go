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
	"github.com/prashantgupta24/mozart/internal/flag"
	"github.com/prashantgupta24/mozart/internal/ws"
	"github.com/spf13/pflag"
)

//Index page of application
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

//WebSocket starts the websocket
func WebSocket(w http.ResponseWriter, r *http.Request) {
	ws.ServeWs(w, r)
}

//Configuration endpoint
func Configuration(w http.ResponseWriter, r *http.Request) {

	flags := getFlags(r.URL.Query())
	commandCenter := command.New(flags)
	commandCenter.CreateSampleConfigFile().
		ParseYaml(command.SampleConfigFileName)

		//make changes
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newConfigurationRequest configurationRequest
	json.Unmarshal(reqBody, &newConfigurationRequest)

	config := commandCenter.Config
	//make changes

	err := config.Validate()
	if err != nil {
		http.Error(w, "error while validating values, err: "+err.Error(), http.StatusBadRequest)
		return
	}

	confFile, _ := flags.GetString(flag.ConfigurationFile)
	commandCenter.CreateFromConfig(confFile)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)

}

//ExecuteDir installs a module
func ExecuteDir(w http.ResponseWriter, r *http.Request) {

	commandCenter := command.New(getFlags(r.URL.Query()))
	err := commandCenter.ParseYaml(command.SampleConfigFileName).Error

	if err != nil {
		http.Error(w, "error with configuration, err: "+err.Error(), http.StatusBadRequest)
		return
	}

	go func() {
		commandCenter.
			GenerateConfigFilesFromDir("test"). //TODO replace with actual module
			RunBashScripts()
	}()

	fmt.Fprint(w, "Success!")
}

//GetState gets the state of a module
func GetState(w http.ResponseWriter, r *http.Request) {

	commandCenter := command.New(getFlags(r.URL.Query())).ReturnStateForDir("test")
	stateMap := commandCenter.ReturnStateMap

	regexFile := regexp.MustCompile("([0-9]+-)")
	regexDir := regexp.MustCompile("([a-z0-9/]+-)")
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
	flags := pflag.NewFlagSet("REST", pflag.ContinueOnError)
	err := command.New(flags).
		StopRunningCommand().
		Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
