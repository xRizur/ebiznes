package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func logAPIResponse(t *testing.T, resp *http.Response) {
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	body := string(bodyBytes)

	if len(body) > 200 && strings.Contains(body, "<!doctype html>") {
		body = body[:200] + "... [truncated HTML response]"
	}

	t.Logf("API Response Status: %d, Body: %s", resp.StatusCode, body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
}

func isJSONResponse(body []byte) bool {
	return !strings.Contains(string(body), "<!doctype html>") && !strings.Contains(string(body), "<html>")
}

type Product struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageUrl"`
}

type CartItem struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"productId"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
}

type Payment struct {
	ID         uint    `json:"id"`
	Amount     float64 `json:"amount"`
	CardNumber string  `json:"cardNumber"`
	CardHolder string  `json:"cardHolder"`
	ExpiryDate string  `json:"expiryDate"`
	CVV        string  `json:"cvv"`
	Status     string  `json:"status"`
}

// Test 11: Test GET /products endpoint - positive scenario
func TestGetProducts(t *testing.T) {
	t.Logf("Testing API endpoint: %s/products", apiBaseURL)

	client := &http.Client{}

	resp, err := client.Get(apiBaseURL + "/products")
	if err != nil {
		t.Fatalf("Failed to make request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var products []Product
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&products); err != nil {
		t.Fatalf("Failed to decode response: %s", err)
	}

	if len(products) == 0 {
		t.Errorf("Expected products, got empty array")
	}

	product := products[0]
	if product.ID == 0 {
		t.Errorf("Product ID is missing")
	}
	if product.Name == "" {
		t.Errorf("Product name is empty")
	}
	if product.Description == "" {
		t.Errorf("Product description is empty")
	}
	if product.Price <= 0 {
		t.Errorf("Product price is invalid: %f", product.Price)
	}
	if product.ImageURL == "" {
		t.Errorf("Product image URL is empty")
	}
}

// Test 12: Test GET /products endpoint - negative scenario (server unavailable)
func TestGetProductsNegative(t *testing.T) {
	client := &http.Client{}

	resp, err := client.Get(apiBaseURL + "/products/nonexistent")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Errorf("Expected non-200 status code, got %d", resp.StatusCode)
	}
}

// Test 13: Test GET /cart endpoint - positive scenario
func TestGetCart(t *testing.T) {
	client := &http.Client{}

	// Make GET request to /cart
	resp, err := client.Get(apiBaseURL + "/cart")
	if err != nil {
		t.Fatalf("Failed to make request: %s", err)
	}
	defer resp.Body.Close()

	// Assert status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Parse response body
	var cartItems []CartItem
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&cartItems); err != nil {
		t.Fatalf("Failed to decode response: %s", err)
	}

	// Cart might be empty, which is a valid scenario
	// Just ensure we can parse the response correctly
}

// Test 14: Test POST /cart endpoint - positive scenario
func TestAddToCart(t *testing.T) {
	client := &http.Client{}

	// First get a product to add to cart
	respProducts, err := client.Get(apiBaseURL + "/products")
	if err != nil {
		t.Fatalf("Failed to get products: %s", err)
	}
	defer respProducts.Body.Close()

	var products []Product
	decoder := json.NewDecoder(respProducts.Body)
	if err := decoder.Decode(&products); err != nil {
		t.Fatalf("Failed to decode products: %s", err)
	}

	if len(products) == 0 {
		t.Fatalf("No products available for test")
	}

	productID := products[0].ID

	// Create a new cart item
	cartItem := CartItem{
		ProductID: productID,
		Quantity:  1,
	}

	// Convert to JSON
	cartItemJSON, err := json.Marshal(cartItem)
	if err != nil {
		t.Fatalf("Failed to marshal cart item: %s", err)
	}

	// Make POST request to /cart
	resp, err := client.Post(apiBaseURL+"/cart", "application/json", bytes.NewBuffer(cartItemJSON))
	if err != nil {
		t.Fatalf("Failed to make request: %s", err)
	}
	defer resp.Body.Close()

	// Assert status code is 201 Created
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	// Parse response body
	var createdItem CartItem
	decoder = json.NewDecoder(resp.Body)
	if err := decoder.Decode(&createdItem); err != nil {
		t.Fatalf("Failed to decode response: %s", err)
	}

	// Assert created item has expected fields
	if createdItem.ID == 0 {
		t.Errorf("Created item ID is missing")
	}
	if createdItem.ProductID != productID {
		t.Errorf("Expected product ID %d, got %d", productID, createdItem.ProductID)
	}
	if createdItem.Quantity != 1 {
		t.Errorf("Expected quantity 1, got %d", createdItem.Quantity)
	}
}

