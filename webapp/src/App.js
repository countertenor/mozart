import React from 'react';
import {
  HashRouter as Router,
  Switch,
  Route,
} from 'react-router-dom';
// Need to use hashrouter to serve as static html file

import './App.scss';
import Configuration from './components/Configuration/Configuration';

export default function App() {

  return (
    <Router>
      <div>
        <Switch>
            {/* <Route path="/install">
              <Install notificationDispatch={notificationDispatch} />
            </Route> */}
            <Route path="/">
              <Configuration />
            </Route>
          </Switch>
      </div>
    </Router>
  );
}
