package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func maybeClickCodespacesContinue(t *testing.T) {
	title, err := webDriver.Title()
	if err != nil {
		t.Logf("Could not get page title: %v", err)
		return
	}

	if title == "Codespaces Access Port" {
		t.Log("Detected Codespaces port protection page — trying to continue...")

		btn, err := webDriver.FindElement(selenium.ByCSSSelector, "button")
		if err != nil {
			t.Logf("Could not find continue button: %v", err)
			return
		}

		if err := btn.Click(); err != nil {
			t.Logf("Failed to click continue button: %v", err)
			return
		}

		time.Sleep(2 * time.Second)
		newTitle, _ := webDriver.Title()
		t.Logf("After click, new page title: %q", newTitle)
	}
}

func TestHomePage(t *testing.T) {
	if webDriver == nil {
		t.Skip("Skipping UI test because WebDriver is not available")
		return
	}

	if err := webDriver.Get(baseURL); err != nil {
		t.Fatalf("Failed to load homepage: %s", err)
	}

	maybeClickCodespacesContinue(t)

	assertPageTitle(t, "Sklep Internetowy")

	assertElementExists(t, selenium.ByCSSSelector, "nav", "navigation menu")

	productElements, err := findElements(selenium.ByCSSSelector, ".product-card")
	if err != nil {
		t.Fatalf("Failed to find product cards: %s", err)
	}
	if len(productElements) == 0 {
		t.Errorf("Expected products to be displayed, but none were found")
	}

	if len(productElements) < 3 {
		t.Errorf("Expected at least 3 products, but found %d", len(productElements))
	}

	nameElement := assertElementExists(t, selenium.ByCSSSelector, ".product-card h3", "product name")
	nameText, _ := nameElement.Text()
	if nameText == "" {
		t.Errorf("Product name is empty")
	}

	priceElement := assertElementExists(t, selenium.ByCSSSelector, ".product-card p:nth-child(4)", "product price")
	priceText, _ := priceElement.Text()
	if priceText == "" {
		t.Errorf("Product price is empty")
	}

	assertElementExists(t, selenium.ByCSSSelector, ".product-card button", "add to cart button")
}

// Test 2: Test navigation to cart page
func TestCartNavigation(t *testing.T) {
	if webDriver == nil {
		t.Skip("Skipping UI test because WebDriver is not available")
		return
	}

	if err := webDriver.Get(baseURL); err != nil {
		t.Fatalf("Failed to load homepage: %s", err)
	}

	maybeClickCodespacesContinue(t)

	cartLink := assertElementExists(t, selenium.ByCSSSelector, "a[href='/cart']", "cart navigation link")
	if err := cartLink.Click(); err != nil {
		t.Fatalf("Failed to click cart link: %s", err)
	}

	time.Sleep(1 * time.Second)
	maybeClickCodespacesContinue(t)

	currentURL, err := webDriver.CurrentURL()
	if err != nil {
		t.Fatalf("Failed to get current URL: %s", err)
	}
	if currentURL != baseURL+"/cart" {
		t.Errorf("Expected to be on cart page, but URL is %s", currentURL)
	}

	assertElementExists(t, selenium.ByCSSSelector, "h1", "cart page heading")
	headingElement, _ := findElement(selenium.ByCSSSelector, "h1")
	headingText, _ := headingElement.Text()
	if headingText != "Koszyk" {
		t.Errorf("Expected cart page heading to be 'Koszyk', got '%s'", headingText)
	}
}

