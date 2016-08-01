var React = require('react');
var ReactDOM = require('react-dom');

var Introduction = React.createClass({
  render: function() {
    return (
      <div className="introduction">
        <p>Placeholder for introduction text / documentation and a GIF for how to plugin and setup a scout.</p>
      </div>
    )
  }
})

var Location = React.createClass({
  render: function() {
    return (
      <div className="location">
        <header className="location-header">
          <h2 className="location-title">{this.props.location.name}</h2>
          <p className="location-meta"><a href="edit">[<i className="fa fa-pencil"></i> edit</a>]</p>
        </header>
        <div id="location-details">
          <h3>0 VISITORS</h3>
          <img className="pure-img" alt='test' src='img/off-frame.gif'/>
          <p className="location-meta"><a href="activate">[<i className="fa fa-power-off"></i> activate</a>]</p>
        </div>
      </div>
    );
  }
});

var NavItem = React.createClass({
  render: function() {
    return (
      <li className="navItem"><a href="#">[{this.props.name}]</a></li>
    );
  }
});

var NavList = React.createClass({
  render: function() {
    var navNodes = this.props.data.map(function(location) {
      return (
        <NavItem name={location.name} key={location.id} />
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
    var Httpreq = new XMLHttpRequest(); // a new request
    Httpreq.open("GET", "http://"+window.location.host+"/scouts", true);
    Httpreq.send(null);

    Httpreq.onreadystatechange = function() {
      if (Httpreq.readyState == 4 && Httpreq.status == 200) {
        this.setState({data: JSON.parse(Httpreq.responseText)});
      }
    }.bind(this);
  },

  getInitialState: function() {
    return {data: []};
  },

  componentDidMount: function() {
    this.loadFromServer();
  },

  render: function() {
    var mainContent = <Introduction />
    if (this.state.data.length) {
      mainContent = <Location location={this.state.data[0]} />
    }

    return (
      <div className="application">
        <div className="sidebar pure-u-1 pure-u-md-1-4">
          <div className="header">
            <h1 className="brand"><img className="pure-img" alt='Measure the Future logo' src='/img/logo.gif' /></h1>
            <nav className="nav"><NavList data={this.state.data} /></nav>
          </div>
        </div>
        <div className="content pure-u-1 pure-u-md-3-4" id="results">
          {mainContent}
        </div>
      </div>
    )
  }
})

ReactDOM.render(
  <Application />,
  document.getElementById('layout')
);