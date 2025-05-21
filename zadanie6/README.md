
- 20+ przypadków testowych
- Ponad 50 asercji
- Testy UI dla frontendu
- Testy API dla backendu
- Konfiguracja do uruchamiania testów na BrowserStack

## Konfiguracja

1. Ustaw zmienne środowiskowe:

```bash
export BROWSERSTACK_USERNAME="twoja_nazwa_użytkownika"
export BROWSERSTACK_ACCESS_KEY="twój_klucz_dostępu"
```

2. Uruchom najpierw aplikację z zadania 5:

## Uruchamianie testów

```bash
cd zadanie6
go test ./tests -v
```
## Lista testów

### Testy UI
1. TestHomePage - Weryfikacja poprawnego ładowania strony głównej
2. TestCartNavigation - Test nawigacji do koszyka
3. TestAddToCart - Test dodawania produktu do koszyka
4. TestCartQuantityIncrement - Test zwiększania ilości w koszyku
5. TestRemoveFromCart - Test usuwania produktu z koszyka
6. TestPaymentNavigation - Test nawigacji do strony płatności
7. TestPaymentFormValidation - Test walidacji formularza płatności
8. TestSuccessfulPayment - Test pomyślnego procesu płatności
9. TestProductDetails - Test poprawnego wyświetlania szczegółów produktu
10. TestResponsiveness - Test responsywności na urządzeniach mobilnych

### Testy API
11. TestGetProducts - Test pobierania produktów (pozytywny)
12. TestGetProductsNegative - Test pobierania produktów (negatywny)
13. TestGetCart - Test pobierania koszyka (pozytywny)
14. TestAddToCart - Test dodawania do koszyka (pozytywny)
15. TestAddToCartNegative - Test dodawania do koszyka (negatywny)
16. TestGetPayments - Test pobierania płatności (pozytywny)
17. TestCreatePayment - Test tworzenia płatności (pozytywny)
18. TestCreatePaymentNegative - Test tworzenia płatności (negatywny)
19. TestCartClearedAfterPayment - Test, czy koszyk jest czyszczony po płatności
20. TestProductFieldsIntegrity - Test integralności pól produktu

### Testy integracyjne
21. TestProductLoadingFromAPI - Test integracji wyświetlania produktów z API
22. TestCartPersistenceBetweenPages - Test persystencji koszyka między stronami
23. TestPaymentValidationConsistency - Test spójności walidacji formularza płatności
