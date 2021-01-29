import axios from 'axios';
import { STATUS_DUMMY } from './constants/MockData';

const USE_MOCK = false;

export function getStatus(moduleName, callback) { // eslint-disable-line
  if (USE_MOCK) {
    callback(null, STATUS_DUMMY);
  } else {
    axios.get(`http://localhost:8080/api/v1/state`, {params:{"module":moduleName}}).then((res) => {
      callback(null, res.data);
    }).catch((err) => {
      console.log(err);
      callback('Error getting status');
    });
  }
}

export function setupWS(fp, callback) { // Expect this to be called multiple times
  if (window.WebSocket) {
    const conn = new WebSocket(`ws://127.0.0.1:8080/api/v1/ws?filename=${fp}`);
    conn.onmessage = (evt) => {
      console.log(evt.data);
      callback(null, evt.data);
    };
    return conn; // Return in order to handle close
    // conn.onclose = (evt) => {
    //   callback(null, { event: 'closed'})
    //   var item = document.createElement("div");
    //   item.innerHTML = "<b>Connection closed.</b>";
    //   appendLog(item);
    // };
  }

  callback('Websockets must be supported by your browser in order to view logs.');
}
