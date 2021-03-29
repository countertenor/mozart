import React, { useState, useReducer } from "react";
import { Button } from "carbon-components-react";
import "./TabView.scss";
// import { server_hostname, server_port } from "../../constants/Constants";
import Configuration from "../Configuration/Configuration";
import Execution from "../Status/Install";

export default function TabView (){
  let [activeTab, setActiveTab] = useState("configuration")
  let [moduleName, setModuleName] = useState("")

  let switchActiveTab = (tab) => {
    if(tab === "execution" || tab === "configuration"){
      setActiveTab(tab)
      console.log("heyyyyyy: ", tab)
    }
    else{
      setActiveTab(tab.target.value);
      console.log("heyyyyyy: ", tab.target.value)
    }
  };

  let passModuleName = (moduleName) => {
      setModuleName(moduleName)
      console.log("moduleName: ", moduleName)
  };

  let getModuleName = () =>{
    return moduleName;
  }

    return (
      <div className="screen-style">
        <h2 className="heading">Mozart</h2>
        <div className="left-col"></div>
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
            <Configuration switchActiveTab={switchActiveTab} passModuleName={passModuleName}/>
          ) : activeTab === "execution" ? (
            <Execution switchActiveTab={switchActiveTab} getModuleName={getModuleName} />
          ) : (
            <div></div>
          )}
        </div>
        <div className="right-col">
        </div>
        <div className="footer"></div>
      </div>
    );
  }
