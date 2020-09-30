package config

//Instance is the main config struct
type Instance struct {
	Metadata struct {
		OS string `yaml:"os" json:"os"`
	} `yaml:"metadata" json:"metadata"`
	Values struct {
		Required struct {
			Value1 string `yaml:"value1" validate:"required" json:"value1"`
		} `yaml:"required"`
	} `yaml:"values"`
}
