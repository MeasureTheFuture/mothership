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
"use strict;"
var React = require('react');

function updateLocations(store, locations, id) {
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

var DeactivateAction = React.createClass({
  handleDeactivate: function() {
  	const { store } = this.context;

    var state = store.getState();
    var location = Object.assign({}, state.locations[state.active]);

    location.authorised = false;
    state.locations[state.active] = location;

    updateLocations(store, state.locations, state.active)
  },

  render: function() {
    return (
      <a href="#" className="warning" onClick={this.handleDeactivate}>[<i className="fa fa-power-off"></i> deactivate]</a>
    )
  }
});
DeactivateAction.contextTypes = {
	store: React.PropTypes.object
};

var ActivateAction = React.createClass({
  handleActivate: function() {
  	const { store } = this.context;

    var state = store.getState();
    var location = Object.assign({}, state.locations[state.active]);

    location.authorised = true;
    state.locations[state.active] = location;

    updateLocations(store, state.locations, state.active)
  },

  render: function() {
    return (
      <a href="#" onClick={this.handleActivate}>[<i className="fa fa-power-off"></i> activate]</a>
    );
  }
});
ActivateAction.contextTypes = {
	store: React.PropTypes.object
};

var MeasureAction = React.createClass({
  handleMeasure: function() {
  	const { store } = this.context;

    var state = store.getState();
    var location = Object.assign({}, state.locations[state.active]);

    location.state = 'measuring';
    state.locations[state.active] = location;

    updateLocations(store, state.locations, state.active);
  },

  render: function() {
    return (
      <a href="#" onClick={this.handleMeasure}>[<i className="fa fa-line-chart"></i> measure]</a>
    );
  }
});
MeasureAction.contextTypes = {
	store: React.PropTypes.object
};

var CalibrateAction = React.createClass({
  handleCalibrate: function() {
  	const { store } = this.context;

    var state = store.getState();
    var location = Object.assign({}, state.locations[state.active]);

    location.state = 'calibrating';
    state.locations[state.active] = location;

    updateLocations(store, state.locations, state.active);
  },

  render: function() {
  	const { store } = this.context;

    var state = store.getState();
    var location = state.locations[state.active];
    var label = ((location.state == 'idle') ? "calibrate" : "recalibrate");

    return (
      <a href="#" onClick={this.handleCalibrate}>[<i className="fa fa-wrench"></i> {label}]</a>
    );
  }
});
CalibrateAction.contextTypes = {
	store: React.PropTypes.object
};

export default class PrimaryActions extends React.Component {
  render() {
  	const { store } = this.context;

    var state = store.getState();
    var location = state.locations[state.active];

    var onOff = (location.authorised ? <DeactivateAction /> : <ActivateAction />);
    var calibrate = ((location.authorised && (location.state == 'idle' || location.state == 'calibrated')) ? <CalibrateAction /> : "");
    var measure = ((location.authorised && location.state == 'calibrated') ? <MeasureAction /> : "");

    return (
      <div className="location-meta">{onOff} {calibrate} {measure}</div>
    );
  }
}
PrimaryActions.contextTypes = {
	store: React.PropTypes.object
};