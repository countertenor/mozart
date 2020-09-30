package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("test", test)
}

func test(fl validator.FieldLevel) bool {
	value1 := fl.Parent().FieldByName("value1")
	//fmt.Println("value1 : ", value1)
	num := fl.Field().Len()
	//fmt.Println("num : ", num)
	if int64(num) == value1.Int() {
		return true
	}
	return false
}

//Validate the struct
func (config *Instance) Validate() error {
	err := validate.Struct(config)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			// fmt.Println("err.Nam", err.Namespace())
			// fmt.Println("err.Fie", err.Field())
			// fmt.Println("err.Str", err.StructNamespace())
			// fmt.Println("err.Str", err.StructField())
			// fmt.Println("err.Tag", err.Tag())
			// fmt.Println("err.Act", err.ActualTag())
			// fmt.Println("err.Kin", err.Kind())
			// fmt.Println("err.Typ", err.Type())
			// fmt.Println("err.Val", err.Value())
			// fmt.Println("err.Par", err.Param())
			// fmt.Println()
			switch err.Tag() {
			case "test":
				return fmt.Errorf("err: value not correct")
			case "oneof":
				return fmt.Errorf("%v - possible values: [%v]", err, strings.Join(strings.Split(err.Param(), " "), ","))
			case "min":
				return fmt.Errorf("%v - possible values: [%v]", err, strings.Join(strings.Split(err.Param(), " "), ","))
			}
		}
		return err
	}
	return nil
}

//Print the struct
func (config *Instance) Print() error {
	return printPretty(config)
}

func printPretty(c interface{}) error {
	//Marshal
	json, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error printing config : %v", err)
	}
	fmt.Printf("\nParsed configuration: %s\n", string(json))
	return nil
}

//PreconfigureFields preconfigures some fields for config struct
func (config *Instance) PreconfigureFields() {

	var executionCommand string
	var extension string

	switch config.Metadata.Execution {
	case "bash":
		executionCommand = "/bin/bash"
		extension = ".sh"
	case "python":
		executionCommand = "python"
		extension = ".py"
	}

	config.Metadata.ExecutionCommand = executionCommand
	config.Metadata.Extension = extension
}
