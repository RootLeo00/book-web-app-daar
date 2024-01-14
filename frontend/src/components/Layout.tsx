import React from 'react';
import PropTypes from 'prop-types';
import '../styles/layout.css'; 

const Layout = ({ children }:any) => {
  return (
    <div className="layout">
      {children}
    </div>
  );
};

Layout.propTypes = {
  children: PropTypes.node.isRequired,
};

export default Layout;
