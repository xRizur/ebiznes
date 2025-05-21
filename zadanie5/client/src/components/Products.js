import React from 'react';
import PropTypes from 'prop-types';

const Products = ({ products, loading, addToCart }) => {
  if (loading) {
    return <div>Ładowanie produktów...</div>;
  }

  return (
    <div>
      <h1>Produkty</h1>
      <div className="product-list">
        {products.map((product) => (
          <div key={product.id} className="product-card">
            <img src={product.imageUrl} alt={product.name} style={{ maxWidth: '100%' }} />
            <h3>{product.name}</h3>
            <p className="description">{product.description}</p>
            <p className="price">Cena: {product.price.toFixed(2)} zł</p>
            <button onClick={() => addToCart(product)}>Dodaj do koszyka</button>
          </div>
        ))}
      </div>
    </div>
  );
};

Products.propTypes = {
  products: PropTypes.arrayOf(PropTypes.shape({
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
    description: PropTypes.string.isRequired,
    price: PropTypes.number.isRequired,
    imageUrl: PropTypes.string.isRequired
  })),
  loading: PropTypes.bool,
  addToCart: PropTypes.func.isRequired
};

Products.defaultProps = {
  products: [],
  loading: false
};

export default Products;
