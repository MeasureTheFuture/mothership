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
  toI: function(v) {
    return v | 0;
  },

  lerp: function(l, r, t) {
    return l + (r - l) * t
  },

  generateFill: function(t) {
    if (t < 0.001) {
      return "rgba(0, 0, 0, 0)"
    }

    if (t < 0.5) {
      return "rgba("+this.toI(this.lerp(19, 250, t))+","
        +this.toI(this.lerp(27, 212, t))+","
        +this.toI(this.lerp(66, 12, t))+","
        +this.lerp(0.4, 0.5, t)+")"
    }

    return "rgba("+this.toI(this.lerp(250, 186, t))+","
      +this.toI(this.lerp(212, 8, t))+","
      +this.toI(this.lerp(12, 16, t))+","
      +this.lerp(0.5, 0.5, t)+")"
  },

  maxTime: function(buckets) {
    var maxT = 0.0;

    buckets.map(function(i) {
      i.map(function(j) {
        maxT = Math.max(maxT, j);
      })
    })

    return maxT
  },

  render: function() {
    const { store } = this.context;
    var url = 'scouts/'+ActiveLocation(store).id+'/frame.jpg?d=' + new Date().getTime();
    var buckets = ActiveLocation(store).summary.VisitTimeBuckets;
    var w = 1920;
    var h = 1080;
    var iBuckets = buckets.length;
    var jBuckets = buckets[0].length;
    var bucketW = w/iBuckets;
    var bucketH = h/jBuckets;
    var maxT = this.maxTime(buckets);
    var viewBox="0 0 " + w + " " + h;

    var data = []
    for (var i = 0; i < iBuckets; i++) {
      for (var j = 0; j < jBuckets; j++) {
        var t = 0.0;
        if (maxT > 0.0) {
          t = buckets[i][j] / maxT;
        }

        data.push(<rect key={i*iBuckets+j} x={i*bucketW} y={j*bucketH} width={bucketW} height={bucketH} style={{fill:this.generateFill(t)}} />);
        //data.push(<text x={(i*bucketW) + (bucketW / 2)} y={(j*bucketH) + (bucketH / 2)} fontFamily="Verdana" fontSize="20">{t.toFixed(3)}</text>)
      }
    }

    return (
      <svg xmlns="http://www.w3.org/2000/svg" xmlnsXlink="http://www.w3.org/1999/xlink" viewBox={viewBox}>
        <image x="0" y="0" width={w} height={h} xlinkHref={url}/>
        {data}
      </svg>
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