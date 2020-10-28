// In renderer process (web page).
import React, { useCallback, useEffect, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import { useHistory, useParams } from 'react-router-dom';
// import { ipcRenderer, shell } from 'electron';
import {
  Button, CodeSnippet, InlineLoading, Loading, Modal, Tile
} from 'carbon-components-react';
import { CheckmarkFilled16, Misuse16, View16 } from '@carbon/icons-react';
import styles from './Install.module.scss';
import { getStatus, setupWS } from './InstallUtils';
// import { getProvisioningStatus, restartProvision, restartUninstall } from '../../utils/installer-utils';

function LogModal({ onRequestClose, open, log }) {
  return (
    <Modal
      onRequestClose={onRequestClose}
      passiveModal
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
  duration, taskName, state, openLogModal
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
        <p>{taskName}</p>
        {duration && <p>{`Completion time: ${duration} minutes.`}</p>}
      </div>
      {statusIcon}
      <Button
        renderIcon={View16}
        hasIconOnly
        kind="ghost"
        iconDescription="View logs"
        tooltipPosition="top"
        onClick={() => openLogModal(taskName)}
      />
    </li>
  );
}

Task.propTypes = {
  duration: PropTypes.number.isRequired,
  taskName: PropTypes.string.isRequired,
  state: PropTypes.string.isRequired,
  openLogModal: PropTypes.func.isRequired
};

export default function Install({ notificationDispatch, uninstall }) {
  const intervalRef = useRef(null);
  const logRef = useRef('');
  const history = useHistory();
  const { provider } = useParams();
  const [steps, setSteps] = useState([]);
  const [curStatus, setCurStatus] = useState('loading');
  const [logModalIsOpen, setLogModalIsOpen] = useState(false);
  const [logSocket, setLogSocket] = useState(null);

  const openLogModal = useCallback((logFilePath) => {
    // TODO: Change the path to passed in
    const socket = setupWS('/Users/aaronleong/Desktop/projects/yellowstone_webapp/test_log.txt', (err, streamData) => {
      if (err) {
        return console.log('Handle error TODO');
      }

      logRef.current = logRef.current.concat(`${streamData}\n`);
      setLogModalIsOpen(true);
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
    getStatus((err, data) => {
      if (err) {
        notificationDispatch({
          type: 'setNotification',
          notification: {
            title: 'Error',
            subtitle: 'Unable to fetch the current status, will attemp to retry.',
            caption: '',
            kind: 'error',
          }
        });
      }
      setSteps(data.steps);
      setCurStatus(''); // TODO: Figure out
    });
    // switch (resCurStatus) {
    // case 'complete':
    //   if (uninstall) {
    //     history.push('/uninstallcomplete');
    //   } else {
    //     history.push('/complete'); // We will let the cleanup function handle the interval (see useEffect(() => {}, []))
    //   }
    //   break;
    // case 'failed':
    //   notificationDispatch({
    //     type: 'setNotification',
    //     notification: {
    //       title: 'Error',
    //       subtitle: `${uninstall ? 'Uninstall' : 'Install'} failed. An error has occured while ${uninstall ? 'uninstalling' : 'installing'}. You can either restart or go back to re-enter your configuration. Err: ${response.data.errMsg}.`,
    //       caption: '',
    //       kind: 'error',
    //     }
    //   });
    //   clearStatusInterval();
    //   break;
    // case 'waiting':
    //   notificationDispatch({
    //     type: 'setNotification',
    //     notification: {
    //       title: 'Progress stalled',
    //       subtitle: 'A step has been completed that requires the user to restart the script to continue.',
    //       caption: '',
    //       kind: 'info',
    //     }
    //   });
    //   clearStatusInterval();
    //   break;
    // default:
    //   break;
    // }
    // setStatuses(response.data.stepItems);
    // setCurStatus(response.data.status.toLowerCase());
  }, [uninstall, clearStatusInterval, history, notificationDispatch]);

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

  // const restart = useCallback(async () => {
  //   setCurStatus('loading');
  //   if (uninstall) {
  //     await restartUninstall(provider); // TODO: error handling for everything
  //   } else {
  //     await restartProvision(provider); // TODO: error handling for everything
  //   }
  //   notificationDispatch({
  //     type: 'setNotification',
  //     notification: {
  //       title: 'Success',
  //       subtitle: `${uninstall ? 'Uninstall' : 'Installation'} has restarted successfully.`,
  //       caption: '',
  //       kind: 'success',
  //       settimeout: 3
  //     }
  //   });
  //   getData();
  //   setStatusInterval(); // this handles clearing for us
  // }, [uninstall, setStatusInterval, getData, notificationDispatch, provider]);

  useEffect(() => {
    getData();
    setStatusInterval();
    return clearStatusInterval;
  }, []); // eslint-disable-line

  const header = (
    <h1>
      Deploying configuration
    </h1>
  );

  return (
    <div style={{ marginLeft: "10%" }}>
      <Tile className={styles.Install}>
        {curStatus === "loading" ? (
          <>
            {header}
            <Loading active small={false} />
          </>
        ) : (
          <>
            {header}
            {steps.map((module) => {
              return (
                <div key={module.module} className={styles.module}>
                  <h3>{module.module}</h3>
                  {module.tasks.map((task) => (
                    <ul className={styles.loadersHolder} key={task.taskName}>
                      <Task
                        duration="TODO"
                        taskName={task.taskName}
                        state={task.status.state}
                        openLogModal={openLogModal}
                        closeLogModal={closeLogModal}
                      />
                    </ul>
                  ))}
                </div>
              );
            })}
            <br />
            {/* (curStatus === 'failed' || curStatus === 'waiting') && (
            <div id="install-restartPromptHolder" className={styles.restartPromptHolder}>
              {curStatus === 'failed' ? (
                <>
                  <p>{`There was an error ${uninstall ? 'when uninstalling' : 'in the installation'}.`}</p>
                  <p>You can view this issue in detail in the log file: nz-cloud/nz_install_[timestamp].log.</p>
                </>
              ) : (
                <p>{`The ${uninstall ? 'uninstall' : 'installer'} requires the user to continue the process.`}</p>
              )}
              <div className="ButtonRow">
                <Button
                  id="install-goback-button"
                  onClick={() => history.goBack()}
                  kind="secondary"
                >
                  Back
                </Button>
                <Button
                  id="install-restart-button"
                  onClick={restart}
                >
                  {curStatus === 'failed' ? 'Restart' : 'Continue'}
                </Button>
              </div>
            </div>
              ) */}
            <div className="ButtonRow">
              <Button onClick={history.goBack} kind="secondary">
                Cancel
              </Button>
              <Button onClick={() => history.push("/enclosureConfig")}>
                Next
              </Button>
            </div>
          </>
        )}
        <LogModal
          open={logModalIsOpen}
          log={logRef.current}
          onRequestClose={closeLogModal}
        />
      </Tile>
    </div>
  );
}

Install.defaultProps = {
  uninstall: false
};

Install.propTypes = {
  uninstall: PropTypes.bool
};
