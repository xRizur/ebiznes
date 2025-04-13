import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import axios from 'axios';
import Products from './components/Products';
import Cart from './components/Cart';
import Payment from './components/Payment';
import Navbar from './components/Navbar';
import './App.css';

// Konfiguracja Axios
axios.defaults.baseURL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function App() {
  const [products, setProducts] = useState([]);
  const [cart, setCart] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Pobieranie produktów
  useEffect(() => {
    const fetchProducts = async () => {
      setLoading(true);
      try {
        const response = await axios.get('/products');
        setProducts(response.data);
        setError(null);
      } catch (err) {
        console.error('Error fetching products:', err);
        setError('Nie udało się pobrać produktów. Spróbuj ponownie później.');
      } finally {
        setLoading(false);
      }
    };

    fetchProducts();
  }, []);

  // Dodawanie produktu do koszyka
  const addToCart = async (product) => {
    try {
      // Sprawdź czy produkt już jest w koszyku
      const existingItem = cart.find(item => item.productId === product.id);
      
      if (existingItem) {
        // Aktualizuj istniejący element w koszyku
        const updatedCart = cart.map(item =>
          item.productId === product.id 
            ? { ...item, quantity: item.quantity + 1 } 
            : item
        );
        setCart(updatedCart);
      } else {
        // Dodaj nowy element do koszyka
        const newCartItem = {
          productId: product.id,
          product: product,
          quantity: 1
        };
        
        // Zapisz do API
        await axios.post('/cart', newCartItem);
        
        // Zaktualizuj stan lokalny
        setCart([...cart, newCartItem]);
      }
      
      setError(null);
    } catch (err) {
      console.error('Error adding to cart:', err);
      setError('Nie udało się dodać produktu do koszyka.');
    }
  };

  // Przetwarzanie płatności
  const processPayment = async (paymentDetails) => {
    try {
      const response = await axios.post('/payments', {
        ...paymentDetails,
        amount: cart.reduce((total, item) => total + (item.product.price * item.quantity), 0)
      });
      
      // Po udanej płatności - wyczyść koszyk
      if (response.data.status === 'completed') {
        setCart([]);
      }
      
      return response.data;
    } catch (err) {
      console.error('Error processing payment:', err);
      setError('Błąd przetwarzania płatności. Spróbuj ponownie.');
      throw err;
    }
  };

  return (
    <Router>
      <div className="App">
        <Navbar cartItemsCount={cart.reduce((acc, item) => acc + item.quantity, 0)} />
        
        {error && <div className="error-message">{error}</div>}
        
        <Routes>
          <Route path="/" element={
            <Products 
              products={products} 
              loading={loading} 
              addToCart={addToCart} 
            />
          } />
          <Route path="/cart" element={
            <Cart 
              cart={cart}
              setCart={setCart}
            />
          } />
          <Route path="/payment" element={
            <Payment 
              cart={cart} 
              processPayment={processPayment} 
            />
          } />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
