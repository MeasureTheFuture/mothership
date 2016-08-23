/*
 * Copyright (C) 2016 Clinton Freeman
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
"use strict";

import React from 'react';
import ReactDOM from 'react-dom';
import { createStore } from 'redux';
import { Provider } from 'react-redux';
import reducers from './reducers/index.js'
import PrimaryActions from './components/primaryActions.jsx';

const s = createStore(reducers.Mothership);

var Introduction = React.createClass({
  render: function() {
    return (
      <div className="introduction">
        <p>Placeholder for introduction text / documentation and a GIF for how to plugin and setup a scout.</p>
      </div>
    )
  }
})

var LocationEdit = React.createClass({
  render: function() {
    const { store } = this.context;

    var state = store.getState();
    var location = state.locations[state.active];

    return (
      <header className="locationLabel">
        <form className="pure-form">
          <h2 className="location-title"><input id="locationInput" className="location-title" type="text" placeholder={location.name} /></h2>
          <p className="location-meta"><a href="#" onClick={this.props.callBack}>[<i className="fa fa-save"></i> save</a>]</p>
        </form>
      </header>
    )
  }
});
LocationEdit.contextTypes = {
  store: React.PropTypes.object
}

var LocationLabel = React.createClass({
  render: function() {
    const { store } = this.context;

    var state = store.getState();
    var location = state.locations[state.active];

    return (
      <header className="locationLabel">
          <h2 className="location-title">{location.name}</h2>
          <p className="location-meta"><a href="#" onClick={this.props.callBack}>[<i className="fa fa-pencil"></i> edit</a>]</p>
      </header>
    )
  }
});
LocationLabel.contextTypes = {
  store: React.PropTypes.object
}

var Location = React.createClass({
  getInitialState: function() {
    return {locationEdit: false};
  },

  saveCallBack: function() {
    const { store } = this.context;
    reducers.UpdateActiveLocation(store, "name", document.getElementById('locationInput').value);
    this.setState({locationEdit: false})
    render();
  },

  editCallBack: function() {
    console.log("EDIT: ");

    this.setState({locationEdit: true});
    render();
  },

  getFrameURL: function() {
    const { store } = this.context;

    var state = store.getState();
    var location = state.locations[state.active];

    if (!location.authorised) {
      return 'img/off-frame.gif';
    }

    if (location.state == 'measuring' || location.state == 'calibrated') {
      return 'scouts/'+location.id+'/frame.jpg?d=' + new Date().getTime();
    } else if (location.state == 'calibrating') {
      return 'img/calibrating-frame.gif';
    }

    return 'img/uncalibrated-frame.gif';
  },

  render: function() {
    const { store } = this.context;
    var state = store.getState();
    var location = state.locations[state.active];

    var locationName = (this.state.locationEdit ? <LocationEdit callBack={this.saveCallBack} /> : <LocationLabel callBack={this.editCallBack} /> )

    return (
      <div className="location">
        <div id="locationName">{ locationName }</div>

        <div id="location-details">
          <h3>0 VISITORS</h3>
          <img className="pure-img" alt='test' src={this.getFrameURL()}/>
          <div className="location-meta"><PrimaryActions /></div>
        </div>
      </div>
    );
  }
});
Location.contextTypes = {
  store: React.PropTypes.object
};

var NavItem = React.createClass({
  handleClick: function() {
    store.dispatch({ type:'SET_ACTIVE', active:this.props.idx})
  },

  render: function() {
    return (
      <li className="navItem">
        <a href="#" onClick={this.handleClick}>[{this.props.name}]</a>
      </li>
    );
  }
});

var NavList = React.createClass({
  render: function() {
    var navNodes = this.props.data.map(function(location, index) {
      return (
        <NavItem name={location.name} key={index} idx={index} />
      )
    });

    return (
      <ul className="navList">
        {navNodes}
      </ul>
    )
  }
})

var Application = React.createClass({
  loadFromServer: function () {
    const { store } = this.context;
    reducers.GetLocations(store);
  },

  componentDidMount: function() {
    this.loadFromServer();
    setInterval(this.loadFromServer, 1000);
  },

  render: function() {
    const { store } = this.context;

    var state = store.getState();
    var mainContent = ((state.locations.length) ? <Location /> : <Introduction />);

    return (
      <div className="pure-g">
        <div className="sidebar pure-u-1 pure-u-md-1-4">
          <div className="header">
            <h1 className="brand"><img className="pure-img" alt='Measure the Future logo' src='/img/logo.gif' /></h1>
            <nav className="nav"><NavList data={state.locations} /></nav>
          </div>
        </div>
        <div className="content pure-u-1 pure-u-md-3-4" id="results">
          {mainContent}
        </div>
      </div>
    )
  }
})
Application.contextTypes = {
  store: React.PropTypes.object
};

function render() {
  ReactDOM.render(
    <Provider store={s}>
      <Application />
    </Provider>,
    document.getElementById('application')
  );
}

render();
s.subscribe(render)