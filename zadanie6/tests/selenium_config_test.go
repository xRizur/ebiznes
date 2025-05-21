package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/tebeka/selenium"
)

const (
	// URLs for testing
	baseURL    = "https://verbose-carnival-x6pr4pv9pwh6g6v-3000.app.github.dev"
	apiBaseURL = "https://verbose-carnival-x6pr4pv9pwh6g6v-8080.app.github.dev"
)

var (
	webDriver selenium.WebDriver
)

// SetupBrowserStack sets up the WebDriver with BrowserStack capabilities
func SetupBrowserStack(username, accessKey string) (selenium.WebDriver, error) {
	log.Printf("BrowserStack credentials - Username: %s, Access Key length: %d", username, len(accessKey))

	if username == "" || accessKey == "" {
		log.Println("BROWSERSTACK_USERNAME and BROWSERSTACK_ACCESS_KEY must be set")
		return nil, fmt.Errorf("missing BrowserStack credentials")
	}

	caps := selenium.Capabilities{
		"browserName":    "Chrome",
		"browserVersion": "latest",
		"bstack:options": map[string]interface{}{
			"os":              "Windows",
			"osVersion":       "10",
			"resolution":      "1920x1080",
			"projectName":     "Zadanie6",
			"sessionName":     "E-commerce Tests",
			"local":           "false",
			"debug":           "true",
			"consoleLogs":     "verbose",
			"networkLogs":     "true",
			"seleniumVersion": "4.0.0",
			"userName":        username,
			"accessKey":       accessKey,
		},
	}

	log.Println("Connecting to BrowserStack...")
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("https://%s:%s@hub-cloud.browserstack.com/wd/hub", username, accessKey))
	if err != nil {
		log.Printf("Failed to connect to BrowserStack: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to BrowserStack")
	return driver, nil
}

func TestMain(m *testing.M) {
	username := "xrizur_fGCWjX"
	accessKey := "QVLsz6JFNhdxfS5ZVHwB"

	log.Println("Initializing BrowserStack WebDriver...")

	var err error
	webDriver, err = SetupBrowserStack(username, accessKey)
	if err != nil {
		log.Printf("Failed to set up BrowserStack: %v", err)
	} else {
		log.Println("WebDriver successfully initialized")
		webDriver.SetImplicitWaitTimeout(10000)
		defer webDriver.Quit()
	}

	code := m.Run()
	os.Exit(code)
}

// Helper functions for tests
func findElement(by, value string) (selenium.WebElement, error) {
	return webDriver.FindElement(by, value)
}

func findElements(by, value string) ([]selenium.WebElement, error) {
	return webDriver.FindElements(by, value)
}

func assertElementExists(t *testing.T, by, value, description string) selenium.WebElement {
	element, err := findElement(by, value)
	if err != nil {
		t.Fatalf("Failed to find %s: %s", description, err)
	}
	return element
}

func assertElementText(t *testing.T, element selenium.WebElement, expected, description string) {
	text, err := element.Text()
	if err != nil {
		t.Fatalf("Failed to get text from %s: %s", description, err)
	}
	if text != expected {
		t.Errorf("Expected %s to be '%s', got '%s'", description, expected, text)
	}
}

func assertElementCount(t *testing.T, by, value string, expected int, description string) {
	elements, err := findElements(by, value)
	if err != nil {
		t.Fatalf("Failed to find %s: %s", description, err)
	}
	if len(elements) != expected {
		t.Errorf("Expected %d %s, got %d", expected, description, len(elements))
	}
}

func assertPageTitle(t *testing.T, expected string) {
	title, err := webDriver.Title()
	if err != nil {
		t.Fatalf("Failed to get page title: %s", err)
	}
	if title != expected {
		t.Errorf("Expected page title '%s', got '%s'", expected, title)
	}
}
