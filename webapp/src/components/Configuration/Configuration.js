import React, { useState, useEffect } from "react";
import {
  FormLabel,
  TextInput,
  Tooltip,
  Dropdown,
  Link,
  Form,
  Button,
  RadioButton,
  RadioButtonGroup,
  FormGroup,
  TextArea,
  FileUploader,
} from "carbon-components-react";
import axios from "axios";

const { PasswordInput } = TextInput;

export default function Configuration() {
  const [moduleName, setModuleName] = useState("default-module-name");
  // const [apiKey, setApiKey] = useState("some apiKey?");
  // const [resourceGroup, setResourceGroup] = useState("some resourceGroup?");
  // const [publicVLAN, setPublicVLAN] = useState("some publicVLAN?");
  // const [privateVLAN, setPrivateVLAN] = useState("some privateVLAN?");

  const [regionAPIResponseOptions] = useState([
    {
      label: "option-0",
      value: "Option 0",
    },
    {
      label: "option-1",
      value: "Option 1",
    },
    {
      label: "option-2",
      value: "Option 2",
    },
    {
      label: "option-3",
      value: "Option 3",
    },
  ]);

  const allSourceFileTypes = ["Python", "Bash"];
  const allExtensions = ["py", "sh"];
  const allOS = ["Darwin", "Linux"];
  const alltypesOfRun = ["Parallel", "Synchronous"];

  const [region, setRegion] = useState(regionAPIResponseOptions[0]);

  const [modules, setModules] = useState([""]);
  let [selectedModule] = useState(modules.length > 0 ? [modules[0]] : [""]);

  const [sourceFileTypes, setSourceFileTypes] = useState(allSourceFileTypes)
  const [extensions, setExtensions] = useState(allExtensions)
  const [os, setOS] = useState(allOS)
  const [typeOfRun, setTypeOfRun] = useState(alltypesOfRun)

  //   const updateFieldChanged = e => {
  //     console.log("hey: ",items);
  //     let newArr = [...items]; // copying the old datas array
  //     newArr[newArr.length] = e.selectedItem; // replace e.target.value with whatever you want to change it to
  //     setRegion(newArr); // ??
  //     console.log("hello", items);
  // }

  const makeSampleAPICall = (e) => {
    const data = {
      moduleName: moduleName,
      // apiKey: apiKey,
      // resourceGroup: resourceGroup,
      // publicVLAN: publicVLAN,
      // privateVLAN: privateVLAN,
      // region: region.value,
      selectedModule: selectedModule,
    };
    console.log("options: ", region);
    console.log("data: ", data);

    axios
      .post("http://localhost:8080/api/v1/", data)
      .then((res) => {
        console.log(res.data);
      })
      .catch((err) => {
        console.log(err);
      });
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
        <h1>Configure Install</h1>
        <p>Subtext goes here</p>
      </div>
      <div>
        <Form>
          <FormGroup>
            <FormLabel>
              <Tooltip triggerText="Select whether you want to create a new module or run an existing module">
                Create a suitable name for your module.
              </Tooltip>
            </FormLabel>
            <RadioButtonGroup
              // defaultSelected="default-selected"
              valueSelected={(e) => {
                console.log("check: ", e);
              }}
              legend="Group Legend"
            >
              <RadioButton
                value="default-selected"
                labelText="existing module"
                id="existing"
                // checked={(e) => {
                //   console.log(e.checked);
                // }}
                // onChange={(e) => {
                //   console.log(e.target.value);
                // }}
              />
              {/* <RadioButton value="standard" labelText="new module" id="new" /> */}
            </RadioButtonGroup>
          </FormGroup>

          {/* <FormGroup>
            <FormLabel>
              <Tooltip triggerText="New Module">
                Create a suitable name for your module.
              </Tooltip>
            </FormLabel>
            <TextInput
              id="ibmConfiguration-textInput-moduleName"
              placeholder="Enter new module name"
              onChange={(e) => {
                console.log("moduleName: ", e.target.value);
                setModuleName(e.target.value);
              }}
            />
          </FormGroup> */}

          <FormGroup>
            <FormLabel>
              <Tooltip triggerText="Module">
                <p id="tooltip-body">IBM Cloud locations.</p>
                <div className="bx--tooltip__footer">
                  <Link
                    href="https://cloud.ibm.com/docs/containers?topic=containers-regions-and-zones"
                    target="_blank"
                  >
                    Learn more
                  </Link>
                </div>
              </Tooltip>
            </FormLabel>
            <Dropdown
              items={modules}
              label="Select an existing module"
              defaultValue={modules[0]}
              onChange={(e) => {
                selectedModule = e.selectedItem;
                console.log(e.selectedItem);
              }}
            />
          </FormGroup>

          <FormGroup>
            <TextArea placeholder="Paste JSON here or upload json file"></TextArea>
            {/* <Button kind="ghost">Upload JSON file</Button> */}
            <FileUploader
              multiple
              accept={[".json"]}
              buttonKind="ghost"
              buttonLabel="Upload .json files at 500mb or less"
              filenameStatus="edit"
            />
          </FormGroup>

          <FormGroup>
            <FormLabel>
              <Tooltip triggerText="Type of execution">
                Shows what scripts will run, but does not run the scripts
              </Tooltip>
            </FormLabel>
            <RadioButtonGroup
              // defaultSelected="default-selected"
              valueSelected={(e) => {
                console.log("check: ", e);
              }}
            >
              <RadioButton
                value="default-selected"
                labelText="dry run"
                id="d"
              />
              <RadioButton value="standard" labelText="run" id="r" />
            </RadioButtonGroup>
          </FormGroup>

          <FormGroup>
            <FormLabel>
              <Tooltip triggerText="Source file type">
                Select source file type [python, bash]
              </Tooltip>
            </FormLabel>
            <Dropdown
              items={sourceFileTypes}
              label="Select your source file type"
              defaultValue={sourceFileTypes[0]}
              onChange={(e) => {
                selectedModule = e.selectedItem;
                console.log(e.selectedItem);
              }}
            />
          </FormGroup>

          <FormGroup>
            <FormLabel>
              <Tooltip triggerText="File extension">
                Select file extension [py, sh]
              </Tooltip>
            </FormLabel>
            <Dropdown
              items={extensions}
              label="Select your file extension"
              defaultValue={extensions[0]}
              onChange={(e) => {
                selectedModule = e.selectedItem;
                console.log(e.selectedItem);
              }}
            />
          </FormGroup>
          <FormGroup>
            <FormLabel>
              <Tooltip triggerText="OS">
                Select OS [Darwin, Linux]
              </Tooltip>
            </FormLabel>
            <Dropdown
              items={os}
              label="Select your OS"
              defaultValue={os[0]}
              onChange={(e) => {
                selectedModule = e.selectedItem;
                console.log(e.selectedItem);
              }}
            />
          </FormGroup>
          <FormGroup>
            <FormLabel>
              <Tooltip triggerText="Type of Run">
                Select type of run [Parallel, Synchronous]
              </Tooltip>
            </FormLabel>
            <Dropdown
              items={typeOfRun}
              label="Select your type of run"
              defaultValue={typeOfRun[0]}
              onChange={(e) => {
                selectedModule = e.selectedItem;
                console.log(e.selectedItem);
              }}
            />
          </FormGroup>

          {/* <FormLabel>
            <Tooltip triggerText="API key">
              <p id="tooltip-body">IBM Cloud user API key.</p>
              <div className="bx--tooltip__footer">
                <Link
                  href="https://cloud.ibm.com/docs/iam?topic=iam-userapikey"
                  target="_blank"
                >
                  Learn more
                </Link>
              </div>
            </Tooltip>
          </FormLabel>
          <PasswordInput
            id="ibmConfiguration-textInput-apiKey"
            placeholder="IBM Cloud API key"
            showPasswordLabel="Show"
            hidePasswordLabel="Hide"
            defaultValue={apiKey}
            onChange={(e) => {
              console.log("apiKey: ", e.target.value);
              setApiKey(e.target.value);
            }}
          />

          <FormLabel>
            <Tooltip triggerText="Resource group (Optional)">
              <p id="tooltip-body">
                Provide your existing resource group name to organize your
                account resources
              </p>
              <div className="bx--tooltip__footer">
                <Link
                  href="https://cloud.ibm.com/account/resource-groups"
                  target="_blank"
                >
                  Learn more
                </Link>
              </div>
            </Tooltip>
          </FormLabel>
          <TextInput
            id="ibmConfiguration-textInput-resourceGroup"
            placeholder="Resource group"
            defaultValue={resourceGroup}
            onChange={(e) => {
              console.log("resourceGroup: ", e.target.value);
              setResourceGroup(e.target.value);
            }}
          />

          <FormLabel>
            <Tooltip triggerText="Public VLAN (Optional)">
              <p id="tooltip-body">Provide existing Public VLAN id</p>
              <div className="bx--tooltip__footer">
                <Link
                  href="https://cloud.ibm.com/classic/network/vlans"
                  target="_blank"
                >
                  Learn more
                </Link>
              </div>
            </Tooltip>
          </FormLabel>
          <TextInput
            id="ibmConfiguration-textInput-publicVLAN"
            placeholder="Public VLAN"
            defaultValue={publicVLAN}
            onChange={(e) => {
              console.log("publicVLAN: ", e.target.value);
              setPublicVLAN(e.target.value);
            }}
          />

          <FormLabel>
            <Tooltip triggerText="Private VLAN (Optional)">
              <p id="tooltip-body">Provide existing Private VLAN id</p>
              <div className="bx--tooltip__footer">
                <Link
                  href="https://cloud.ibm.com/classic/network/vlans"
                  target="_blank"
                >
                  Learn more
                </Link>
              </div>
            </Tooltip>
          </FormLabel>
          <TextInput
            id="ibmConfiguration-textInput-privateVLAN"
            placeholder="Private VLAN"
            defaultValue={privateVLAN}
            onChange={(e) => {
              console.log("privateVLAN: ", e.target.value);
              setPrivateVLAN(e.target.value);
            }}
          /> */}
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