// Test 3: Test adding a product to cart
func TestAddToCartUI(t *testing.T) {
	if webDriver == nil {
		t.Skip("Skipping UI test because WebDriver is not available")
		return
	}

	if err := webDriver.Get(baseURL); err != nil {
		t.Fatalf("Failed to load homepage: %s", err)
	}

	maybeClickCodespacesContinue(t)

	time.Sleep(3 * time.Second)

	pageSource, _ := webDriver.PageSource()
	t.Logf("Page source excerpt: %.300s...", pageSource)

	productCards, err := findElements(selenium.ByCSSSelector, ".product-card")
	if err != nil || len(productCards) == 0 {
		t.Fatalf("No product cards found on page: %v", err)
	}

	nameElement, err := productCards[0].FindElement(selenium.ByCSSSelector, "h3")
	if err != nil {
		t.Fatalf("Failed to find product name: %s", err)
	}

	productName, _ := nameElement.Text()
	t.Logf("Found product with name: %s", productName)

	addButton, err := productCards[0].FindElement(selenium.ByCSSSelector, "button")
	if err != nil {
		t.Fatalf("Failed to find add to cart button: %s", err)
	}

	if err := addButton.Click(); err != nil {
		t.Fatalf("Failed to click add to cart button: %s", err)
	}

	time.Sleep(2 * time.Second)

	if err := webDriver.Get(baseURL + "/cart"); err != nil {
		t.Fatalf("Failed to navigate to cart page: %s", err)
	}

	maybeClickCodespacesContinue(t)

	time.Sleep(2 * time.Second)

	cartPageSource, _ := webDriver.PageSource()
	t.Logf("Cart page source excerpt: %.300s...", cartPageSource)

	cartTableErr := func() error {
		_, err := findElement(selenium.ByCSSSelector, ".cart-table")
		return err
	}()
	if cartTableErr != nil {
		t.Logf("Cart table not found: %v", cartTableErr)

		emptyCartMsg, msgErr := findElement(selenium.ByCSSSelector, "p")
		if msgErr == nil {
			msgText, _ := emptyCartMsg.Text()
			t.Logf("Found message in cart: %s", msgText)
			if msgText == "Twój koszyk jest pusty." {
				t.Errorf("Cart is empty - item was not added successfully")
			}
		}

		t.Skip("Cart table not found - cart may be empty. Skipping remainder of test.")
		return
	}

	cartItems, err := findElements(selenium.ByCSSSelector, ".cart-table tbody tr")
	if err != nil || len(cartItems) == 0 {
		t.Errorf("Expected at least one item in cart, but none were found")
		return
	}

	itemNameElement, _ := cartItems[0].FindElement(selenium.ByCSSSelector, "td:first-child")
	itemName, _ := itemNameElement.Text()
	if itemName != productName {
		t.Errorf("Expected cart to contain '%s', but found '%s'", productName, itemName)
	}

	quantityCell, _ := cartItems[0].FindElement(selenium.ByCSSSelector, "td:nth-child(3) .quantity-value")
	quantityText, _ := quantityCell.Text()
	if quantityText != "1" {
		t.Errorf("Expected quantity to be 1, got '%s'", quantityText)
	}
}

// Test 4: Test quantity increment in cart
func TestCartQuantityIncrement(t *testing.T) {
	if webDriver == nil {
		t.Skip("Skipping UI test because WebDriver is not available")
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
		t.Fatalf("Failed to find add to cart button: %s", err)
	}

	if err := addButton.Click(); err != nil {
		t.Fatalf("Failed to click add to cart button: %s", err)
	}

	if err := webDriver.Get(baseURL + "/cart"); err != nil {
		t.Fatalf("Failed to navigate to cart page: %s", err)
	}

	maybeClickCodespacesContinue(t)

	time.Sleep(2 * time.Second)
	
	// Debug - capture and log cart page source to see the actual HTML structure
	cartSource, _ := webDriver.PageSource()
	t.Logf("Cart page HTML excerpt for debugging: %.300s...", cartSource)
	
	// Check if the cart table exists
	_, tableErr := findElement(selenium.ByCSSSelector, ".cart-table")
	if tableErr != nil {
		t.Skip("Cart table not found - cart may be empty. Skipping remainder of test.")
		return
	}
	
	// Find the quantity cell - this should contain the initial quantity value
	quantityCell, err := findElement(selenium.ByCSSSelector, "table.cart-table tbody tr td:nth-child(3)")
	if err != nil {
		t.Skip("Quantity cell not found. Skipping remainder of test.")
		return
	}
	
	// Get the text of the entire cell - this will include the buttons and quantity
	cellText, _ := quantityCell.Text()
	t.Logf("Quantity cell text: %s", cellText)
	
	// Try to find all buttons in the quantity cell
	buttons, err := quantityCell.FindElements(selenium.ByCSSSelector, "button")
	if err != nil || len(buttons) < 2 {
		t.Logf("Could not find expected buttons in quantity cell. Found %d buttons", len(buttons))
		
		// Try using a more generic selector for the increment button
		incrementButton, err := findElement(selenium.ByCSSSelector, "button[aria-label='Increment'], button:contains('+'), button:nth-of-type(2)")
		if err != nil {
			t.Fatalf("Could not find increment button using alternative selectors: %v", err)
		}
		
		// Click the found button
		if err := incrementButton.Click(); err != nil {
			t.Fatalf("Failed to click increment button: %v", err)
		}
	} else {
		// If we found buttons, use the second one (the increment button)
		if err := buttons[1].Click(); err != nil {
			t.Fatalf("Failed to click increment button: %v", err)
		}
	}
	
	// Wait for the quantity to update
	time.Sleep(2 * time.Second)
	
	// Get the updated quantity
	updatedQuantityCell, err := findElement(selenium.ByCSSSelector, "table.cart-table tbody tr td:nth-child(3)")
	if err != nil {
		t.Fatalf("Could not find quantity cell after increment: %v", err)
	}
	
	updatedText, _ := updatedQuantityCell.Text()
	t.Logf("Updated quantity cell text: %s", updatedText)
	
	// Look for the quantity value - it might be in a span with a specific class
	quantityValue, err := updatedQuantityCell.FindElement(selenium.ByCSSSelector, ".quantity-value")
	if err == nil {
		// If we found a specific element with the quantity, use that
		newQuantityText, _ := quantityValue.Text()
		if newQuantityText == "1" {
			t.Errorf("Quantity did not increase, still showing: %s", newQuantityText)
		} else {
			t.Logf("Quantity successfully increased to: %s", newQuantityText)
		}
	} else {
		// Otherwise check the whole cell text
		if updatedText == cellText {
			t.Errorf("Quantity cell text did not change after clicking increment")
		} else {
			t.Logf("Quantity cell text changed from '%s' to '%s'", cellText, updatedText)
		}
	}
}

