package config

//Instance is the main config struct
type Instance struct {
	Value1 string `yaml:"value1" validate:"required" json:"value1"`
}
