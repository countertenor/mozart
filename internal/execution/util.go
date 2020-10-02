package execution

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

type status struct {
	fullDirPath string
	fileName    string
	logFileName string
}

var cancelRunningCommandFunc context.CancelFunc
var waitGroup sync.WaitGroup

type saveState func() error

//Init the execution instance
func (i *Instance) Init() {
	switch i.ExecutionSource {
	case "bash":
		i.ExecutionSource = "/bin/bash"
		i.ExecFileExtension = ".sh"
	case "python":
		i.ExecutionSource = "python"
		i.ExecFileExtension = ".py"
	default:
		i.ExecFileExtension = "." + i.ExecFileExtension
	}
	i.WaitGroup = &waitGroup
	i.DirExecStatusMap = makeStatusMap()
	i.initState()
}

//RunScriptsInDir handles running of script files inside a directory
func (i *Instance) RunScriptsInDir(fullDirPath string) {
	i.loadState()
	if i.Error != nil {
		return
	}
	i.handleRunScripts(fullDirPath)
	if i.DoRunParallel {
		i.WaitGroup.Wait()
	}
}

//recursively run all scripts inside the directory
func (i *Instance) handleRunScripts(fullDirPath string) {
	filesInDir, err := ioutil.ReadDir(fullDirPath)
	if err != nil {
		i.Error = fmt.Errorf("error while reading files from dir %v, err : %v", fullDirPath, err)
		return
	}
	for _, file := range filesInDir {
		if file.IsDir() { //recursion into dir
			i.handleRunScripts(fullDirPath + "/" + file.Name())
		} else {
			if i.DoRunParallel {
				i.WaitGroup.Add(1)
				go func(fullDirPath, filename string) {
					i.runScript(fullDirPath, filename)
				}(fullDirPath, file.Name())
			} else {
				err := i.runScript(fullDirPath, file.Name())
				if err != nil {
					i.Error = err
				}
			}
		}
	}
}

func (i *Instance) runScript(fullDirPath, filename string) error {
	if i.DoRunParallel {
		defer i.WaitGroup.Done()
	}
	if i.SkipRun {
		return nil
	}
	i.PrintSeparator()
	fmt.Printf("\nRunning file : %v\n\n", fullDirPath+"/"+filename)
	//precautionary step so that scripts don't run locally
	osRunning := runtime.GOOS
	if osRunning != i.OS { //scripts run only on OS defined
		fmt.Printf("(Skipping execution since OS is %v. Scripts only run on %v)\n", osRunning, i.OS)
		return nil
	}
	if i.DirExecStatusMap[fullDirPath][filename].State == RunningState {
		return fmt.Errorf("Already running")
	}
	if !i.ReRun && i.DirExecStatusMap[fullDirPath][filename].State == SuccessState {
		fmt.Println("Skipping file since it ran successfully in last execution.")
		return nil
	}
	if i.DryRunEnabled {
		fmt.Printf("(Skipping execution because dry-run was selected)\n")
		i.updateNotStartedState(fullDirPath, filename, "")
		return nil
	}
	if i.Error != nil {
		fmt.Println("Skipping file since previous file had errors or was terminated.")
		i.updateSkipState(fullDirPath, filename, "")
		return nil
	}

	args := []string{fullDirPath + "/" + filename}
	ctx, cancelFunc := context.WithTimeout(context.Background(), i.TimeoutInterval)
	defer cancelFunc()
	cancelRunningCommandFunc = cancelFunc

	command := exec.CommandContext(ctx, i.ExecutionSource, args...)

	logFile, err := createLogFile(filename, i.LogDir, i.ExecFileExtension)
	if err != nil {
		i.updateErrorState(fullDirPath, filename, logFile.Name())
		return err
	}
	defer logFile.Close()

	if i.DoRunParallel {
		command.Stdout = logFile
		command.Stderr = logFile
	} else {
		writeToStdOutputAndFile := io.MultiWriter(logFile, os.Stdout)
		command.Stdout = writeToStdOutputAndFile
		command.Stderr = writeToStdOutputAndFile
	}

	i.updateRunningStatus(fullDirPath, filename, logFile.Name())
	err = command.Run()
	if err != nil {
		errMessage := ""
		if ctx.Err() == context.DeadlineExceeded {
			i.updateTimeoutState(fullDirPath, filename, logFile.Name())
			errMessage = "script timeout"
		} else if ctx.Err() == context.Canceled {
			i.updateCancelState(fullDirPath, filename, logFile.Name())
			errMessage = "script canceled"
		} else {
			i.updateErrorState(fullDirPath, filename, logFile.Name())
		}
		return fmt.Errorf("error while running script %v, %v, err: %v", filename, errMessage, err)
	}
	i.PrintSeparator()
	fmt.Printf("File ran successfully : %v\n", filename)
	i.updateSuccessState(fullDirPath, filename, logFile.Name())
	return nil
}

