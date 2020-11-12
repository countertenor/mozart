import React, { useState, useEffect } from "react";
import {
  FormLabel,
  Tooltip,
  Dropdown,
  Form,
  Button,
  FormGroup,
  TextArea,
  Checkbox,
  Modal,
} from "carbon-components-react";
import axios from "axios";

export default function Configuration(props) {
  let [jsonObject, setJsonObject] = useState("")

  let [modules, setModules] = useState([""]);
  let [moduleName, setModuleName] = useState(modules.length > 0 ? [modules[0]] : []);

  let [reRun, setReRun] = useState(false);
  let [parallel, setParallel] = useState(false);

  let [networkError, setNetworkError] = useState("");
  let [openModal, setOpenModal] = useState(false);

  const someProps = {
    invalid: true,
    invalidText: "This value cannot be empty. You must enter a valid json object here.",
  };

  const propsForModuleName = {
    invalid: true,
    invalidText: "This value cannot be empty. The module list is not populated. Refresh the page. If issue persists please contact Mozart team.",
  };

  let [validateTextArea, setValidateTextArea] = useState(false);
  let [validateModuleList, setModuleList] = useState(false);

  function IsJsonString(str) {
    try {
        JSON.parse(str);
    } catch (e) {
        return false;
    }
    return true;
}

  const makeSampleAPICall = (e) => {
    if (networkError.length > 0) {
      setOpenModal(true);
    }
    else{
    e.preventDefault();
    console.log(Object.keys(jsonObject).length)
    console.log(typeof jsonObject)
    if (IsJsonString(jsonObject) === false || Object.keys(jsonObject).length === 0) {
      console.log("ERROR!");
      setValidateTextArea(true);
    } 
    if(moduleName.length <= 1){
      console.log("ERROR!", moduleName.length);
      setModuleList(true)
    }
    else {
      setValidateTextArea(false);
      setModuleList(false)
      let configDataObject = JSON.parse(jsonObject || "{}");
      const dataBodyObj = {
        moduleName: moduleName
      };
      axios
        .post(`http://localhost:8080/api/v1/config`,configDataObject)
        .then((res) => {
          console.log("config response: ", res.data);
          axios
            .post(`http://localhost:8080/api/v1/execute?re-run=${reRun}&parallel=${parallel}}`,dataBodyObj)
            .then((res) => {
              console.log("execute response: ", res.data);
              props.switchActiveTab("execution")
            })
            .catch((err) => {
              console.log(err);
            });
        })
        .catch((err) => {
          console.log(err);
        });
    }
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
        setNetworkError("err");
        console.log(err);
      });
  };

  useEffect(() => {
    getModulesAPI();
  }, []);

  return (
    <div>
      <div style={{ marginBottom: "2%", marginTop: "2%" }}>
      </div>
      <div style={{backgroundColor: "#f4f4f4", minWidth:"8rem", minHeight:"4rem", padding: "1.625rem 1.625rem 2.125rem 1.625rem"}}>
        <Form>

          <FormGroup>
            <Tooltip triggerText="Config file details">
              Enter the JSON for your config file.
            </Tooltip>
            <TextArea
              placeholder="Paste JSON here or upload json file"
              defaultValue="{}"
              onChange={(e) => {
                setJsonObject(e.target.value);
              }}
              {...validateTextArea ===true ? {...someProps} : ""}
            ></TextArea>
          </FormGroup>
            <FormGroup>
              <FormLabel>
                <Tooltip triggerText="Module">
                  Select a module you want to run
                </Tooltip>
              </FormLabel>
              <Dropdown
                items={modules}
                label="Select a module to run"
                defaultValue={modules[0]}
                defaultSelected={modules[0]}
                onChange={(e) => {
                if (networkError.length > 0) {
                  setOpenModal(true);
                }
                  setModuleName(e.selectedItem);
                }}
                {...validateModuleList ===true ? {...propsForModuleName} : ""}
              />
            </FormGroup>
          <FormGroup>
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
                reRun === false ? setReRun(true) : setReRun(false);
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
          
        </Form>

        <div style={{ marginTop: "2%" }}>
          <Button kind="secondary">Cancel</Button>
          <span style={{ marginLeft: "2%" }}>
            <Button onClick={makeSampleAPICall}>Deploy</Button>
          </span>
        </div>
      </div>
      <div>
        {openModal === true ? (
          <div>
            <Modal
              className="error-modal"
              iconDescription="Close"
              modalHeading="NETWORK ERROR!"
              onRequestClose={() => {
                setOpenModal(false);
              }}
              open
              passiveModal
            >
              <p>You are not connected to the server!</p>
            </Modal>
          </div>
        ) : (
          <div></div>
        )}
      </div>
    </div>
  );
}
