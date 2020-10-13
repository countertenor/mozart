import React, { useState } from "react";
import {
  FormLabel,
  TextInput,
  Tooltip,
  Dropdown,
  Link,
  Form,
  Button
} from "carbon-components-react";
import axios from 'axios';

const { PasswordInput } = TextInput;

export default function Configuration () {
  const [clusterName, setClusterName] = useState("some cluster name?");
  const [apiKey, setApiKey] = useState("some apiKey?");
  const [resourceGroup, setResourceGroup] = useState("some resourceGroup?");
  const [publicVLAN, setPublicVLAN] = useState("some publicVLAN?");
  const [privateVLAN, setPrivateVLAN] = useState("some privateVLAN?");

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
  ])
  const [region, setRegion] = useState(regionAPIResponseOptions[0])


//   const updateFieldChanged = e => {
//     console.log("hey: ",items);
//     let newArr = [...items]; // copying the old datas array
//     newArr[newArr.length] = e.selectedItem; // replace e.target.value with whatever you want to change it to
//     setRegion(newArr); // ??
//     console.log("hello", items);
// }

const makeSampleAPICall = e =>{
  const data = {
    "clusterName":clusterName,
    "apiKey":apiKey,
    "resourceGroup":resourceGroup,
    "publicVLAN":publicVLAN,
    "privateVLAN":privateVLAN,
    "region":region.value
  }
  console.log("options: ",region)
  console.log("data: ",data)

  axios.post('http://localhost:8080/api/v1/', data).then((res) => {
      console.log(res.data);
    }).catch((err) => {
      console.log(err);
    });
}

  return (
    <div style={{ marginLeft: "10%", width: "80%" }}>
      <div style={{ marginBottom: "2%", marginTop: "2%", textAlign: "center" }}>
        <h1>Configure Install</h1>
        <p>Subtext goes here</p>
      </div>
      <div>
        <Form>
          <FormLabel>
            <Tooltip triggerText="Cluster name">
              Create a suitable name for your cluster.
            </Tooltip>
          </FormLabel>
          <TextInput
            id="ibmConfiguration-textInput-clusterName"
            placeholder="Cluster name"
            defaultValue={clusterName}
            onChange={(e) => {
              console.log("clustername: ",e.target.value);
              setClusterName(e.target.value);
            }}
          />

          <FormLabel>
          <Tooltip triggerText="Region">
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
          items={regionAPIResponseOptions}
          label="Select a region"
          defaultValue={region}
          onChange={(e) => {
            // updateFieldChanged(e)
            setRegion(e.selectedItem);
            console.log(e.selectedItem);
          }}
        />

          <FormLabel>
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
          />
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