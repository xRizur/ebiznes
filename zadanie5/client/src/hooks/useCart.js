import React, { createContext, useContext, useReducer } from 'react';

// Create a Context
const CartContext = createContext();

// Initial state
const initialState = {
  items: [],
  total: 0
};

// Actions
const ADD_TO_CART = 'ADD_TO_CART';
const REMOVE_FROM_CART = 'REMOVE_FROM_CART';
const CLEAR_CART = 'CLEAR_CART';

// Reducer function
function cartReducer(state, action) {
  switch (action.type) {
    case ADD_TO_CART: {
      const existingItem = state.items.find(item => item.id === action.payload.id);
      
      if (existingItem) {
        // If item exists, increase quantity
        const updatedItems = state.items.map(item => 
          item.id === action.payload.id ? {...item, quantity: item.quantity + 1} : item
        );
        return {
          ...state,
          items: updatedItems,
          total: state.total + action.payload.price
        };
      } else {
        // Add new item with quantity 1
        return {
          ...state,
          items: [...state.items, {...action.payload, quantity: 1}],
          total: state.total + action.payload.price
        };
      }
    }
    
    case REMOVE_FROM_CART: {
      const itemToRemove = state.items.find(item => item.id === action.payload);
      if (!itemToRemove) return state;

      if (itemToRemove.quantity > 1) {
        // If quantity > 1, decrease quantity
        const updatedItems = state.items.map(item => 
          item.id === action.payload ? {...item, quantity: item.quantity - 1} : item
        );
        return {
          ...state,
          items: updatedItems,
          total: state.total - itemToRemove.price
        };
      } else {
        // If quantity is 1, remove the item completely
        return {
          ...state,
          items: state.items.filter(item => item.id !== action.payload),
          total: state.total - itemToRemove.price
        };
      }
    }
    
    case CLEAR_CART:
      return initialState;
      
    default:
      return state;
  }
}

// Provider component
export const CartProvider = ({ children }) => {
  const [state, dispatch] = useReducer(cartReducer, initialState);
  
  // Action creators
  const addToCart = (product) => {
    dispatch({ type: ADD_TO_CART, payload: product });
  };
  
  const removeFromCart = (productId) => {
    dispatch({ type: REMOVE_FROM_CART, payload: productId });
  };
  
  const clearCart = () => {
    dispatch({ type: CLEAR_CART });
  };

  return (
    <CartContext.Provider value={{ 
      cart: state.items, 
      total: state.total, 
      addToCart, 
      removeFromCart, 
      clearCart 
    }}>
      {children}
    </CartContext.Provider>
  );
};

// Custom hook
export const useCart = () => {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
};
