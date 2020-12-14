import React, { useCallback, useEffect, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import {
  Button, CodeSnippet, InlineLoading, Loading, Modal, Tile, Accordion, AccordionItem
} from 'carbon-components-react';
import { CheckmarkFilled16, Misuse16, View16, WarningAltFilled16 } from '@carbon/icons-react';
import styles from './Install.module.scss';
import { getStatus, setupWS } from './InstallUtils';
import axios from "axios";

function LogModal({ onRequestClose, open, log }) {
  return (
    <Modal
      onRequestClose={onRequestClose}
      passiveModal
      modalHeading="Logs"
      open={open}
      size="lg"
    >
      <CodeSnippet
        type="multi"
        light
        className={styles.codeSnippet}
      >
        {log}
      </CodeSnippet>
    </Modal>
  );
}

LogModal.propTypes = {
  onRequestClose: PropTypes.func.isRequired,
  open: PropTypes.bool.isRequired,
  log: PropTypes.string.isRequired
};

function Task({
  duration, taskName, state, openLogModal, logFilePath
}) {
  let statusIcon = null;
  switch (state) {
  case 'running':
    statusIcon = (
      <div className={styles.loaderIconHolder}>
        <InlineLoading />
      </div>
    );
    break;

  case 'error':
    statusIcon = (
      <div className={styles.loaderIconHolder}>
        <Misuse16 className={styles.errorFilled} />
      </div>
    );
    break;
  case 'success':
    statusIcon = (
      <div className={styles.loaderIconHolder}>
        <CheckmarkFilled16 className={styles.successFilled} />
      </div>
    );
    break;
  case 'timeout':
    statusIcon = (
      <div className={styles.loaderIconHolder}>
        <WarningAltFilled16 className={styles.warningFilled} />
      </div>
    );
    break;
  case 'skipped':
    statusIcon = (
      <div className={styles.loaderIconHolder}>
        <Misuse16 className={styles.notStartedFilled} />
      </div>
    );
    break;
  default:
    statusIcon = (
      <div className={styles.loaderIconHolder}>
        <CheckmarkFilled16 className={styles.notStartedFilled} />
      </div>
    );
  }

  return (
    <li className={styles.Task}>
      <div className={styles.statusMsg}>
        <p>{taskName}</p>   {/*steps*/}
        {/* <p>{logFilePath}</p> */}
        {/* {duration && <p>{`Completion time: ${duration} minutes.`}</p>} */}
      </div>
      {statusIcon}
      <Button
        renderIcon={View16}
        hasIconOnly
        kind="ghost"
        iconDescription="View logs"
        tooltipPosition="top"
        onClick={() => openLogModal(logFilePath)}
      />
    </li>
  );
}

Task.propTypes = {
  duration: PropTypes.number.isRequired,
  taskName: PropTypes.string.isRequired,
  logFilePath: PropTypes.string.isRequired,
  state: PropTypes.string.isRequired,
  openLogModal: PropTypes.func.isRequired,
};

export default function Install(props) {
  const intervalRef = useRef(null);
  const logRef = useRef('');
  const [steps, setSteps] = useState([]);
  const [state, setState] = useState("");
  const [moduleName] = useState(props.getModuleName()+"");
  let stringPrint=" - "
  for(let i=0; i<30; i++){
    stringPrint+=" - "
  }
  const [print] = useState(stringPrint)
  const [percentage, setPercentage] = useState({});
  const [moduleState, setModuleState] = useState({});
  const [logs, setLogs] = useState("")
  
  // steps is an array of objects with keys directory, module, tasks
  // tasks is an array of objects with keys taskName and status
  const [curStatus, setCurStatus] = useState('loading');
  const [logModalIsOpen, setLogModalIsOpen] = useState(false);
  const [logSocket, setLogSocket] = useState(null);

  const openLogModal = useCallback((logFilePath) => {
    // TODO: Change the path to passed in
    console.log(`logFilePath: ≥≥≥ ${__filename}`);
    const socket = setupWS(`${logFilePath}`, (err, streamData) => {
      if (err) {
        return console.log('Handle error TODO');
      }

      logRef.current = logRef.current.concat(`${streamData}\n`);
      setLogModalIsOpen(true);
      setLogs(logRef.current)
    });

    setLogSocket(socket);
  }, [setLogModalIsOpen]);

  const closeLogModal = useCallback(() => {
    if (logSocket) {
      logSocket.close();
      setLogSocket(null);
    }

    logRef.current = '';
    setLogModalIsOpen(false);
  }, [logSocket, setLogModalIsOpen]);

  const clearStatusInterval = useCallback(() => {
    if (intervalRef.current) {
      console.log(`clear interval: ${intervalRef.current}`);
      clearInterval(intervalRef.current);
      intervalRef.current = null; // now that it's clear no longer relevant
    } else {
      console.log('status already cleared');
    }
  }, []);

  const getData = useCallback(async () => {
    let countObj={}
    let moduleStateObj={}
    getStatus(moduleName, (err, data) => {
      console.log("check props?: ",moduleName)
      if (err) {
        console.log(err)
      }
      setSteps((data||{}).steps||[]);
      const st = (data||{}).state||"";
      console.log(st, " :: ", intervalRef.current)     //every 3 seconds
      setState(st)
      console.log("hello",(data||{}).steps||[])
      setCurStatus(''); // TODO: Figure out
      ((data||{}).steps||[]).map( e => {
        countObj = {...countObj,
            [e.directory]: 0
        }
        moduleStateObj = {...moduleStateObj,
          [e.directory]: "pending"
        }
        let count = 0;
        e.tasks.map(item =>{
          if(item.status.state === "success"){
            countObj = {...countObj,
              [e.directory]: Math.trunc(++count * 100/e.tasks.length)
            }
          }
          if(["error","timeout","skipped","canceled"].indexOf(moduleStateObj[e.directory])<0){
            if(!!item.status.state){
              moduleStateObj = {...moduleStateObj,
                [e.directory]: item.status.state
              }
            }
          }
        })
      })
      setModuleState(moduleStateObj)
      setPercentage(countObj);
    });
  }, [clearStatusInterval]);

  // const getModulesAPI = (e) => {
  //   axios
  //     .get("http://localhost:8080/api/v1/modules")
  //     .then((res) => {
  //       console.log(res.data);
  //       setModules(res.data);
  //     })
  //     .catch((err) => {
  //       setNetworkError("err");
  //       console.log(err);
  //     });
  // };

  const cancel = (e) =>{
    console.log("Cancel pressed!", props)
    // e.preventDefault();
    axios
    .put(`http://localhost:8080/api/v1/cancel`)
    .then((res) => {
      console.log("response cancel??: ", res.data);
    })
    .catch((err) => {
      console.log("error cancel",err);
    });
  }

  const setStatusInterval = useCallback(() => {
    if (intervalRef.current) {
      // Clear any current status interval if there is one
      console.log(`interval already set at: ${intervalRef.current}, clearing`);
      clearInterval(intervalRef.current);
    }
    const intervalId = setInterval(getData, 3000);
    console.log(`set interval: ${intervalId}`);
    intervalRef.current = intervalId;
  }, [getData]);

  useEffect(() => {
    getData();        //gets data.steps array of tasks
    setStatusInterval();
    return clearStatusInterval;
  }, []); // eslint-disable-line

  useEffect(() => {
    if(state === "completed"){
      clearStatusInterval();
    }
  }, [state])

  const header =
    steps.length > 0 ? (
      <h6>
        Status: {state} {moduleName}
      </h6>
    ) : <h6>
    Go to Configuration tab to run your script, you will then see the status here.
  </h6>;

  return (
    <div style={{ marginBottom: "2%", marginTop: "2%" }}>
      <Tile className={styles.Install}>
        {curStatus === "loading" ? (
          <>
            {header}
            <Loading active small={false} />
          </>
        ) : (
          <>
            {header}
            <Accordion size="xl" align="start">
              {steps.map((module) => {
                return (
                  <div key={module.directory} className={styles.module}>
                    <AccordionItem title= {<div>{module.directory.split("/").slice(1).join(" | ")}<span style={{float:"right"}}>{moduleState[module.directory]+ " | "}{percentage[module.directory]+"% complete"}</span></div>}>
                      {module.tasks.map((task) => (
                        <ul
                          className={styles.loadersHolder}
                          key={task.taskName}
                        >
                          <Task
                            duration="TODO"
                            taskName={task.taskName}
                            logFilePath={task.status.logFilePath}
                            state={task.status.state}
                            openLogModal={openLogModal}
                            closeLogModal={closeLogModal}
                          />
                          {task.status.state === "success" ? 
                          <div>Last Success Time: {new Date(task.status.lastSuccessTime.split(".")[0]).toLocaleString()}
                          <div>Time Taken: {task.status.timeTaken}</div>
                          </div>
                          :
                          task.status.state === "error" ? 
                          <div>Last Error Time: {new Date(task.status.lastErrorTime.split(".")[0]).toLocaleString()}
                          <div>Time Taken: {task.status.timeTaken}</div>
                          </div>
                          :
                          <div>Last Executed Time: {new Date(task.status.lastErrorTime.split(".")[0]).toLocaleString()}
                          <div>Time Taken: {task.status.timeTaken}</div>
                          </div>}
                        </ul>
                      ))}
                    </AccordionItem>
                  </div>
                );
              })}
            </Accordion>
            <br />
            {steps.length>0 ?
            <div className="ButtonRow">
              <Button onClick={cancel} kind="secondary">
                Cancel
              </Button>
            </div>
            :null
            }
          </>
        )}
        <LogModal
          open={logModalIsOpen}
          log={logs}
          onRequestClose={closeLogModal}
        />
      </Tile>
    </div>
  );
}
