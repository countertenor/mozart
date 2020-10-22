import React, { useState, useEffect } from "react";
import { useHistory } from "react-router-dom";
import {
  FormLabel,
  TextInput,
  Tooltip,
  Dropdown,
  Form,
  Button,
  RadioButton,
  RadioButtonGroup,
  FormGroup,
  TextArea,
  FileUploader,
  Checkbox,
} from "carbon-components-react";
import axios from "axios";

export default function Configuration() {
  let history = useHistory();

  const allSourceFileTypes = ["Bash", "Python"];
  const allOS = ["Darwin", "Linux"];
  const allTypesOfModules = ["existing", "new"];

  let [configFileName, setConfigFileName] = useState("mozart-test.yaml")

  let [jsonObject, setJsonObject] = useState("")
  let [jsonFile, setJsonFile] = useState("")

  let [typeOfModule, setTypeOfModule] = useState(allTypesOfModules[0])

  let [modules, setModules] = useState([""]);
  let [selectedModule, setSelectedModule] = useState(modules.length > 0 ? [modules[0]] : [""]);
  const [newModuleName, setNewModuleName] = useState("");

  let [dryRun, setDryRun] = useState(false);
  let [reRun, setReRun] = useState(false);
  let [parallel, setParallel] = useState(false);

  let [source, setSourceFileTypes] = useState(allSourceFileTypes[0])
  let [os, setOS] = useState(allOS[0])

  const someProps = {
    invalid: true,
    invalidText: "This value cannot be empty. You must enter a valid json object here.",
  };

  let [validateTextArea, setValidateTextArea] = useState(false);

  //   const updateFieldChanged = e => {
  //     console.log("hey: ",items);
  //     let newArr = [...items]; // copying the old datas array
  //     newArr[newArr.length] = e.selectedItem; // replace e.target.value with whatever you want to change it to
  //     setRegion(newArr); // ??
  //     console.log("hello", items);
  // }
  function IsJsonString(str) {
    try {
        JSON.parse(str);
    } catch (e) {
        return false;
    }
    return true;
}

  const makeSampleAPICall = (e) => {
    e.preventDefault();
    console.log(!jsonFile)
    console.log(Object.keys(jsonObject).length)
    console.log(typeof jsonObject)
    if (
      IsJsonString(jsonObject) === false ||
      (Object.keys(jsonObject).length === 0 && !jsonFile)
    ) {
      console.log("ERROR!");
      setValidateTextArea(true);
    } else {
      setValidateTextArea(false);
      jsonObject = JSON.parse(jsonObject || "{}");
      let moduleName =
        selectedModule.length > 0 ? selectedModule : newModuleName;

      const dataBodyObj = {
        moduleName: moduleName,
        os: os,
      };
      const queryParamsObj = {
        "re-run": reRun,
        // "dry-run": dryRun,
        parallel: parallel,
        source: source.toLowerCase(),
      };

      console.log(dataBodyObj);
      console.log(queryParamsObj);

      let data = {};
      if (Object.keys(jsonObject).length === 0) {
        data = jsonFile;
      } else {
        data = jsonObject;
      }
      console.log("data: ", data);

      axios
        .post(
          `http://localhost:8080/api/v1/config?conf=${configFileName}`,
          data
        )
        .then((res) => {
          console.log("response111: ", res.data);
          axios
            .post(
              // `http://localhost:8080/api/v1/execute?re-run=${reRun}&conf=${configFileName}&parallel=${parallel}&source=${source}&dry-run=${dryRun}`,
              `http://localhost:8080/api/v1/execute?re-run=${reRun}&conf=${configFileName}&parallel=${parallel}&source=${source.toLowerCase()}`,
              dataBodyObj
            )
            .then((res) => {
              console.log("response2222: ", res.data);
            })
            .catch((err) => {
              console.log(err);
            });
        })
        .catch((err) => {
          console.log(err);
        });
    }
  };

  const getModulesAPI = (e) => {
    axios
      .get("http://localhost:8080/api/v1/modules")
      .then((res) => {
        console.log(res.data);
        setModules(res.data);
      })
      .catch((err) => {
        console.log(err);
      });
  };

  useEffect(() => {
    getModulesAPI();
  }, []);

  return (
    <div style={{ marginLeft: "10%", width: "80%" }}>
      <div style={{ marginBottom: "2%", marginTop: "2%", textAlign: "center" }}>
        <h1>Mozart</h1>
        {/* <p>Subtext goes here</p> */}
      </div>
      <div>
        <Form>
          {/* <FormGroup>
            <FormLabel>
              <Tooltip triggerText="Config file name">
                Create a suitable name for your YAML config file.
              </Tooltip>
            </FormLabel>
            <TextInput
              id="ibmConfiguration-textInput-newModuleName"
              placeholder="Enter a name for your config/YAML file here"
              onChange={(e) => {
                setConfigFileName(e.target.value);
              }}
            />
          </FormGroup> */}

          <FormGroup>
            <Tooltip triggerText="Config file details">
              Enter the JSON for your config file.
            </Tooltip>
            <TextArea
              placeholder="Paste JSON here or upload json file"
              onChange={(e) => {
                setJsonObject(e.target.value);
              }}
              {...validateTextArea ===true ? {...someProps} : ""}
            ></TextArea>
            {/* <FileUploader
              multiple
              accept={[".json"]}
              buttonKind="ghost"
              buttonLabel="Upload .json files at 500mb or less"
              filenameStatus="edit"
              onChange={(e) => {
                setJsonFile(e.target.files[0] || {});
              }}
            /> */}
          </FormGroup>

          {/* <FormGroup>
            <FormLabel>
              <Tooltip triggerText="Select whether you want to create a new module or run an existing module">
                Create a suitable name for your module.
              </Tooltip>
            </FormLabel>
            <RadioButtonGroup defaultSelected="existing">
              <RadioButton
                value="existing"
                labelText="existing module"
                id="existing"
                onClick={(e) => {
                  setTypeOfModule(e.target.value);
                }}
              />
              <RadioButton
                value="new"
                labelText="new module"
                id="new"
                onClick={(e) => {
                  setTypeOfModule(e.target.value);
                }}
              />
            </RadioButtonGroup>
          </FormGroup> */}

          {/* {typeOfModule === "existing" ? ( */}
            <FormGroup>
              <FormLabel>
                <Tooltip triggerText="Module">
                  Select a module you want to run
                </Tooltip>
              </FormLabel>
              <Dropdown
                items={modules}
                // label="Select an existing module"
                label={modules[0]}
                defaultValue={modules[0]}
                defaultSelected={modules[0]}
                onChange={(e) => {
                  setSelectedModule(e.selectedItem);
                }}
              />
            </FormGroup>
          {/* ) : ( */}
            {/* <FormGroup>
              <FormLabel>
                <Tooltip triggerText="New Module">
                  Create a suitable name for your module.
                </Tooltip>
              </FormLabel>
              <TextInput
                id="ibmConfiguration-textInput-newModuleName"
                placeholder="Enter new module name"
                onChange={(e) => {
                  setNewModuleName(e.target.value);
                }}
              />
            </FormGroup> */}
          {/* )} */}

          <FormGroup>
            {/* <Tooltip triggerText="Type of execution">
              Dry Run shows what scripts will run, but does not run the scripts.
            </Tooltip>
            <Checkbox
              labelText="Dry Run"
              id="dry-run"
              onClick={(e) => {
                dryRun === false ? setDryRun(true) :setDryRun(false)
              }}
            /> */}
            <Tooltip triggerText="Type of execution">
              Re Run runs all the scripts from initial state, ignoring
              previously saved state. (Check re-run if you have already ran your
              scripts once and want to run it again ignoring its previous
              state.)
            </Tooltip>
            <Checkbox
              labelText="Re Run"
              id="re-run"
              onClick={(e) => {
                reRun === false ? setReRun(true) :setReRun(false)
              }}
            />
          </FormGroup>

          <FormGroup>
            <FormLabel>
              <Tooltip triggerText="Type of Run">
              Select how you want to run your files, sequentially or parallelly. Default runs sequentially.
              </Tooltip>
            </FormLabel>
            <Checkbox
              labelText="Parallel"
              id="parallel"
              onClick={(e) => {
                parallel === false ? setParallel(true) :setParallel(false)
              }}
            />
          </FormGroup>
          
          <FormGroup>
            <FormLabel>
              <Tooltip triggerText="Source file type">
                Select source file type [python, bash]
              </Tooltip>
            </FormLabel>
            <Dropdown
              items={allSourceFileTypes}
              // label="Select your source file type"
              label={allSourceFileTypes[0]}
              defaultValue={allSourceFileTypes[0]}
              onChange={(e) => {
                setSourceFileTypes(e.selectedItem);
              }}
            />
          </FormGroup>

          <FormGroup>
            <FormLabel>
              <Tooltip triggerText="OS">Select OS [Darwin, Linux]</Tooltip>
            </FormLabel>
            <Dropdown
              items={allOS}
              // label="Select your OS"
              label={allOS[0]}
              defaultValue={allOS[0]}
              onChange={(e) => {
                setOS(e.selectedItem);
              }}
            />
          </FormGroup>
          
        </Form>

        <div style={{ marginTop: "2%" }}>
          <Button kind="secondary">Cancel</Button>
          <span style={{ marginLeft: "2%" }}>
            <Button onClick={makeSampleAPICall}>Deploy</Button>
          </span>
        </div>
      </div>
    </div>
  );
}
