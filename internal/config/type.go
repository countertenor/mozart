package config

//Instance is the main config struct
type Instance struct {
	Values struct {
		Required struct {
			Value1   string   `yaml:"value1" validate:"required" json:"value1"`
			OS       string   `yaml:"os" validate:"required" json:"os"`
			Machines []string `yaml:"machines" validate:"required" json:"machines"`
		} `yaml:"required" json:"required"`
		Optional struct {
			Value1 string `yaml:"value1" json:"value1"`
		} `yaml:"optional" json:"optional"`
	} `yaml:"values" json:"values"`
}
