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
            {/* <Button className="logout" onClick={this.logout}>
              Logout
            </Button> */}
          </div>
          <hr />
          {activeTab === "configuration" ? (
            <Configuration switchActiveTab={switchActiveTab} />
          ) : activeTab === "execution" ? (
            <Execution notificationDispatch={notificationDispatch} />
          ) : (
            <div></div>
          )}
        </div>
        <div class="right-col">
          {/* {this.state.notification.result !== "RUNNING" && this.state.catch_error === false  ? (
            <ToastNotification
              caption=""
              hideCloseButton={false}
              iconDescription="describes the close button"
              kind={this.state.notification.type}
              notificationType="toast"
              onCloseButtonClick={function noRefCheck() {}}
              role="alert"
              statusIconDescription="describes the status icon"
              style={{
                marginBottom: "0.5rem",
                maxWidth: "20rem",
              }}
              subtitle={
                this.state.notification.type === "info" ?
                `Your build is taking too long! You will recieve an email once it is complete.` :
                this.state.notification.build_number === null || this.state.notification.build_number === undefined
                  ? "Your build was a " + this.state.notification.result
                  : "Your build " +
                    this.state.notification.build_number +
                    " was a " +
                    this.state.notification.result
              }
              timeout={0}
              title={this.state.notification.result}
            />
          ) : null} */}
          {/* {this.state.error === true ? (
            <Modal
            className="error"
            iconDescription="Close"
            modalHeading="ERROR!"
            onRequestClose={() => {
              this.setState({ error: false });
            }}
            open
            passiveModal
          >
            <p>Oops something went wrong!</p>
          </Modal>
          ) : null} */}
        </div>
        <div className="footer"></div>
      </div>
    );
  }
