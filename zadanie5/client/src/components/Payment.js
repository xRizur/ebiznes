import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const Payment = ({ cart, processPayment }) => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    cardNumber: '',
    cardHolder: '',
    expiryDate: '',
    cvv: ''
  });
  const [formErrors, setFormErrors] = useState({});
  const [isProcessing, setIsProcessing] = useState(false);

  const totalAmount = cart.reduce((total, item) => total + (item.product.price * item.quantity), 0);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
  };

  const validate = () => {
    const errors = {};
    if (!formData.cardNumber.trim()) {
      errors.cardNumber = 'Numer karty jest wymagany';
    } else if (!/^\d{16}$/.test(formData.cardNumber.replace(/\s/g, ''))) {
      errors.cardNumber = 'Nieprawidłowy format numeru karty';
    }

    if (!formData.cardHolder.trim()) {
      errors.cardHolder = 'Imię i nazwisko jest wymagane';
    }

    if (!formData.expiryDate.trim()) {
      errors.expiryDate = 'Data ważności jest wymagana';
    } else if (!/^(0[1-9]|1[0-2])\/\d{2}$/.test(formData.expiryDate)) {
      errors.expiryDate = 'Nieprawidłowy format daty (MM/YY)';
    }

    if (!formData.cvv.trim()) {
      errors.cvv = 'Kod CVV jest wymagany';
    } else if (!/^\d{3,4}$/.test(formData.cvv)) {
      errors.cvv = 'Nieprawidłowy format CVV';
    }

    return errors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    const errors = validate();
    if (Object.keys(errors).length > 0) {
      setFormErrors(errors);
      return;
    }

    setIsProcessing(true);
    try {
      await processPayment({
        ...formData,
        amount: totalAmount
      });
      
      // Po udanej płatności przekieruj do strony głównej
      alert('Płatność została zrealizowana pomyślnie!');
      navigate('/');
    } catch (error) {
      console.error('Błąd podczas przetwarzania płatności:', error);
      setFormErrors({ submit: 'Błąd przetwarzania płatności. Spróbuj ponownie.' });
    } finally {
      setIsProcessing(false);
    }
  };

  if (cart.length === 0) {
    return (
      <div>
        <h1>Płatność</h1>
        <p>Twój koszyk jest pusty, nie możesz dokonać płatności.</p>
        <button onClick={() => navigate('/')}>Wróć do zakupów</button>
      </div>
    );
  }

  return (
    <div>
      <h1>Płatność</h1>
      <div className="payment-summary">
        <h2>Podsumowanie zamówienia</h2>
        <p>Całkowita kwota: {totalAmount.toFixed(2)} zł</p>
      </div>

      <form className="payment-form" onSubmit={handleSubmit}>
        <div>
          <label>Numer karty:</label>
          <input
            type="text"
            name="cardNumber"
            value={formData.cardNumber}
            onChange={handleChange}
            placeholder="1234 5678 9012 3456"
          />
          {formErrors.cardNumber && <p className="error">{formErrors.cardNumber}</p>}
        </div>

        <div>
          <label>Imię i nazwisko:</label>
          <input
            type="text"
            name="cardHolder"
            value={formData.cardHolder}
            onChange={handleChange}
            placeholder="Jan Kowalski"
          />
          {formErrors.cardHolder && <p className="error">{formErrors.cardHolder}</p>}
        </div>

        <div>
          <label>Data ważności:</label>
          <input
            type="text"
            name="expiryDate"
            value={formData.expiryDate}
            onChange={handleChange}
            placeholder="MM/YY"
          />
          {formErrors.expiryDate && <p className="error">{formErrors.expiryDate}</p>}
        </div>

        <div>
          <label>CVV:</label>
          <input
            type="text"
            name="cvv"
            value={formData.cvv}
            onChange={handleChange}
            placeholder="123"
          />
          {formErrors.cvv && <p className="error">{formErrors.cvv}</p>}
        </div>

        {formErrors.submit && <p className="error">{formErrors.submit}</p>}

        <button type="submit" disabled={isProcessing}>
          {isProcessing ? 'Przetwarzanie...' : 'Zapłać'}
        </button>
      </form>
    </div>
  );
};

export default Payment;
