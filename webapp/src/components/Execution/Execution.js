import React, { useState, useEffect, useReducer } from "react";
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
  Tabs,
  Tab,
} from "carbon-components-react";
import axios from "axios";
import Install from '../Status/Install';
import Configuration from '../Configuration/Configuration';
import { NOTIFICATION_INIT, NOTIFICATION_REDUCER } from '../../constants/reducer-constants';

export default function Execution() {
    const [notification, notificationDispatch] = useReducer(NOTIFICATION_REDUCER, NOTIFICATION_INIT);

  
  return (
    <div style={{ marginLeft: "10%", width: "80%" }}>
      <div style={{ marginBottom: "2%", marginTop: "2%", textAlign: "center" }}>
        <h1>Mozart Status SPA Template</h1>
        {/* <p>Subtext goes here</p> */}
      </div>
      <div>
        <div>
          <Tabs>
            <Tab id="tab-1" label="Tab label 1">
              <div className="some-content">
              <Configuration />
              </div>
            </Tab>
            <Tab id="tab-2" label="Tab label 2">
              <div className="some-content">
              <Install notificationDispatch={notificationDispatch}/>
              </div>
            </Tab>
          </Tabs>
        </div>

        <div style={{ marginTop: "2%" }}>
          <Button kind="secondary">Cancel</Button>
          <span style={{ marginLeft: "2%" }}>
            <Button>Deploy</Button>
          </span>
        </div>
      </div>
    </div>
  );
}
