import React from 'react';
import {
  HashRouter as Router,
  Switch,
  Route,
} from 'react-router-dom';
// Need to use hashrouter to serve as static html file

import './App.scss';
import TabView from './components/TabView/TabView';

export default function App() {

  return (
    <Router>
      <div>
        <Switch>
          <Route path="/">
            <TabView/>
          </Route>
        </Switch>
      </div>
    </Router>
  );
}
