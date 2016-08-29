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
import { UpdateActiveLocation, ActiveLocation } from '../reducers/index.js'
import PrimaryActions from './primaryActions.jsx';

var LocationEdit = React.createClass({
  render: function() {
    const { store } = this.context;

    return (
      <header className="locationLabel">
        <form className="pure-form">
          <h2 className="location-title"><input id="locationInput" className="location-title" type="text" placeholder={ActiveLocation(store).name} /></h2>
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

    return (
      <header className="locationLabel">
          <h2 className="location-title">{ActiveLocation(store).name}</h2>
          <p className="location-meta"><a href="#" onClick={this.props.callBack}>[<i className="fa fa-pencil"></i> edit</a>]</p>
      </header>
    )
  }
});
LocationLabel.contextTypes = {
  store: React.PropTypes.object
}

var Placeholder = React.createClass({
  getFrameURL: function() {
    const { store } = this.context;

    if (!ActiveLocation(store).authorised) {
      return 'img/off-frame.gif';
    }

    if (ActiveLocation(store).state == 'calibrated') {
      return 'scouts/'+ActiveLocation(store).id+'/frame.jpg?d=' + new Date().getTime();
    } else if (ActiveLocation(store).state == 'calibrating') {
      return 'img/calibrating-frame.gif';
    }

    return 'img/uncalibrated-frame.gif';
  },

  render: function() {
    return (
      <img className="pure-img" alt='test' src={this.getFrameURL()}/>
    )
  }
});
Placeholder.contextTypes = {
  store: React.PropTypes.object
}

var Heatmap = React.createClass({
  render: function() {
    const { store } = this.context;
    url = 'scouts/'+ActiveLocation(store).id+'/frame.jpg?d=' + new Date().getTime();

    return (
      <img className="pure-img" alt='test' src={url}/>
    )
  }
});
Heatmap.contextTypes = {
  store: React.PropTypes.object
}

Location = React.createClass({
  getInitialState: function() {
    return {locationEdit: false};
  },

  saveCallBack: function() {
    const { store } = this.context;
    UpdateActiveLocation(store, "name", document.getElementById('locationInput').value);
    this.setState({locationEdit: false})
    this.render();
  },

  editCallBack: function() {
    this.setState({locationEdit: true});
    this.render();
  },

  render: function() {
    const { store } = this.context;
    var locationName = (this.state.locationEdit ? <LocationEdit callBack={this.saveCallBack} /> : <LocationLabel callBack={this.editCallBack} /> )
    var heatmap = (ActiveLocation(store).state == 'measuring') ? <Heatmap /> : <Placeholder />

    return (
      <div className="location">
        <div id="locationName">{ locationName }</div>

        <div id="location-details">
          <h3>{ ActiveLocation(store).summary.VisitorCount } VISITORS</h3>
          { heatmap }
          <div className="location-meta"><PrimaryActions /></div>
        </div>
      </div>
    );
  }
});
Location.contextTypes = {
  store: React.PropTypes.object
};

export default Location;