// Test 15: Test POST /cart endpoint - negative scenario (invalid product ID)
func TestAddToCartNegative(t *testing.T) {
	client := &http.Client{}

	cartItem := CartItem{
		ProductID: 9999,
		Quantity:  1,
	}

	cartItemJSON, err := json.Marshal(cartItem)
	if err != nil {
		t.Fatalf("Failed to marshal cart item: %s", err)
	}

	resp, err := client.Post(apiBaseURL+"/cart", "application/json", bytes.NewBuffer(cartItemJSON))
	if err != nil {
		t.Fatalf("Failed to make request: %s", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	t.Logf("API Response Status: %d, Body: %s", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 404 or 400, got %d", resp.StatusCode)
	}
}

// Test 16: Test GET /payments endpoint - positive scenario
func TestGetPayments(t *testing.T) {
	client := &http.Client{}

	resp, err := client.Get(apiBaseURL + "/payments")
	if err != nil {
		t.Fatalf("Failed to make request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var payments []Payment
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&payments); err != nil {
		t.Fatalf("Failed to decode response: %s", err)
	}

}

// Test 17: Test POST /payments endpoint - positive scenario
func TestCreatePayment(t *testing.T) {
	client := &http.Client{}

	payment := Payment{
		Amount:     99.99,
		CardNumber: "4242424242424242",
		CardHolder: "Test User",
		ExpiryDate: "12/25",
		CVV:        "123",
	}

	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		t.Fatalf("Failed to marshal payment: %s", err)
	}

	resp, err := client.Post(apiBaseURL+"/payments", "application/json", bytes.NewBuffer(paymentJSON))
	if err != nil {
		t.Fatalf("Failed to make request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var createdPayment Payment
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&createdPayment); err != nil {
		t.Fatalf("Failed to decode response: %s", err)
	}

	if createdPayment.ID == 0 {
		t.Errorf("Created payment ID is missing")
	}
	if createdPayment.Amount != payment.Amount {
		t.Errorf("Expected amount %f, got %f", payment.Amount, createdPayment.Amount)
	}
	if createdPayment.CardHolder != payment.CardHolder {
		t.Errorf("Expected card holder '%s', got '%s'", payment.CardHolder, createdPayment.CardHolder)
	}
	if createdPayment.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", createdPayment.Status)
	}
}

// Test 18: Test POST /payments endpoint - negative scenario (invalid card data)
func TestCreatePaymentNegative(t *testing.T) {
	client := &http.Client{}

	payment := Payment{
		Amount:     99.99,
		CardNumber: "invalid",
		CardHolder: "",
		ExpiryDate: "invalid",
		CVV:        "invalid",
	}

	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		t.Fatalf("Failed to marshal payment: %s", err)
	}

	resp, err := client.Post(apiBaseURL+"/payments", "application/json", bytes.NewBuffer(paymentJSON))
	if err != nil {
		t.Fatalf("Failed to make request: %s", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	t.Logf("API Response Status: %d, Body: %s", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode == http.StatusCreated {
		t.Errorf("Expected non-201 status code, got %d", resp.StatusCode)
	}
}

// Test 19: Test cart is cleared after successful payment
func TestCartClearedAfterPayment(t *testing.T) {
	client := &http.Client{}

	respProducts, err := client.Get(apiBaseURL + "/products")
	if err != nil {
		t.Fatalf("Failed to get products: %s", err)
	}

	var products []Product
	decoder := json.NewDecoder(respProducts.Body)
	if err := decoder.Decode(&products); err != nil {
		t.Fatalf("Failed to decode products: %s", err)
	}
	respProducts.Body.Close()

	if len(products) == 0 {
		t.Fatalf("No products available for test")
	}

	cartItem := CartItem{
		ProductID: products[0].ID,
		Quantity:  1,
	}

	cartItemJSON, err := json.Marshal(cartItem)
	if err != nil {
		t.Fatalf("Failed to marshal cart item: %s", err)
	}

	_, err = client.Post(apiBaseURL+"/cart", "application/json", bytes.NewBuffer(cartItemJSON))
	if err != nil {
		t.Fatalf("Failed to add item to cart: %s", err)
	}

	payment := Payment{
		Amount:     products[0].Price,
		CardNumber: "4242424242424242",
		CardHolder: "Test User",
		ExpiryDate: "12/25",
		CVV:        "123",
	}

	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		t.Fatalf("Failed to marshal payment: %s", err)
	}

	_, err = client.Post(apiBaseURL+"/payments", "application/json", bytes.NewBuffer(paymentJSON))
	if err != nil {
		t.Fatalf("Failed to process payment: %s", err)
	}

	respCart, err := client.Get(apiBaseURL + "/cart")
	if err != nil {
		t.Fatalf("Failed to get cart: %s", err)
	}
	defer respCart.Body.Close()

	var cartItems []CartItem
	decoder = json.NewDecoder(respCart.Body)
	if err := decoder.Decode(&cartItems); err != nil {
		t.Fatalf("Failed to decode cart: %s", err)
	}

	if len(cartItems) > 0 {
		t.Errorf("Expected empty cart after payment, got %d items", len(cartItems))
	}
}

// Test 20: Test all product fields are returned correctly
func TestProductFieldsIntegrity(t *testing.T) {
	client := &http.Client{}

	resp, err := client.Get(apiBaseURL + "/products")
	if err != nil {
		t.Fatalf("Failed to get products: %s", err)
	}
	defer resp.Body.Close()

	var products []Product
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&products); err != nil {
		t.Fatalf("Failed to decode products: %s", err)
	}

	if len(products) == 0 {
		t.Fatalf("No products available for test")
	}

	for i, product := range products {
		t.Run(fmt.Sprintf("Product-%d", i), func(t *testing.T) {
			if product.ID == 0 {
				t.Errorf("Product ID is missing")
			}

			if product.Name == "" {
				t.Errorf("Product name is empty")
			}

			if product.Description == "" {
				t.Errorf("Product description is empty")
			}

			if product.Price <= 0 {
				t.Errorf("Product price is invalid: %f", product.Price)
			}

			if product.ImageURL == "" {
				t.Errorf("Product image URL is empty")
			}
		})
	}
}
