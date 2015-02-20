// config package wraps configuration accessors
package config

import (
	"appengine"
)

// AppConfig is all global application configuration values
// It has a test mode and a production as per the datastore's appengine
// environment.
type AppConfig struct {
	GoogleClientId       string
	GoogleClientSecret   string
	Host                 string
	SSLHost              string
	StripeKey            string
	StripePublishableKey string
}

// newTestAppConfig returns the AppConfig for a test environment
func newTestAppConfig() *AppConfig {
	appConfig := new(AppConfig)
	appConfig.GoogleClientId = "ENV_TEST_CLIENT_ID"
	appConfig.GoogleClientSecret = "ENV_TEST_CLIENT_SECRET"
	appConfig.Host = "localhost:8080"
	appConfig.SSLHost = "http://localhost:8080"
	appConfig.StripeKey = "ENV_TEST_STRIPE_KEY"
	appConfig.StripePublishableKey = "ENV_TEST_STRIPE_PUBLISHABLE_KEY"

	return appConfig
}

// newProdAppConfig returns the AppConfig for the production environment
func newProdAppConfig() *AppConfig {
	appConfig := new(AppConfig)
	appConfig.GoogleClientId = "ENV_PROD_CLIENT_ID"
	appConfig.GoogleClientSecret = "ENV_PROD_CLIENT_SECRET"
	appConfig.Host = "www.mygluk.it"
	appConfig.SSLHost = "https://glukit.appspot.com"
	appConfig.StripeKey = "ENV_PROD_STRIPE_KEY"
	appConfig.StripePublishableKey = "ENV_PROD_STRIPE_PUBLISHABLE_KEY"

	return appConfig
}

// NewAppConfig returns the AppConfig that matches the current environment (test or prod)
// as returned by appengine.IsDevAppServer()
func NewAppConfig() *AppConfig {
	if appengine.IsDevAppServer() {
		return newTestAppConfig()
	} else {
		return newProdAppConfig()
	}
}
