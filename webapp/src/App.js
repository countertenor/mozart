import React, {useReducer} from 'react';
import {
  HashRouter as Router,
  Switch,
  Route,
} from 'react-router-dom';
// Need to use hashrouter to serve as static html file

import './App.scss';
import Configuration from './components/Configuration/Configuration';
import Execution from './components/Execution/Execution';
import Install from './components/Status/Install';
import { NOTIFICATION_INIT, NOTIFICATION_REDUCER } from './constants/reducer-constants';


export default function App() {
  const [notification, notificationDispatch] = useReducer(NOTIFICATION_REDUCER, NOTIFICATION_INIT);


  return (
    <Router>
      <div>
        <Switch>
          {/* <Route path="/install">
              <Install notificationDispatch={notificationDispatch} />
            </Route> */}
          <Route path="/status">
            <Execution />
          </Route>
          <Route path="/">
            <Execution />
          </Route>
        </Switch>
      </div>
    </Router>
  );
}