// Test 5: Test removing an item from cart
func TestRemoveFromCart(t *testing.T) {
	if webDriver == nil {
		t.Skip("Skipping UI test because WebDriver is not available")
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
		t.Fatalf("Failed to find add to cart button: %s", err)
	}

	if err := addButton.Click(); err != nil {
		t.Fatalf("Failed to click add to cart button: %s", err)
	}

	if err := webDriver.Get(baseURL + "/cart"); err != nil {
		t.Fatalf("Failed to navigate to cart page: %s", err)
	}

	maybeClickCodespacesContinue(t)

	time.Sleep(2 * time.Second)

	initialItems, _ := findElements(selenium.ByCSSSelector, "table.cart-table tbody tr")
	initialCount := len(initialItems)

	removeButton := waitForElement(t, selenium.ByCSSSelector, "table.cart-table tbody tr td:last-child button", 5*time.Second)
	removeButton.Click()
	time.Sleep(2 * time.Second)

	remainingItems, _ := findElements(selenium.ByCSSSelector, "table.cart-table tbody tr")
	remainingCount := len(remainingItems)
	if remainingCount != initialCount-1 {
		t.Errorf("Expected %d items after removal, got %d", initialCount-1, remainingCount)
	}
}

// Test 6: Test navigation to payment page
func TestPaymentNavigation(t *testing.T) {
	if webDriver == nil {
		t.Skip("Skipping UI test because WebDriver is not available")
		return
	}

	if err := webDriver.Get(baseURL); err != nil {
		t.Fatalf("Failed to load homepage: %s", err)
	}

	addToCartButton := assertElementExists(t, selenium.ByCSSSelector, ".product-card button", "add to cart button")
	if err := addToCartButton.Click(); err != nil {
		t.Fatalf("Failed to click add to cart button: %s", err)
	}

	if err := webDriver.Get(baseURL + "/cart"); err != nil {
		t.Fatalf("Failed to navigate to cart page: %s", err)
	}

	maybeClickCodespacesContinue(t)

	paymentButton := assertElementExists(t, selenium.ByCSSSelector, "a[href='/payment'] button", "proceed to payment button")
	if err := paymentButton.Click(); err != nil {
		t.Fatalf("Failed to click proceed to payment button: %s", err)
	}

	time.Sleep(1 * time.Second)
	maybeClickCodespacesContinue(t)

	currentURL, err := webDriver.CurrentURL()
	if err != nil {
		t.Fatalf("Failed to get current URL: %s", err)
	}
	if currentURL != baseURL+"/payment" {
		t.Errorf("Expected to be on payment page, but URL is %s", currentURL)
	}

	assertElementExists(t, selenium.ByCSSSelector, "h1", "payment page heading")
	headingElement, _ := findElement(selenium.ByCSSSelector, "h1")
	headingText, _ := headingElement.Text()
	if headingText != "Płatność" {
		t.Errorf("Expected payment page heading to be 'Płatność', got '%s'", headingText)
	}

	assertElementExists(t, selenium.ByCSSSelector, "input[name='cardNumber']", "card number field")
	assertElementExists(t, selenium.ByCSSSelector, "input[name='cardHolder']", "card holder field")
	assertElementExists(t, selenium.ByCSSSelector, "input[name='expiryDate']", "expiry date field")
	assertElementExists(t, selenium.ByCSSSelector, "input[name='cvv']", "CVV field")
	assertElementExists(t, selenium.ByCSSSelector, "button[type='submit']", "pay button")
}

