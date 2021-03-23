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

	errChanMain := make(chan error)
	errChanPipe := make(chan error)

	go func() {
		err := commandCenter.RunScripts().Error
		errChanPipe <- err
	}()

	/*
		This go function is created so that the go function above
		can exit successfully. The above go function needs a
		valid, open channel to send its message to once it completes,
		otherwise it will be stuck in sending state forever (or panic if
		trying to send to closed channel).

		We cannot entrust the above go function to close the channel
		since it is busy executing the main function.

		This go function starts a timer, and makes sure that the above go function
		only gets to propagate its error if it returns before the timer
		finishes. If not, this go function takes the err and discards it, since
		the error is irrelevant at this stage.

		Both cases in this go function need to close the main error channel,
		since the response is waiting for that channel to close.
	*/
	go func() {
		timer := time.NewTimer(time.Second)
		timerDone := false
		for {
			select {
			case <-timer.C:
				// fmt.Println("time to close")
				close(errChanMain)
				timer.Stop()
				timerDone = true
			case err := <-errChanPipe:
				if !timerDone {
					errChanMain <- err
					close(errChanMain)
				}
				//else discard error
				return
			}
		}
	}()

	for err := range errChanMain {
		if err != nil {
			// fmt.Println("returned error : ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	// fmt.Println("done")

	fmt.Fprint(w, "Success!")
	return
}

//GetState gets the state of a module
func GetState(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("module : ", r.URL.Query().Get("module"))
	moduleName := strings.Join(strings.Split(r.URL.Query().Get("module"), " "), "/")
	commandCenter := command.New(getFlags(r.URL.Query())).
		ReturnStateForDir(moduleName)
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
			taskInstance.TaskName = regexFile.ReplaceAllString(file[:strings.LastIndex(file, ".")], "")
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
	if stateOfExecution == "" && len(stepList) > 0 {
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
