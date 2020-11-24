package execution

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

var cancelRunningCommandFunc context.CancelFunc
var waitGroup sync.WaitGroup

//Init the execution instance
func (i *Instance) Init() {
	i.WaitGroup = &waitGroup
	i.DirExecStatusMap = makeStatusMap()
	i.initState()
	i.Interrupter = i.configureInterrupter()
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

	fileMetadata := fileMetadata{
		fullDirPath: fullDirPath,
		filename:    filename,
	}

	if i.DoRunParallel {
		defer i.WaitGroup.Done()
	}
	if i.DirExecStatusMap[fullDirPath][filename].State == RunningState {
		if time.Since(i.DirExecStatusMap[fullDirPath][filename].StartTime) > i.TimeoutInterval {
			i.updateTimeoutState(fileMetadata)
		} else {
			return fmt.Errorf("Already running")
		}
	}
	if !i.ReRun && i.DirExecStatusMap[fullDirPath][filename].State == SuccessState {
		fmt.Printf("Skipping %v since it ran successfully in last execution.\n", filename)
		return nil
	}
	if i.DryRunEnabled {
		i.updateNotStartedState(fileMetadata)
		return nil
	}
	if i.Error != nil {
		fmt.Printf("Skipping %v since previous file had errors or was terminated.\n", filename)
		i.updateSkipState(fileMetadata)
		return nil
	}
	i.PrintSeparator()
	fmt.Printf("\nRunning file : %v\n\n", fullDirPath+"/"+filename)
	args := []string{fullDirPath + "/" + filename}
	ctx, cancelFunc := context.WithTimeout(context.Background(), i.TimeoutInterval)
	defer cancelFunc()
	cancelRunningCommandFunc = cancelFunc

	command := exec.CommandContext(ctx, getSource(filename), args...)

	logFile, logfilePath, err := createLogFile(filename, i.LogDir)
	if err != nil {
		i.updateErrorState(fileMetadata)
		return err
	}
	defer logFile.Close()

	fileMetadata.logfilePath = logfilePath
	if i.DoRunParallel {
		command.Stdout = logFile
		command.Stderr = logFile
	} else {
		writeToStdOutputAndFile := io.MultiWriter(logFile, os.Stdout)
		command.Stdout = writeToStdOutputAndFile
		command.Stderr = writeToStdOutputAndFile
	}

	i.updateRunningStatus(fileMetadata)

	err = command.Run()
	if err != nil {
		errMessage := ""
		if ctx.Err() == context.DeadlineExceeded {
			i.updateTimeoutState(fileMetadata)
			errMessage = "script timeout"
		} else if ctx.Err() == context.Canceled {
			i.updateCancelState(fileMetadata)
			errMessage = "script canceled"
		} else {
			i.updateErrorState(fileMetadata)
		}
		signal.Stop(i.Interrupter)
		return fmt.Errorf("error while running script %v, %v, err: %v", filename, errMessage, err)
	}
	i.PrintSeparator()
	fmt.Printf("File ran successfully : %v\n", filename)
	i.updateSuccessState(fileMetadata)
	return nil
}

func (i *Instance) configureInterrupter() chan os.Signal {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case m := <-c:
			fmt.Println("canceling due to : ", m)
			i.StopRunningCmd()
			signal.Stop(c)
		}
	}()
	return c
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

func createLogFile(filename, logDir string) (*os.File, string, error) {

	//check if log dir exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			return nil, "", fmt.Errorf("error while creating log dir %v, err : %v", logDir, err)
		}
	}

	//create log file
	timeNow := time.Now().Format("2006-01-02--15-04-05.000")
	logfilePath := logDir + timeNow + "-" + filename[:strings.LastIndex(filename, ".")] + ".log"

	absPath, err := filepath.Abs(logfilePath)
	if err != nil {
		return nil, "", fmt.Errorf("error while getting absolute dir for logfile %v, err : %v", filename, err)
	}
	logFile, err := os.Create(logfilePath)
	if err != nil {
		return nil, "", fmt.Errorf("cannot create logfile for %v, err : %v", logfilePath, err)
	}
	return logFile, absPath, nil
}

func getSource(filename string) (source string) {
	fileExt := filename[strings.LastIndex(filename, "."):]
	switch fileExt {
	case ".sh":
		source = "/bin/bash"
	case ".py":
		source = "python"
	}
	return
}

//PrintSeparator prints a separator
func (i *Instance) PrintSeparator() {
	fmt.Println("=================================================================================================================")
}

func (i *Instance) updateRunningStatus(fileMetadata fileMetadata) {
	i.updateState(fileMetadata, RunningState)
}

func (i *Instance) updateSkipState(fileMetadata fileMetadata) {
	i.updateState(fileMetadata, skipped)
}

func (i *Instance) updateNotStartedState(fileMetadata fileMetadata) {
	i.updateState(fileMetadata, notStarted)
}

func (i *Instance) updateErrorState(fileMetadata fileMetadata) {
	i.updateState(fileMetadata, ErrorState)
}

func (i *Instance) updateTimeoutState(fileMetadata fileMetadata) {
	i.updateState(fileMetadata, timeout)
}

func (i *Instance) updateCancelState(fileMetadata fileMetadata) {
	i.updateState(fileMetadata, canceled)
}

func (i *Instance) updateSuccessState(fileMetadata fileMetadata) {
	i.updateState(fileMetadata, SuccessState)
}

func (i *Instance) updateState(fileMetadata fileMetadata, state stateType) {
	dirExecStatusMap := i.DirExecStatusMap //will update according to execution
	// fmt.Println("dirExecStatusMap : ", dirExecStatusMap)

	i.Mutex.Lock()
	defer i.Mutex.Unlock()

	var fileExecStatus FileExecStatus
	if fileExecStatusMap, exists := dirExecStatusMap[fileMetadata.fullDirPath]; exists {
		fileExecStatus = fileExecStatusMap[fileMetadata.filename]
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
	fileExecStatus.LogFilePath = fileMetadata.logfilePath

	if fileExecStatusMap, exists := dirExecStatusMap[fileMetadata.fullDirPath]; exists {
		fileExecStatusMap[fileMetadata.filename] = fileExecStatus
	} else {
		fileExecStatusMap := make(map[string]FileExecStatus)
		fileExecStatusMap[fileMetadata.filename] = fileExecStatus
		dirExecStatusMap[fileMetadata.fullDirPath] = fileExecStatusMap
	}

	i.saveState()
}
