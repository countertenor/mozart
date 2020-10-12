package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/prashantgupta24/mozart/internal/execution"
	"github.com/prashantgupta24/mozart/internal/flag"
	"github.com/prashantgupta24/mozart/internal/template"
	"github.com/prashantgupta24/mozart/internal/yaml"
	"github.com/prashantgupta24/mozart/statik"

	"github.com/spf13/pflag"
)

var logDirPathFromEnv string  //This will be set through the build command, see Makefile
var stateDBPathFromEnv string //This will be set through the build command, see Makefile

//constants needed
const (
	sampleConfigFileName = "mozart-sample.yaml"

	defaultConfigFileName = "mozart-defaults.yaml"
	stateFileDefaultName  = "mozart-state.db"

	generatedDir    = "generated"
	templateDir     = "/templates"
	templateFileExt = ".tmpl"
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

	executionInstance := execution.Instance{
		LogDir:            logDir,
		GeneratedDir:      generatedDir,
		TemplateDir:       templateDir,
		OS:                getStringFlagValue(flags, flag.OS),
		ExecutionSource:   getStringFlagValue(flags, flag.ExecutionSource),
		ExecFileExtension: getStringFlagValue(flags, flag.ExecFileExtension),
		DoRunParallel:     getBoolFlagValue(flags, flag.DoRunParallel),
		DryRunEnabled:     getBoolFlagValue(flags, flag.DryRun),
		ReRun:             getBoolFlagValue(flags, flag.ReRun),
		TimeoutInterval:   time.Minute * 15, //change later
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
		i.Error = fmt.Errorf("error while parsing YAML file: %v", err)
		return i
	}
	if getBoolFlagValue(i.Flags, flag.Verbose) {
		fmt.Println("config : ", i.Config)
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
		configDir, err = statik.GetActualDirName(statik.Template, dirToGenerateFrom, templateDir)
		if err != nil {
			i.Error = fmt.Errorf("could not get ActualDirName for dir %v, err : %v ", dirToGenerateFrom, err)
			return i
		}
		if configDir == "" {
			i.Error = fmt.Errorf("could not find directory or directory is empty %v", dirToGenerateFrom)
			return i
		}
	}
	//fmt.Println("actual dir : ", configDir)
	i.ConfigDir = configDir
	if !noGenerate {
		//cleaning up all scripts in dir if it exists
		if _, err := os.Stat(generatedDir + configDir); !os.IsNotExist(err) {
			filesDeleted, err := cleanupFilesInDir(generatedDir+configDir, i.ExecFileExtension)
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
			templateFileExt,
			i.ExecFileExtension,
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
	fullPath := generatedDir + i.ConfigDir
	// fmt.Println("fullPath : ", fullPath)

	//skip execution the first time to populate state obj
	i.DryRunEnabled = true
	i.RunScriptsInDir(fullPath)
	i.DryRunEnabled = false
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
	dirs, err := statik.GetAllDirsInDir(statik.Template, templateDir)
	if err != nil {
		return nil, err
	}
	// fmt.Println("dirs : ", dirs)
	return dirs, nil
}

//returns number of files cleaned up, along with error (if nil)
func cleanupFilesInDir(directory, fileExt string) (int, error) {
	filesDeleted := 0
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == fileExt {
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
