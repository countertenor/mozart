package command

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/countertenor/mozart/internal/execution"
	"github.com/countertenor/mozart/internal/flag"
	"github.com/countertenor/mozart/internal/template"
	"github.com/countertenor/mozart/internal/yaml"
	"github.com/countertenor/mozart/static"

	"github.com/spf13/pflag"
)

var logDirPathFromEnv string  //This will be set through the build command, see Makefile
var stateDBPathFromEnv string //This will be set through the build command, see Makefile

//constants needed
const (
	sampleConfigFileName = "mozart-sample.yaml"
	commonDirName        = "common"

	defaultConfigFileName = "mozart-defaults.yaml"
	stateFileDefaultName  = "mozart-state.db"

	generatedDir = "generated"
	templateDir  = "templates"
)

func init() {
	if _, err := os.Stat(generatedDir); os.IsNotExist(err) {
		// fmt.Printf("\n%v directory does not exist, creating ...\n\n", generatedDir)
		err := os.Mkdir(generatedDir, 0755)
		if err != nil {
			log.Fatalf("init error while creating %v directory, cannot proceed. err: %v", generatedDir, err)
		}
	}
}

//New creates a new instance for command execution
func New(flags *pflag.FlagSet) *Instance {
	stateFilePath := "./"
	if stateDBPathFromEnv != "" {
		stateDBPathFromEnv = parsePath(stateDBPathFromEnv)
		stateFilePath = stateDBPathFromEnv
	}

	logDir := "logs/"
	if logDirPathFromEnv != "" {
		logDirPathFromEnv = parsePath(logDirPathFromEnv)
		logDir = logDirPathFromEnv
	}

	executionInstance := &execution.Instance{
		LogDir:          logDir,
		GeneratedDir:    generatedDir,
		TemplateDir:     templateDir,
		DoRunParallel:   getBoolFlagValue(flags, flag.DoRunParallel),
		DryRunEnabled:   getBoolFlagValue(flags, flag.DryRun),
		ReRun:           getBoolFlagValue(flags, flag.ReRun),
		TimeoutInterval: time.Minute * 15, //change later
		State: execution.State{
			StateFilePath:        stateFilePath,
			StateFileDefaultname: stateFileDefaultName,
		},
	}
	executionInstance.Init()

	return &Instance{
		Config:    make(map[string]interface{}),
		Flags:     flags,
		Instance:  executionInstance,
		StartTime: time.Now(),
	}
}

//CreateSampleConfigFile creates sample config file
func (i *Instance) CreateSampleConfigFile() *Instance {
	if i.Error != nil {
		return i
	}
	err := yaml.CreateSampleConfigFile(sampleConfigFileName)
	if err != nil {
		i.Error = err
		return i
	}
	fmt.Println("\nGenerated sample file : ", sampleConfigFileName)
	return i
}

//ParseConfig parses the file passed in through flags
func (i *Instance) ParseConfig() *Instance {
	if i.Error != nil {
		return i
	}
	confFile := getStringFlagValue(i.Flags, flag.ConfigurationFile)
	err := yaml.ParseFile(i.Config, confFile)
	if err != nil {
		i.Error = fmt.Errorf("error while parsing %v YAML file: %v", confFile, err)
		return i
	}
	//optional values from config file
	err = i.parseConfigParams()
	if err != nil {
		i.Error = err
		return i
	}
	//common folder files
	err = yaml.ParseCommonFolder(i.Config, commonDirName)
	if err != nil {
		i.Error = fmt.Errorf("error while parsing %v common dir: %v", commonDirName, err)
		return i
	}
	if getBoolFlagValue(i.Flags, flag.Verbose) {
		i.printConfig()
	}
	return i
}

//CreateFromConfig writes config struct to file
func (i *Instance) CreateFromConfig() *Instance {
	if i.Error != nil {
		return i
	}
	confFile := getStringFlagValue(i.Flags, flag.ConfigurationFile)
	err := yaml.CreateFromConfig(i.Config, confFile)
	if err != nil {
		i.Error = err
		return i
	}
	return i
}

//GenerateConfigFilesFromDir generates all config files
func (i *Instance) GenerateConfigFilesFromDir(dirToGenerateFrom string) *Instance {
	noGenerate := getBoolFlagValue(i.Flags, flag.NoGenerate)
	if i.Error != nil {
		return i
	}
	var configDir string
	var err error
	if dirToGenerateFrom != "" {
		configDir, err = static.GetActualDirName(static.ResourceType, dirToGenerateFrom, templateDir)
		if err != nil {
			i.Error = fmt.Errorf("could not get ActualDirName for dir %v, err : %v ", dirToGenerateFrom, err)
			return i
		}
		if configDir == "" {
			i.Error = fmt.Errorf("could not find directory or directory is empty %v", dirToGenerateFrom)
			return i
		}
	}
	// fmt.Println("configDir : ", configDir)
	i.ConfigDir = configDir
	if !noGenerate {
		//cleaning up all scripts in dir if it exists
		if _, err := os.Stat(generatedDir + configDir); !os.IsNotExist(err) {
			filesDeleted, err := cleanupFilesInDir(generatedDir + configDir)
			if err != nil {
				i.Error = fmt.Errorf("could not delete files in %v directory, err: %v", generatedDir+configDir, err)
				return i
			}
			if getBoolFlagValue(i.Flags, flag.Verbose) {
				fmt.Printf("\n%v directory exists, cleaned up %v files inside...\n\n", generatedDir+configDir, filesDeleted)
			}
		}
		err := template.Generate(i.Config,
			configDir,
			templateDir,
			generatedDir)
		if err != nil {
			i.Error = fmt.Errorf("error while creating configuration : %v", err)
			return i
		}
	}
	i.PrintSeparator()
	return i
}

