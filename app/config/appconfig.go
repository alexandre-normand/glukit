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
	StripeKey            string
	StripePublishableKey string
}

// newTestAppConfig returns the AppConfig for a test environment
func newTestAppConfig() *AppConfig {
	appConfig := new(AppConfig)
	appConfig.GoogleClientId = "***REMOVED***"
	appConfig.GoogleClientSecret = "***REMOVED***"
	appConfig.Host = "localhost:8080"
	appConfig.StripeKey = "***REMOVED***"
	appConfig.StripePublishableKey = "***REMOVED***"

	return appConfig
}

// newProdAppConfig returns the AppConfig for the production environment
func newProdAppConfig() *AppConfig {
	appConfig := new(AppConfig)
	appConfig.GoogleClientId = "***REMOVED***"
	appConfig.GoogleClientSecret = "***REMOVED***"
	appConfig.Host = "www.mygluk.it"
	appConfig.StripeKey = "***REMOVED***"
	appConfig.StripePublishableKey = "***REMOVED***"

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
