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
import { UpdateActiveLocation, ActiveLocation } from '../reducers/index.js';

var DeactivateAction = React.createClass({
  handleDeactivate: function() {
    const { store } = this.context;
    UpdateActiveLocation(store, "state", 'idle');
    UpdateActiveLocation(store, "authorised", false);
  },

  render: function() {
    return (
      <a id="deactivate" href="#" className="warning" onClick={this.handleDeactivate}>[<i className="fa fa-power-off"></i> deactivate]</a>
    )
  }
});
DeactivateAction.contextTypes = {
  store: React.PropTypes.object
};

var ActivateAction = React.createClass({
  handleActivate: function() {
    const { store } = this.context;
    UpdateActiveLocation(store, "authorised", true);
  },

  render: function() {
    return (
      <a id="activate" href="#" onClick={this.handleActivate}>[<i className="fa fa-power-off"></i> activate]</a>
    );
  }
});
ActivateAction.contextTypes = {
  store: React.PropTypes.object
};

var MeasureAction = React.createClass({
  handleMeasure: function() {
    const { store } = this.context;
    UpdateActiveLocation(store, "state", 'measuring');
  },

  render: function() {
    return (
      <a id="measure" href="#" onClick={this.handleMeasure}>[<i className="fa fa-line-chart"></i> measure]</a>
    );
  }
});
MeasureAction.contextTypes = {
  store: React.PropTypes.object
};

var CalibrateAction = React.createClass({
  handleCalibrate: function() {
    const { store } = this.context;
    UpdateActiveLocation(store, "state", 'calibrating');
  },

  render: function() {
    const { store } = this.context;
    var label = ((ActiveLocation(store).state == 'idle') ? "calibrate" : "recalibrate");

    return (
      <a id="calibrate" href="#" onClick={this.handleCalibrate}>[<i className="fa fa-wrench"></i> {label}]</a>
    );
  }
});
CalibrateAction.contextTypes = {
  store: React.PropTypes.object
};

var EditAction = React.createClass({
  render: function() {
    var txt = (this.props.locationEdit ? "save" : "edit");
    var icon = (this.props.locationEdit ? "fa fa-save" : "fa fa-pencil");

    return (
      <a href="#" onClick={this.props.editCallBack}>[<i className={icon}></i> {txt}]</a>
    );
  }
});
EditAction.contextTypes = {
  store: React.PropTypes.object
};

var PrimaryActions = React.createClass({
  render: function() {
    const { store } = this.context;
    var onOff = (ActiveLocation(store).authorised ? <DeactivateAction /> : <ActivateAction />);
    var calibrate = ((ActiveLocation(store).authorised && (ActiveLocation(store).state == 'idle' || ActiveLocation(store).state == 'calibrated')) ? <CalibrateAction /> : "");
    var measure = ((ActiveLocation(store).authorised && ActiveLocation(store).state == 'calibrated') ? <MeasureAction /> : "");

    return (
      <p className="location-meta">{onOff} {calibrate} {measure}<EditAction locationEdit={this.props.locationEdit} editCallBack={this.props.editCallBack}/></p>
    );
  }
});
PrimaryActions.contextTypes = {
  store: React.PropTypes.object
};

export default PrimaryActions;