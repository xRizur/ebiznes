import React from 'react';
import { Link } from 'react-router-dom';

const Cart = ({ cart, setCart }) => {
  const removeFromCart = (itemId) => {
    setCart(cart.filter(item => item.productId !== itemId));
  };

  const updateQuantity = (productId, newQuantity) => {
    if (newQuantity <= 0) {
      removeFromCart(productId);
      return;
    }

    setCart(cart.map(item => 
      item.productId === productId 
        ? { ...item, quantity: newQuantity } 
        : item
    ));
  };

  const totalAmount = cart.reduce((total, item) => total + (item.product.price * item.quantity), 0);

  if (cart.length === 0) {
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
          {cart.map((item) => (
            <tr key={item.productId}>
              <td>{item.product.name}</td>
              <td>{item.product.price.toFixed(2)} zł</td>
              <td>
                <button onClick={() => updateQuantity(item.productId, item.quantity - 1)}>-</button>
                {item.quantity}
                <button onClick={() => updateQuantity(item.productId, item.quantity + 1)}>+</button>
              </td>
              <td>{(item.product.price * item.quantity).toFixed(2)} zł</td>
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

export default Cart;
