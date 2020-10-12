import React, { useState } from "react";
import {
  FormLabel,
  TextInput,
  Tooltip,
  Dropdown,
  Link
} from "carbon-components-react";

const { PasswordInput } = TextInput;

export default function Configuration () {

  const [clusterName, setClusterName] = useState("some cluster name?");
  const [items, setItems] = useState([
    {
      label: "option-0",
      value:
        "Lorem, ipsum dolor sit amet consectetur adipisicing elit. Vitae, aliquam. Blanditiis quia nemo enim voluptatibus quos ducimus porro molestiae nesciunt error cumque quaerat, tempore vero unde eum aperiam eligendi repellendus.",
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

//   const updateFieldChanged = e => {
//     console.log("hey: ",items);
//     let newArr = [...items]; // copying the old datas array
//     newArr[newArr.length] = e.selectedItem; // replace e.target.value with whatever you want to change it to
//     setItems(newArr); // ??
//     console.log("hello", items);
// }

  return (
    <div>
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
          console.log(e.target.value)
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
      items={items}
      label="Select a region"
      onChange = {(e)=>{
        // updateFieldChanged(e)
        setItems(e.selectedItem)
        console.log(e.selectedItem)
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
      />
      <FormLabel>
        <Tooltip triggerText="Resource group (Optional)">
          <p id="tooltip-body">
            Provide your existing resource group name to organize your account
            resources
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
      />
    </div>
  );
}