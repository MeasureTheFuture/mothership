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
var React = require('react');
var ReactDOM = require('react-dom');
var Redux = require('redux');
var reducers = require('./reducers');

const store = Redux.createStore(reducers);

function updateLocations(locations, id) {
  var l = locations[id]

  // Push the updated location to the backend.
  var Httpreq = new XMLHttpRequest();
  Httpreq.open("PUT", "http://"+window.location.host+"/scouts/"+l.id, true);
  Httpreq.send(JSON.stringify(l));

  Httpreq.onreadystatechange = function() {
    if (Httpreq.readyState == 4 && Httpreq.status == 200) {
      store.dispatch({ type:'UPDATE_LOCATIONS', locations:locations})
    }
  }
}

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
  handleSave: function() {
    var state = store.getState();
    var location = Object.assign({}, state.locations[state.active]);

    location.name = document.getElementById('locationInput').value;
    console.log("saving: " + location.name);
    state.locations[state.active] = location;

    updateLocations(state.locations, state.active);

    ReactDOM.render(
      <LocationLabel />,
      document.getElementById('locationName')
    )
  },

  render: function() {
    var state = store.getState();
    var location = state.locations[state.active];

    return (
      <header className="locationLabel">
        <form className="pure-form">
          <h2 className="location-title"><input id="locationInput" className="location-title" type="text" placeholder={location.name} /></h2>
          <p className="location-meta"><a href="#" onClick={this.handleSave}>[<i className="fa fa-save"></i> save</a>]</p>
        </form>
      </header>
    )
  }
});

var LocationLabel = React.createClass({
  handleEdit: function() {
    ReactDOM.render(
      <LocationEdit />,
      document.getElementById('locationName'))
  },

  render: function() {
    var state = store.getState();
    var location = state.locations[state.active];

    return (
      <header className="locationLabel">
          <h2 className="location-title">{location.name}</h2>
          <p className="location-meta"><a href="#" onClick={this.handleEdit}>[<i className="fa fa-pencil"></i> edit</a>]</p>
      </header>
    )
  }
});

var DeactivateAction = React.createClass({
  handleDeactivate: function() {
    var state = store.getState();
    var location = Object.assign({}, state.locations[state.active]);

    location.authorised = false;
    state.locations[state.active] = location;

    updateLocations(state.locations, state.active)
  },

  render: function() {
    return (
      <a href="#" className="warning" onClick={this.handleDeactivate}>[<i className="fa fa-power-off"></i> deactivate]</a>
    )
  }
})

var ActivateAction = React.createClass({
  handleActivate: function() {
    var state = store.getState();
    var location = Object.assign({}, state.locations[state.active]);

    location.authorised = true;
    state.locations[state.active] = location;

    updateLocations(state.locations, state.active)
  },

  render: function() {
    return (
      <a href="#" onClick={this.handleActivate}>[<i className="fa fa-power-off"></i> activate]</a>
    );
  }
})

var MeasureAction = React.createClass({
  handleMeasure: function() {
    var state = store.getState();
    var location = Object.assign({}, state.locations[state.active]);

    location.state = 'measuring';
    state.locations[state.active] = location;

    updateLocations(state.locations, state.active);
  },

  render: function() {
    return (
      <a href="#" onClick={this.handleMeasure}>[<i className="fa fa-line-chart"></i> measure]</a>
    );
  }
})

var CalibrateAction = React.createClass({
  handleCalibrate: function() {
    var state = store.getState();
    var location = Object.assign({}, state.locations[state.active]);

    location.state = 'calibrating';
    state.locations[state.active] = location;

    updateLocations(state.locations, state.active);
  },

  render: function() {
    var state = store.getState();
    var location = state.locations[state.active];
    var label = ((location.state == 'idle') ? "calibrate" : "recalibrate");

    return (
      <a href="#" onClick={this.handleCalibrate}>[<i className="fa fa-wrench"></i> {label}]</a>
    );
  }
})

var PrimaryActions = React.createClass({
  render: function() {
    var state = store.getState();
    var location = state.locations[state.active];

    var onOff = (location.authorised ? <DeactivateAction /> : <ActivateAction />);
    var calibrate = ((location.authorised && (location.state == 'idle' || location.state == 'calibrated')) ? <CalibrateAction /> : "");
    var measure = ((location.authorised && location.state == 'calibrated') ? <MeasureAction /> : "");

    return (
      <p className="location-meta">{onOff} {calibrate}</p>
    );
  }
})

var Location = React.createClass({
  getFrameURL: function() {
    var state = store.getState();
    var location = state.locations[state.active];

    if (!location.authorised) {
      return 'img/off-frame.gif';
    }

    if (location.state == 'measuring' || location.state == 'calibrated') {
      return 'scouts/'+location.id+'/frame.jpg';
    } else if (location.state == 'calibrating') {
      return 'img/calibrating-frame.gif';
    }

    return 'img/uncalibrated-frame.gif';
  },

  render: function() {
    var state = store.getState();
    var location = state.locations[state.active];

    return (
      <div className="location">
        <div id="locationName"><LocationLabel /></div>

        <div id="location-details">
          <h3>0 VISITORS</h3>
          <img className="pure-img" alt='test' src={this.getFrameURL()}/>
          <p className="location-meta"><PrimaryActions /></p>
        </div>
      </div>
    );
  }
});

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
    var Httpreq = new XMLHttpRequest();
    Httpreq.open("GET", "http://"+window.location.host+"/scouts", true);
    Httpreq.send(null);

    Httpreq.onreadystatechange = function() {
      if (Httpreq.readyState == 4 && Httpreq.status == 200) {
        var locations = JSON.parse(Httpreq.responseText)
        store.dispatch({ type:'UPDATE_LOCATIONS', locations:locations})
      }
    }.bind(this);
  },

  componentDidMount: function() {
    this.loadFromServer();
  },

  render: function() {
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

function render() {
  ReactDOM.render(
    <Application />,
    document.getElementById('application')
  );
}

render();
store.subscribe(render)