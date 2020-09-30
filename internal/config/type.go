package config

//Instance is the main config struct
type Instance struct {
	Metadata struct {
		OS               string `yaml:"os" json:"os"`
		Execution        string `yaml:"execution" json:"execution" validate:"oneof=bash python"`
		ExecutionCommand string `json:"executionCommand"` //calculated dynamically, based on above
		Extension        string `json:"extension"`        //calculated dynamically, based on above
	} `yaml:"metadata" json:"metadata"`
	Values struct {
		Required struct {
			Value1 string `yaml:"value1" validate:"required" json:"value1"`
		} `yaml:"required" json:"required"`
	} `yaml:"values" json:"values"`
}
