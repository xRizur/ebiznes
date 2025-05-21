import React from 'react';
import { Link } from 'react-router-dom';
import PropTypes from 'prop-types';

const Cart = ({ cart, setCart }) => {
  const safeCart = cart || [];

  console.log('Cart state in Cart component:', safeCart);
  
  const removeFromCart = (productId) => {
    setCart(safeCart.filter(item => item.productId !== productId));
  };

  const updateQuantity = (productId, newQuantity) => {
    if (newQuantity <= 0) {
      removeFromCart(productId);
      return;
    }

    setCart(safeCart.map(item => 
      item.productId === productId 
        ? { ...item, quantity: newQuantity } 
        : item
    ));
  };

  const totalAmount = safeCart.reduce((total, item) => 
    total + (item.product ? item.product.price * item.quantity : 0), 0);

  if (safeCart.length === 0) {
    return (
      <div>
        <h1>Koszyk</h1>
        <p>Twój koszyk jest pusty.</p>
        <Link to="/">Kontynuuj zakupy</Link>
      </div>
    );
  }

  return (
    <div>
      <h1>Koszyk</h1>
      <table className="cart-table">
        <thead>
          <tr>
            <th>Produkt</th>
            <th>Cena</th>
            <th>Ilość</th>
            <th>Suma</th>
            <th>Akcje</th>
          </tr>
        </thead>
        <tbody>
          {safeCart.map((item) => (
            <tr key={item.productId}>
              <td>{item.product ? item.product.name : 'Produkt niedostępny'}</td>
              <td>{item.product ? item.product.price.toFixed(2) : 0} zł</td>
              <td>
                <button onClick={() => updateQuantity(item.productId, item.quantity - 1)}>-</button>
                <span className="quantity-value">{item.quantity}</span>
                <button onClick={() => updateQuantity(item.productId, item.quantity + 1)}>+</button>
              </td>
              <td>{item.product ? (item.product.price * item.quantity).toFixed(2) : 0} zł</td>
              <td>
                <button onClick={() => removeFromCart(item.productId)}>Usuń</button>
              </td>
            </tr>
          ))}
        </tbody>
        <tfoot>
          <tr>
            <td colSpan="3"><strong>Suma całkowita:</strong></td>
            <td>{totalAmount.toFixed(2)} zł</td>
            <td></td>
          </tr>
        </tfoot>
      </table>
      <div style={{ marginTop: '20px' }}>
        <Link to="/payment">
          <button>Przejdź do płatności</button>
        </Link>
      </div>
    </div>
  );
};

Cart.propTypes = {
  cart: PropTypes.arrayOf(PropTypes.shape({
    productId: PropTypes.number.isRequired,
    product: PropTypes.shape({
      name: PropTypes.string.isRequired,
      price: PropTypes.number.isRequired
    }).isRequired,
    quantity: PropTypes.number.isRequired
  })),
  setCart: PropTypes.func.isRequired
};

Cart.defaultProps = {
  cart: []
};

export default Cart;
