import React from 'react';
import { Link } from 'react-router-dom';
import PropTypes from 'prop-types';

const Navbar = ({ cartItemsCount }) => {
  return (
    <nav className="navbar">
      <div className="navbar-container">
        <Link to="/" className="navbar-logo">
          Sklep Internetowy
        </Link>
        <ul className="nav-menu">
          <li className="nav-item">
            <Link to="/" className="nav-links">
              Produkty
            </Link>
          </li>
          <li className="nav-item">
            <Link to="/cart" className="nav-links">
              Koszyk
              {cartItemsCount > 0 && (
                <span className="cart-count">{cartItemsCount}</span>
              )}
            </Link>
          </li>
        </ul>
      </div>
    </nav>
  );
};

Navbar.propTypes = {
  cartItemsCount: PropTypes.number
};

Navbar.defaultProps = {
  cartItemsCount: 0
};

export default Navbar;
