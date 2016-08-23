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

import React from 'react';
import reducers from '../reducers/index.js'

var DeactivateAction = React.createClass({
  handleDeactivate: function() {
    const { store } = this.context;
    reducers.UpdateActiveLocation(store, "authorised", false);
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
    reducers.UpdateActiveLocation(store, "authorised", true);
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
    reducers.UpdateActiveLocation(store, "state", 'measuring');
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
    reducers.UpdateActiveLocation(store, "state", 'calibrating');
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

var PrimaryActions = React.createClass({
  render: function() {
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
});
PrimaryActions.contextTypes = {
  store: React.PropTypes.object
};

export default PrimaryActions;