//StopRunningCmd stops currently running command
func (i *Instance) StopRunningCmd() {
	if cancelRunningCommandFunc != nil {
		cancelRunningCommandFunc()
	} else {
		i.Error = fmt.Errorf("nothing running at the moment")
	}
	return
}

func createLogFile(filename, logDir, fileExt string) (*os.File, error) {

	//check if log dir exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("error while creating log dir %v, err : %v", logDir, err)
		}
	}

	//create log file
	timeNow := time.Now().Format("2006-01-02--15-04-05.000")
	logFile, err := os.Create(logDir + timeNow + "-" + strings.TrimSuffix(filename, fileExt) + ".log")
	if err != nil {
		return nil, fmt.Errorf("cannot write to logfile, err : %v", err)
	}
	return logFile, nil
}

//PrintSeparator prints a separator
func (i *Instance) PrintSeparator() {
	fmt.Println("=================================================================================================================")
}

func (i *Instance) updateRunningStatus(directory, filename, logFilePath string) {
	i.updateState(directory, filename, logFilePath, RunningState)
}

func (i *Instance) updateSkipState(directory, filename, logFilePath string) {
	i.updateState(directory, filename, logFilePath, skipped)
}

func (i *Instance) updateNotStartedState(directory, filename, logFilePath string) {
	i.updateState(directory, filename, logFilePath, notStarted)
}

func (i *Instance) updateErrorState(directory, filename, logFilePath string) {
	i.updateState(directory, filename, logFilePath, ErrorState)
}

func (i *Instance) updateTimeoutState(directory, filename, logFilePath string) {
	i.updateState(directory, filename, logFilePath, timeout)
}

func (i *Instance) updateCancelState(directory, filename, logFilePath string) {
	i.updateState(directory, filename, logFilePath, canceled)
}

func (i *Instance) updateSuccessState(directory, filename, logFilePath string) {
	i.updateState(directory, filename, logFilePath, SuccessState)
}

func (i *Instance) updateState(directory, filename, logFilePath string, state stateType) {
	dirExecStatusMap := i.DirExecStatusMap //will update according to execution
	// fmt.Println("dirExecStatusMap : ", dirExecStatusMap)

	var fileExecStatus FileExecStatus
	if fileExecStatusMap, exists := dirExecStatusMap[directory]; exists {
		fileExecStatus = fileExecStatusMap[filename]
	}
	if state == SuccessState {
		fileExecStatus.LastSuccessTime = time.Now().String()
		fileExecStatus.TimeTaken = time.Since(fileExecStatus.StartTime).String()
	} else if state == ErrorState || state == canceled || state == timeout {
		fileExecStatus.LastErrorTime = time.Now().String()
		fileExecStatus.TimeTaken = time.Since(fileExecStatus.StartTime).String()
	} else if state == RunningState {
		fileExecStatus.StartTime = time.Now()
	}
	fileExecStatus.State = state
	fileExecStatus.LogFilePath = logFilePath

	if fileExecStatusMap, exists := dirExecStatusMap[directory]; exists {
		fileExecStatusMap[filename] = fileExecStatus
	} else {
		fileExecStatusMap := make(map[string]FileExecStatus)
		fileExecStatusMap[filename] = fileExecStatus
		dirExecStatusMap[directory] = fileExecStatusMap
	}

	i.saveState()
}
