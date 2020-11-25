import React, { useCallback, useEffect, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import {
  Button, CodeSnippet, InlineLoading, Loading, Modal, Tile, Accordion, AccordionItem
} from 'carbon-components-react';
import { CheckmarkFilled16, Misuse16, View16 } from '@carbon/icons-react';
import styles from './Install.module.scss';
import { getStatus, setupWS } from './InstallUtils';
import _ from 'lodash';

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
    statusIcon = <div className={styles.loaderIconHolder}><InlineLoading /></div>;
    break;

  case 'failed':
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
  case 'notstarted':
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
  for(let i=0; i<35; i++){
    stringPrint+=" - "
  }
  const [print] = useState(stringPrint)
  const [percentage, setPercentage] = useState({});
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
    // let moduleName = props.getModuleName()
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
        let count = 0;
        e.tasks.map(item =>{
          if(item.status.state === "success"){
            countObj = {...countObj,
              [e.directory]: Math.trunc(++count * 100/e.tasks.length)
            }
          }
        })
      })
      setPercentage(countObj);
    });
  }, [clearStatusInterval]);


  const cancel = (e) =>{
    console.log("Cancel pressed!", props)
    props.switchActiveTab("configuration")
    // e.preventDefault();
    // axios
    // .put(`http://localhost:8080/api/v1/cancel`)
    // .then((res) => {
    //   console.log("response cancel??: ", res.data);
    // })
    // .catch((err) => {
    //   console.log("error cancel",err);
    // });
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

  const header = (
    <h6>
      Status: {state} {moduleName}
    </h6>
  );

  // root : test-module
    // child1 : 00-bash-module
      // child11 : 00-module1
      // child12 : 00-module2
      // child13 : 00-module3
    // child2 : 10-python-module
      // child21 : 00-module1
        // child211 : 01-module1
      // child22 : 00-module2
    // child3 : 20-random-module

  // [test-module, 00-bash-module, 00-module1]
  // test-module/00-bash-module/01-module2
  // test-module/00-bash-module/02-module3
  // test-module/10-python-module/00-module1
  // test-module/10-python-module/00-module1/01-module1
  // test-module/10-python-module/10-module2
  // test-module/20-random-module

  // let ob = {
  //   "test-module": {
  //     "00-bash-module": {
  //       "00-module1": {},
  //       "01-module2":{},
  //       "02-module3":{}
  //     },
  //     "10-python-module":{
  //       "00-module1":{
  //         "hasFile":true,
  //         "01-module1":{}
  //       }
  //     }
  //   },
  // };

  let o = {
    "test-module": {
      "00-bash-module": {
        "00-module1": {},
        "01-module2": {},
        "02-module3": {},
      },
      "10-python-module": {
        "00-module1": { 
          "01-module1": {} 
        },
        "10-module2": {},
      },
      "20-random-module": {},
    },
  };

  let nested_obj ={}
  let nestedObject = steps.map((module) => {
    let hierarchy = module.directory.split("/");    //[test-module, 00-bash-module,00-module1], test-module/00-bash-module/01-module2
    let result = hierarchy.reverse().reduce((res, key) => ({[key]: res}), {});
    nested_obj = _.merge(nested_obj, result)
  })
  console.log(JSON.stringify(nested_obj));

  function NestedAccordions ({ k, v }) {
    if (Object.keys(v).length === 0) {
      return  (
        <AccordionItem title={k} k={module} v={v[module]}>
          {[
                {
                    "taskName": "step1",
                    "status": {
                        "startTime": "2020-11-18T17:11:13.816348-08:00",
                        "timeTaken": "3.020846466s",
                        "lastSuccessTime": "2020-11-18 17:11:16.837249 -0800 PST m=+18.470463882",
                        "lastErrorTime": "",
                        "state": "success",
                        "logFilePath": "/Users/tosha.kamath@ibm.com/IBM/new/mozart/logs/2020-11-18--17-11-13.816-00-step1.log"
                    }
                },
                {
                    "taskName": "step2",
                    "status": {
                        "startTime": "2020-11-18T17:11:16.837818-08:00",
                        "timeTaken": "2.017150379s",
                        "lastSuccessTime": "2020-11-18 17:11:18.855009 -0800 PST m=+20.488177611",
                        "lastErrorTime": "",
                        "state": "success",
                        "logFilePath": "/Users/tosha.kamath@ibm.com/IBM/new/mozart/logs/2020-11-18--17-11-16.837-10-step2.log"
                    }
                }
            ].map((task) => (
            <ul className={styles.loadersHolder} key={task.taskName}>
              <Task
                duration="TODO"
                taskName={task.taskName}
                logFilePath={task.status.logFilePath}
                state={task.status.state}
                openLogModal={openLogModal}
                closeLogModal={closeLogModal}
              />
            </ul>
          ))}
        </AccordionItem>
      );
    } else {
      const nestedAcc = Object.keys(v).map((module) => {
      return <AccordionItem title={module}><NestedAccordions k={module} v={v[module]}/></AccordionItem>;
      });
      return <Accordion size="xl" align="start">{nestedAcc}</Accordion>;
    }
  };

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
              {/* {steps.map((module) => { */}
                {/* return ( */}
                  <div key={module.directory} className={styles.module}>
                    <NestedAccordions k="" v={nested_obj}/>
                    {/* {nestedAccordions({k:null, v:nested_obj})} */}
                        {/* for (const [key, value] of Object.entries(nested_obj)) {
                        console.log(`${key}: ${value}`); */}
                    {/* <AccordionItem title= {module.directory+print+percentage[module.directory]+"% complete"}>
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
                        </ul>
                      ))}
                    </AccordionItem> */}
                  </div>
                {/* ); */}
              {/* })} */}
            </Accordion>
            <br />
            <div className="ButtonRow">
              <Button onClick={cancel} kind="secondary">
                Cancel
              </Button>
            </div>
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
