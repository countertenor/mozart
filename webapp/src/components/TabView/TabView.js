import React, { useState, useEffect, useReducer } from "react";
import { useHistory } from "react-router-dom";
import { Button, Modal, ToastNotification } from "carbon-components-react";
import "./TabView.scss";
import axios from "axios";
import { server_hostname, server_port } from "../../constants/Constants";
import Configuration from "../Configuration/Configuration";
import Execution from "../Status/Install";
import { NOTIFICATION_INIT, NOTIFICATION_REDUCER } from '../../constants/reducer-constants';

export default function TabView (){
  let [activeTab, setActiveTab] = useState("configuration")

  let switchActiveTab = (tab) => {
    if(tab == "execution" || tab == "configuration"){
      setActiveTab(tab)
      console.log("heyyyyyy: ", tab)
    }
    else{
      setActiveTab(tab.target.value);
      console.log("heyyyyyy: ", tab.target.value)
    }
  };

  const [notification, notificationDispatch] = useReducer(NOTIFICATION_REDUCER, NOTIFICATION_INIT);

    return (
      <div className="screen-style">
        <h2 className="heading">Mozart</h2>
        <div class="left-col"></div>
        <div className="form-style">
          <div className="button-toggle">
            {activeTab === "configuration" ? (
              <Button className="b1-active" value="configuration" onClick={switchActiveTab}>Step 1: Configuration</Button>
            ) : (
              <Button className="b1" value="configuration" onClick={switchActiveTab}>Step 1: Configuration</Button>
            )}
            {activeTab === "execution" ? (
              <Button className="b2-active" value="execution" onClick={switchActiveTab}>Step 2: Execution</Button>
            ) : (
              <Button className="b2" value="execution" onClick={switchActiveTab}>Step 2: Execution</Button>
            )}
          </div>
          <hr />
          {activeTab === "configuration" ? (
            <Configuration switchActiveTab={switchActiveTab} />
          ) : activeTab === "execution" ? (
            <Execution notificationDispatch={notificationDispatch} switchActiveTab={switchActiveTab} />
          ) : (
            <div></div>
          )}
        </div>
        <div class="right-col">
        </div>
        <div className="footer"></div>
      </div>
    );
  }