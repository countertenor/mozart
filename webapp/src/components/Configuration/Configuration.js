import React, { Component } from "react";
import {
  FormLabel,
  TextInput,
  Tooltip,
  Dropdown,
  Link
} from "carbon-components-react";

const { PasswordInput } = TextInput;

export default class Configuration extends Component {
  render() {
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
          <Dropdown sel label="Select a region" />
          <FormLabel>
            <Tooltip triggerText="Availability zone">
              <p id="tooltip-body">IBM data centers.</p>
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
          <Dropdown sel label="Select a zone" />
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
}
