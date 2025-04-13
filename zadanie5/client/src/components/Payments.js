import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useCart } from '../hooks/useCart';

function Payments() {
  const { cart, total, clearCart } = useCart();
  const navigate = useNavigate();
  
  const [formData, setFormData] = useState({
    fullName: '',
    email: '',
    address: '',
    cardNumber: '',
    expiryDate: '',
    cvv: ''
  });
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (cart.length === 0) {
      alert('Your cart is empty!');
      navigate('/');
      return;
    }
    
    setLoading(true);
    setError(null);
    
    try {
      // Prepare order data
      const orderData = {
        customer: {
          fullName: formData.fullName,
          email: formData.email,
          address: formData.address
        },
        payment: {
          cardNumber: formData.cardNumber,
          expiryDate: formData.expiryDate,
          cvv: formData.cvv
        },
        items: cart.map(item => ({
          productId: item.id,
          quantity: item.quantity,
          price: item.price
        })),
        total: total
      };
      
      // Send to server with CORS headers
      await axios.post('http://localhost:8080/orders', orderData, {
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      });
      
      alert('Your order has been placed successfully!');
      clearCart();
      navigate('/');
      
    } catch (err) {
      setError('Payment processing failed. Please try again.');
      console.error('Payment error:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2>Payment Information</h2>
      {cart.length === 0 ? (
        <div>
          <p>Your cart is empty. Please add items before proceeding to payment.</p>
          <button onClick={() => navigate('/')}>Browse Products</button>
        </div>
      ) : (
        <>
          <div className="order-summary">
            <h3>Order Summary</h3>
            <p>Total items: {cart.reduce((acc, item) => acc + item.quantity, 0)}</p>
            <p>Total amount: ${total.toFixed(2)}</p>
          </div>
          
          <form onSubmit={handleSubmit} className="payment-form">
            <h3>Customer Information</h3>
            <div>
              <label htmlFor="fullName">Full Name</label>
              <input
                type="text"
                id="fullName"
                name="fullName"
                value={formData.fullName}
                onChange={handleChange}
                required
              />
            </div>
            <div>
              <label htmlFor="email">Email</label>
              <input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                required
              />
            </div>
            <div>
              <label htmlFor="address">Shipping Address</label>
              <input
                type="text"
                id="address"
                name="address"
                value={formData.address}
                onChange={handleChange}
                required
              />
            </div>
            
            <h3>Payment Details</h3>
            <div>
              <label htmlFor="cardNumber">Card Number</label>
              <input
                type="text"
                id="cardNumber"
                name="cardNumber"
                value={formData.cardNumber}
                onChange={handleChange}
                placeholder="XXXX XXXX XXXX XXXX"
                required
              />
            </div>
            <div>
              <label htmlFor="expiryDate">Expiry Date</label>
              <input
                type="text"
                id="expiryDate"
                name="expiryDate"
                value={formData.expiryDate}
                onChange={handleChange}
                placeholder="MM/YY"
                required
              />
            </div>
            <div>
              <label htmlFor="cvv">CVV</label>
              <input
                type="text"
                id="cvv"
                name="cvv"
                value={formData.cvv}
                onChange={handleChange}
                placeholder="123"
                required
              />
            </div>
            
            {error && <p className="error">{error}</p>}
            
            <div style={{ marginTop: '20px' }}>
              <button type="button" onClick={() => navigate('/cart')} style={{ marginRight: '10px', backgroundColor: '#666' }}>
                Back to Cart
              </button>
              <button type="submit" disabled={loading}>
                {loading ? 'Processing...' : 'Place Order'}
              </button>
            </div>
          </form>
        </>
      )}
    </div>
  );
}

export default Payments;