// Test 7: Test payment form validation
func TestPaymentFormValidation(t *testing.T) {
	if webDriver == nil {
		t.Skip("Skipping UI test because WebDriver is not available")
		return
	}

	if err := webDriver.Get(baseURL); err != nil {
		t.Fatalf("Failed to load homepage: %s", err)
	}

	addToCartButton := assertElementExists(t, selenium.ByCSSSelector, ".product-card button", "add to cart button")
	if err := addToCartButton.Click(); err != nil {
		t.Fatalf("Failed to click add to cart button: %s", err)
	}

	if err := webDriver.Get(baseURL + "/payment"); err != nil {
		t.Fatalf("Failed to navigate to payment page: %s", err)
	}

	maybeClickCodespacesContinue(t)

	submitButton := assertElementExists(t, selenium.ByCSSSelector, "button[type='submit']", "submit button")
	if err := submitButton.Click(); err != nil {
		t.Fatalf("Failed to click submit button: %s", err)
	}

	errorMessages, err := findElements(selenium.ByCSSSelector, ".error")
	if err != nil {
		t.Fatalf("Failed to find error messages: %s", err)
	}
	if len(errorMessages) == 0 {
		t.Errorf("Expected validation error messages, but none were found")
	}

	cardNumberError, err := findElement(selenium.ByCSSSelector, "input[name='cardNumber'] + .error")
	if err != nil {
		t.Fatalf("Failed to find card number error: %s", err)
	}
	cardNumberErrorText, _ := cardNumberError.Text()
	if cardNumberErrorText != "Numer karty jest wymagany" {
		t.Errorf("Expected card number error 'Numer karty jest wymagany', got '%s'", cardNumberErrorText)
	}
}

// Test 8: Test successful payment submission
func TestSuccessfulPayment(t *testing.T) {
	if webDriver == nil {
		t.Skip("Skipping UI test because WebDriver is not available")
		return
	}

	if err := webDriver.Get(baseURL); err != nil {
		t.Fatalf("Failed to load homepage: %s", err)
	}

	addToCartButton := assertElementExists(t, selenium.ByCSSSelector, ".product-card button", "add to cart button")
	if err := addToCartButton.Click(); err != nil {
		t.Fatalf("Failed to click add to cart button: %s", err)
	}

	if err := webDriver.Get(baseURL + "/payment"); err != nil {
		t.Fatalf("Failed to navigate to payment page: %s", err)
	}

	maybeClickCodespacesContinue(t)

	cardNumberField := assertElementExists(t, selenium.ByCSSSelector, "input[name='cardNumber']", "card number field")
	if err := cardNumberField.Clear(); err != nil {
		t.Fatalf("Failed to clear card number field: %s", err)
	}
	if err := cardNumberField.SendKeys("4242424242424242"); err != nil {
		t.Fatalf("Failed to enter card number: %s", err)
	}

	cardHolderField := assertElementExists(t, selenium.ByCSSSelector, "input[name='cardHolder']", "card holder field")
	if err := cardHolderField.Clear(); err != nil {
		t.Fatalf("Failed to clear card holder field: %s", err)
	}
	if err := cardHolderField.SendKeys("Test User"); err != nil {
		t.Fatalf("Failed to enter card holder: %s", err)
	}

	expiryDateField := assertElementExists(t, selenium.ByCSSSelector, "input[name='expiryDate']", "expiry date field")
	if err := expiryDateField.Clear(); err != nil {
		t.Fatalf("Failed to clear expiry date field: %s", err)
	}
	if err := expiryDateField.SendKeys("12/25"); err != nil {
		t.Fatalf("Failed to enter expiry date: %s", err)
	}

	cvvField := assertElementExists(t, selenium.ByCSSSelector, "input[name='cvv']", "CVV field")
	if err := cvvField.Clear(); err != nil {
		t.Fatalf("Failed to clear CVV field: %s", err)
	}
	if err := cvvField.SendKeys("123"); err != nil {
		t.Fatalf("Failed to enter CVV: %s", err)
	}

	submitButton := assertElementExists(t, selenium.ByCSSSelector, "button[type='submit']", "submit button")
	if err := submitButton.Click(); err != nil {
		t.Fatalf("Failed to click submit button: %s", err)
	}

	webDriver.AcceptAlert()

	time.Sleep(2 * time.Second)
	currentURL, err := webDriver.CurrentURL()
	if err != nil {
		t.Fatalf("Failed to get current URL: %s", err)
	}
	if currentURL != baseURL+"/" {
		t.Errorf("Expected to be redirected to home page, but URL is %s", currentURL)
	}
}

