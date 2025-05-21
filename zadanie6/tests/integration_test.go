package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

// Test 21: Test API call integration with frontend product display
func TestProductLoadingFromAPI(t *testing.T) {
	if webDriver == nil {
		t.Skip("Skipping integration test because WebDriver is not available")
		return
	}

	if err := webDriver.Get(baseURL); err != nil {
		t.Fatalf("Failed to load homepage: %s", err)
	}

	time.Sleep(2 * time.Second)

	frontendProducts, err := findElements(selenium.ByCSSSelector, ".product-card")
	if err != nil {
		t.Fatalf("Failed to find product cards: %s", err)
	}

	client := &http.Client{}
	resp, err := client.Get(apiBaseURL + "/products")
	if err != nil {
		t.Fatalf("Failed to make API request: %s", err)
	}
	defer resp.Body.Close()

	var apiProducts []Product
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&apiProducts); err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}

	if len(frontendProducts) != len(apiProducts) {
		t.Errorf("Frontend shows %d products but API returned %d products",
			len(frontendProducts), len(apiProducts))
	}

	if len(frontendProducts) > 0 && len(apiProducts) > 0 {
		firstProductNameElement := assertElementExists(t, selenium.ByCSSSelector,
			".product-card:first-child h3", "first product name")
		frontendName, err := firstProductNameElement.Text()
		if err != nil {
			t.Fatalf("Failed to get frontend product name: %s", err)
		}

		apiName := apiProducts[0].Name

		if !strings.EqualFold(strings.TrimSpace(frontendName), strings.TrimSpace(apiName)) {
			t.Errorf("First product name mismatch. Frontend: '%s', API: '%s'",
				frontendName, apiName)
		}
	}
}

// Test 22: Test cart persistence between page navigation
func TestCartPersistenceBetweenPages(t *testing.T) {
	if webDriver == nil {
		t.Skip("Skipping integration test because WebDriver is not available")
		return
	}

	if err := webDriver.Get(baseURL); err != nil {
		t.Fatalf("Failed to load homepage: %s", err)
	}

	maybeClickCodespacesContinue(t)
	time.Sleep(3 * time.Second)

	productCards, err := findElements(selenium.ByCSSSelector, ".product-card")
	if err != nil || len(productCards) == 0 {
		t.Fatalf("No product cards found on page: %v", err)
	}

	addButton, err := productCards[0].FindElement(selenium.ByCSSSelector, "button")
	if err != nil {
		t.Fatalf("Failed to find add to cart button within product card: %s", err)
	}

	if err := addButton.Click(); err != nil {
		t.Fatalf("Failed to click add to cart button: %s", err)
	}

	time.Sleep(2 * time.Second)

	cartLink := waitForElement(t, selenium.ByCSSSelector, "a[href='/cart']", 5*time.Second)
	cartLinkText, _ := cartLink.Text()
	t.Logf("Cart link text after adding item: %s", cartLinkText)

	if err := webDriver.Get(baseURL + "/cart"); err != nil {
		t.Fatalf("Failed to navigate to cart page: %s", err)
	}

	maybeClickCodespacesContinue(t)
	time.Sleep(2 * time.Second)

	emptyMsg, err := findElements(selenium.ByCSSSelector, "p")
	if err == nil && len(emptyMsg) > 0 {
		for _, msg := range emptyMsg {
			text, _ := msg.Text()
			if text == "Twój koszyk jest pusty." {
				t.Errorf("Cart is empty after adding item")
				return
			}
		}
	}

	_, tableErr := findElement(selenium.ByCSSSelector, ".cart-table")
	if tableErr != nil {
		t.Errorf("Cart table not found after adding item: %v", tableErr)
		return
	}

	if err := webDriver.Get(baseURL); err != nil {
		t.Fatalf("Failed to navigate back to home page: %s", err)
	}

	maybeClickCodespacesContinue(t)
	time.Sleep(2 * time.Second)

	cartLink = waitForElement(t, selenium.ByCSSSelector, "a[href='/cart']", 5*time.Second)
	finalCartText, _ := cartLink.Text()
	t.Logf("Cart link text after returning to homepage: %s", finalCartText)

	if strings.Contains(finalCartText, "Koszyk") && !strings.Contains(finalCartText, "Koszyk0") {
		t.Log("Cart persisted as expected after navigation")
	} else {
		t.Errorf("Cart appears to be empty after navigation")
	}
}

// Test 23: Test that payment form validation matches API validation
func TestPaymentValidationConsistency(t *testing.T) {
	maybeClickCodespacesContinue(t)
	if webDriver == nil {
		t.Skip("Skipping integration test because WebDriver is not available")
		return
	}

	if err := webDriver.Get(baseURL); err != nil {
		t.Fatalf("Failed to load homepage: %s", err)
	}

	addToCartButton := assertElementExists(t, selenium.ByCSSSelector,
		".product-card button", "add to cart button")
	if err := addToCartButton.Click(); err != nil {
		t.Fatalf("Failed to click add to cart button: %s", err)
	}

	if err := webDriver.Get(baseURL + "/payment"); err != nil {
		t.Fatalf("Failed to navigate to payment page: %s", err)
	}

	cardNumberField := assertElementExists(t, selenium.ByCSSSelector,
		"input[name='cardNumber']", "card number field")
	if err := cardNumberField.Clear(); err != nil {
		t.Fatalf("Failed to clear card number field: %s", err)
	}
	if err := cardNumberField.SendKeys("1234"); err != nil {
		t.Fatalf("Failed to enter invalid card number: %s", err)
	}

	submitButton := assertElementExists(t, selenium.ByCSSSelector,
		"button[type='submit']", "submit button")
	if err := submitButton.Click(); err != nil {
		t.Fatalf("Failed to click submit button: %s", err)
	}

	frontendError, err := findElement(selenium.ByCSSSelector, ".error")
	if err != nil {
		t.Fatalf("Failed to find error message: %s", err)
	}
	frontendErrorText, _ := frontendError.Text()

	client := &http.Client{}

	payment := Payment{
		Amount:     99.99,
		CardNumber: "1234",
		CardHolder: "Test User",
		ExpiryDate: "12/25",
		CVV:        "123",
	}

	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		t.Fatalf("Failed to marshal payment: %s", err)
	}

	resp, err := client.Post(apiBaseURL+"/payments", "application/json",
		bytes.NewBuffer(paymentJSON))
	if err != nil {
		t.Fatalf("Failed to make API request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		t.Errorf("API accepted invalid card data that frontend rejected")
	}

	t.Logf("Frontend validation error: %s", frontendErrorText)
	t.Logf("API response status code: %d", resp.StatusCode)
}
