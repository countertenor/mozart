package config

//Instance is the main config struct
type Instance struct {
	Values struct {
		Required struct {
			Value1 string `yaml:"value1" validate:"required" json:"value1"`
		} `yaml:"required" json:"required"`
		Optional struct {
			Value1 string `yaml:"value1" json:"value1"`
		} `yaml:"optional" json:"optional"`
	} `yaml:"values" json:"values"`
}