//RunScripts runs all scripts in a directory
func (i *Instance) RunScripts() *Instance {
	if i.Error != nil {
		return i
	}
	fullPath := filepath.Join(generatedDir, i.ConfigDir)
	// fmt.Println("fullPath : ", fullPath)
	i.RunScriptsInDir(fullPath)

	i.Error = i.Instance.Error
	i.PrintSeparator()
	return i
}

//TimeTaken prints time taken for execution
func (i *Instance) TimeTaken() *Instance {
	if i.Error != nil {
		return i
	}
	fmt.Println("Time taken : ", time.Since(i.StartTime))
	return i
}

//DeleteStateForDir deletes state for given dir
func (i *Instance) DeleteStateForDir(directory string) *Instance {
	if i.Error != nil {
		return i
	}
	i.DeleteState(directory)
	return i
}

//ReturnStateForDir prints state for given dir
func (i *Instance) ReturnStateForDir(directory string) *Instance {
	if i.Error != nil {
		return i
	}
	i.ReturnState(directory)
	return i
}

//PrintStateForDir prints state for given dir
func (i *Instance) PrintStateForDir(directory string) *Instance {
	if i.Error != nil {
		return i
	}
	i.PrintState(directory)
	return i
}

//StopRunningCommand stops currently running command
func (i *Instance) StopRunningCommand() *Instance {
	if i.Error != nil {
		return i
	}
	i.StopRunningCmd()
	i.Error = i.Instance.Error
	return i
}

//GetAllDirsInsideTmpl gets all directories inside template folder
func GetAllDirsInsideTmpl() ([]string, error) {
	dirs, err := static.GetAllDirsInDir(static.ResourceType, templateDir)
	if err != nil {
		return nil, err
	}
	// fmt.Println("dirs : ", dirs)
	return dirs, nil
}

//returns number of files cleaned up, along with error (if nil)
func cleanupFilesInDir(directory string) (int, error) {
	filesDeleted := 0
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err = os.Remove(path)
			if err != nil {
				return fmt.Errorf("could not delete file %v, err : %v", info.Name(), err)
			}
			filesDeleted++
		}
		return nil
	})
	return filesDeleted, err
}

func getBoolFlagValue(flags *pflag.FlagSet, flagname string) bool {
	if value, err := flags.GetBool(flagname); err == nil {
		return value
	}
	return false
}

func getStringFlagValue(flags *pflag.FlagSet, flagname string) string {
	if value, err := flags.GetString(flagname); err == nil {
		return value
	}
	return ""
}

func parsePath(path string) string {
	lastChar := path[len(path)-1:]

	if lastChar != "/" {
		path += "/"
	}
	return path
}

func (i *Instance) printConfig() error {
	jsonData, err := json.MarshalIndent(i.Config, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to parse version, err : %v", err)
	}
	fmt.Printf("\nConfig %s\n", jsonData)
	return nil
}

func (i *Instance) parseConfigParams() error {
	if i.Config["log_path"] != nil {
		logPath, parseOk := i.Config["log_path"].(string)
		if parseOk {
			i.LogDir = i.LogDir + parsePath(logPath)
		} else {
			return fmt.Errorf("could not parse log file path in config file")
		}
	}

	if i.Config["exec_source"] != nil {
		errorMsg := "could not parse exec_source in config file"
		source, parseOk := i.Config["exec_source"].(map[interface{}]interface{})
		if parseOk {
			// fmt.Println("source : ", source)
			for key, val := range source {
				// fmt.Printf("key %v val %v", key, val)
				if key != nil && val != nil {
					keyStr, parseKeyOk := key.(string)
					valStr, parseValOk := val.(string)
					if parseKeyOk && parseValOk {
						i.ExecutionSource[keyStr] = valStr
					} else {
						return fmt.Errorf(errorMsg)
					}
				}
			}
		} else {
			return fmt.Errorf(errorMsg)
		}
	}

	if i.Config["args"] != nil {
		errorMsg := "could not parse args in config file"
		argumentMapInterface, parseOk := i.Config["args"].(map[interface{}]interface{})
		if !parseOk {
			return fmt.Errorf(errorMsg)
		}
		for key, val := range argumentMapInterface {
			argArr := []string{}
			keyStr, parseKey := key.(string)
			valInterfaceArr, parseVal := val.([]interface{})
			if parseKey && parseVal {
				for _, s := range valInterfaceArr {
					argStr, parseArgStr := s.(string)
					if parseArgStr {
						argArr = append(argArr, argStr)
					} else {
						return fmt.Errorf(errorMsg)
					}
				}
			} else {
				return fmt.Errorf(errorMsg)
			}
			i.ArgumentMap[keyStr] = argArr
		}
		// fmt.Println("argumentMap : ", argumentMap)
	}
	return nil
}
