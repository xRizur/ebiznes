import React from 'react';
import { Link } from 'react-router-dom';

const Navbar = ({ cartItemsCount }) => {
  return (
    <header className="App-header">
      <nav>
        <ul>
          <li>
            <Link to="/">Produkty</Link>
          </li>
          <li>
            <Link to="/cart">Koszyk ({cartItemsCount})</Link>
          </li>
        </ul>
      </nav>
    </header>
  );
};

export default Navbar;
