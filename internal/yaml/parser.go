package yaml

import (
	"fmt"
	"io/ioutil"
	"os"

	"sigs.k8s.io/yaml"
)

//ParseFile parses file into config
func ParseFile(config map[string]interface{}, filename string) error {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no yaml file called %v found. Run 'init' to generate the sample config yaml file", filename)
		}
		return fmt.Errorf("error while reading conf file %v : %v", filename, err)
	}

	err = yaml.Unmarshal(fileData, &config)
	if err != nil {
		// return errors.WithStackTrace(err)
		return fmt.Errorf("error while unmarshalling yaml %v: %v", filename, err)
	}

	// NOTE: the client-go library can only decode json, so we will first convert the yaml to json before unmarshaling
	// jsonData, err := yaml.YAMLToJSON([]byte(fileData))
	// if err != nil {
	// 	// return errors.WithStackTrace(err)
	// 	fmt.Println("err1 : ", err)
	// }
	// fmt.Println("JSON : ", string(jsonData))

	// // var destinationObj interface{}
	// // destinationObj := make(map[string]interface{})
	// err = json.Unmarshal(jsonData, &config)
	// if err != nil {
	// 	// return errors.WithStackTrace(err)
	// 	fmt.Println("err2 : ", err)
	// }
	// // return nil

	// fmt.Println("destinationObj : ", config)
	// // err = yaml.Unmarshal([]byte(fileData), &config)
	// // if err != nil {
	// // 	return fmt.Errorf("error while unmarshalling yaml %v: %v", filename, err)
	// // }
	return nil
}
