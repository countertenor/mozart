package bash

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/prashantgupta24/mozart/internal/state"
	"github.com/prashantgupta24/mozart/statik"
)

//State is the main struct to hold state object
type State struct {
	StateFilePath        string
	StateFileDefaultname string
	DeploymentID         string
	StateFileName        string
	ReturnStateMap       DirExecStatusMap
	DirExecStatusMap     DirExecStatusMap
}

func (i *Instance) initState() {
	state.InitiateFileTypeInstance(i.StateFilePath + i.StateFileDefaultname)
}

//DeleteState deletes state
func (i *Instance) DeleteState(directory string) *Instance {
	if i.Error != nil {
		return i
	}
	statusMap := makeStatusMap()
	err := state.Load(&statusMap)
	if err != nil {
		i.Error = fmt.Errorf("unable to load previous state , err : %v", err)
		return i
	}

	if directory != "" {
		if !strings.Contains(directory, i.GeneratedDir) {
			actualDir, err := statik.GetActualDirName(directory, i.TemplateDir)
			if err != nil {
				i.Error = fmt.Errorf("unable to get actual dir for %v. err : %v", directory, err)
				return i
			}
			if actualDir == "" {
				fmt.Printf("State for directory %v not found.\n", directory)
				return i
			}

			directory = i.GeneratedDir + actualDir
		}
		// fmt.Println("dir : ", directory)
		if statusMap[directory] == nil {
			for key := range statusMap {
				if strings.Contains(key, directory) {
					// fmt.Println("deleting : ", key)
					delete(statusMap, key)
				}
			}
		} else {
			delete(statusMap, directory)
		}
	} else {
		//delete entire map
		statusMap = makeStatusMap()
	}
	i.DirExecStatusMap = statusMap
	i.saveState()

	if i.Error != nil {
		i.Error = fmt.Errorf("unable to delete state , err : %v", err)
		return i
	}
	return i
}

//ReturnState returns state
func (i *Instance) ReturnState(directory string) *Instance {
	if i.Error != nil {
		return i
	}
	statusMap, err := parseState(directory, i.TemplateDir, i.GeneratedDir)
	if err != nil {
		i.Error = err
		return i
	}
	i.ReturnStateMap = statusMap
	return i
}

//PrintState prints current state
func (i *Instance) PrintState(directory string) *Instance {
	if i.Error != nil {
		return i
	}
	statusMap, err := parseState(directory, i.TemplateDir, i.GeneratedDir)
	if err != nil {
		i.Error = err
		return i
	}
	jsonData, err := json.MarshalIndent(statusMap, "", "  ")
	if err != nil {
		i.Error = fmt.Errorf("error parsing state for printing, err : %v", err)
		return i
	}
	fmt.Printf("\nState: %s\n", string(jsonData))
	return i
}

func parseState(directory, templateDir, generatedDir string) (DirExecStatusMap, error) {

	stateMapToReturn := makeStatusMap()
	statusMap := makeStatusMap()
	err := state.Load(&statusMap)
	if err != nil {
		return nil, fmt.Errorf("unable to load previous state , err : %v", err)
	}

	// var jsonData []byte
	//Marshal
	if directory != "" {
		actualDir, err := statik.GetActualDirName(directory, templateDir)
		if err != nil {
			return nil, fmt.Errorf("unable to get actual dir for %v. err : %v", directory, err)
		}
		if actualDir == "" {
			fmt.Printf("State for directory %v not found.\n", directory)
			return stateMapToReturn, nil
		}

		directory = generatedDir + actualDir
		// fmt.Println("dir : ", directory)
		if statusMap[directory] == nil {
			var keys []string
			for key := range statusMap {
				if strings.Contains(key, directory) {
					keys = append(keys, key)
				}
			}
			// fmt.Println("keys : ", keys)
			tempMap := makeStatusMap()
			for _, key := range keys {
				tempMap[key] = statusMap[key]
			}
			// jsonData, err = json.MarshalIndent(tempMap, "", "  ")
			stateMapToReturn = tempMap
		} else {
			// jsonData, err = json.MarshalIndent(statusMap[directory], "", "  ")
			stateMapToReturn[directory] = statusMap[directory]
		}
	} else {
		// jsonData, err = json.MarshalIndent(statusMap, "", "  ")
		stateMapToReturn = statusMap
	}
	// if err != nil {
	// 	return nil, fmt.Errorf("unable to parse status , err : %v", err)
	// }
	return stateMapToReturn, nil
}

func (i *Instance) loadState() {
	statusMap := makeStatusMap()
	err := state.Load(&statusMap)
	if err != nil {
		i.Error = fmt.Errorf("unable to load previous state , err : %v", err)
		return
	}
	i.DirExecStatusMap = statusMap
}

func (i *Instance) saveState() {
	err := state.Save(i.DirExecStatusMap)
	if err != nil {
		i.Error = fmt.Errorf("unable to save state , err : %v", err)
	}
}