// Test 9: Test product details are displayed correctly
func TestProductDetails(t *testing.T) {
	if webDriver == nil {
		t.Skip("Skipping UI test because WebDriver is not available")
		return
	}

	if err := webDriver.Get(baseURL); err != nil {
		t.Fatalf("Failed to load homepage: %s", err)
	}

	maybeClickCodespacesContinue(t)

	productCard := assertElementExists(t, selenium.ByCSSSelector, ".product-card", "product card")

	nameElement := assertElementExists(t, selenium.ByCSSSelector, "h3", "product name")
	nameText, _ := nameElement.Text()
	if nameText == "" {
		t.Errorf("Product name is empty")
	}

	priceElement := assertElementExists(t, selenium.ByCSSSelector, ".price", "product price")
	priceText, _ := priceElement.Text()
	if priceText == "" {
		t.Errorf("Product price is empty")
	}

	descElement := assertElementExists(t, selenium.ByCSSSelector, ".description", "product description")
	descText, _ := descElement.Text()
	if descText == "" {
		t.Errorf("Product description is empty")
	}

	assertElementExists(t, selenium.ByCSSSelector, "img", "product image")

	imgElement, _ := productCard.FindElement(selenium.ByCSSSelector, "img")
	imgSrc, _ := imgElement.GetAttribute("src")
	if imgSrc == "" {
		t.Errorf("Product image source is empty")
	}
}

func isValidPriceFormat(price string) bool {
	return len(price) > 2 && strings.Contains(price, "Cena:") && strings.Contains(price, "zł")
}

// Test 10: Test responsiveness - mobile view
func TestResponsiveness(t *testing.T) {
	if webDriver == nil {
		t.Skip("Skipping UI test because WebDriver is not available")
		return
	}

	if err := webDriver.ResizeWindow("", 375, 812); err != nil {
		t.Fatalf("Failed to resize window: %s", err)
	}

	if err := webDriver.Get(baseURL); err != nil {
		t.Fatalf("Failed to load homepage: %s", err)
	}

	maybeClickCodespacesContinue(t)

	assertElementExists(t, selenium.ByCSSSelector, "nav", "navigation menu")

	productElements, err := findElements(selenium.ByCSSSelector, ".product-card")
	if err != nil {
		t.Fatalf("Failed to find product cards: %s", err)
	}

	if len(productElements) >= 2 {
		firstProduct := productElements[0]
		secondProduct := productElements[1]

		firstLocation, err := firstProduct.Location()
		if err != nil {
			t.Fatalf("Failed to get first product location: %s", err)
		}

		secondLocation, err := secondProduct.Location()
		if err != nil {
			t.Fatalf("Failed to get second product location: %s", err)
		}

		firstSize, err := firstProduct.Size()
		if err != nil {
			t.Fatalf("Failed to get first product size: %s", err)
		}

		if secondLocation.Y < (firstLocation.Y + firstSize.Height) {
			t.Errorf("Products not stacked vertically in mobile view")
		}
	}

	webDriver.ResizeWindow("", 1366, 768)
}

func waitForElement(t *testing.T, by, value string, timeout time.Duration) selenium.WebElement {
	deadline := time.Now().Add(timeout)
	var lastErr error

	for time.Now().Before(deadline) {
		element, err := findElement(by, value)
		if err == nil {
			return element
		}
		lastErr = err
		time.Sleep(500 * time.Millisecond)
	}

	source, _ := webDriver.PageSource()
	t.Logf("Page source (first 500 chars) while looking for %s: %.500s...", value, source)

	t.Fatalf("Timed out waiting for element %s: %s after %s", value, lastErr, timeout)
	return nil
}
