// config package wraps configuration accessors
package config

import (
	"appengine"
	"github.com/alexandre-normand/glukit/app/secrets"
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
func newTestAppConfig(appSecrets *secrets.AppSecrets) *AppConfig {
	appConfig := new(AppConfig)
	appConfig.GoogleClientId = appSecrets.LocalGoogleClientId
	appConfig.GoogleClientSecret = appSecrets.LocalGoogleClientSecret
	appConfig.Host = "localhost:8080"
	appConfig.SSLHost = "http://localhost:8080"
	appConfig.StripeKey = appSecrets.LocalStripeKey
	appConfig.StripePublishableKey = appSecrets.LocalStripePublishableKey

	return appConfig
}

// newProdAppConfig returns the AppConfig for the production environment
func newProdAppConfig(appSecrets *secrets.AppSecrets) *AppConfig {
	appConfig := new(AppConfig)
	appConfig.GoogleClientId = appSecrets.ProdGoogleClientId
	appConfig.GoogleClientSecret = appSecrets.ProdGoogleClientSecret
	appConfig.Host = "www.mygluk.it"
	appConfig.SSLHost = "https://glukit.appspot.com"
	appConfig.StripeKey = appSecrets.ProdStripeKey
	appConfig.StripePublishableKey = appSecrets.ProdStripePublishableKey

	return appConfig
}

// NewAppConfig returns the AppConfig that matches the current environment (test or prod)
// as returned by appengine.IsDevAppServer()
func NewAppConfig() *AppConfig {
	appSecrets := secrets.NewAppSecrets()
	if appengine.IsDevAppServer() {
		return newTestAppConfig(appSecrets)
	} else {
		return newProdAppConfig(appSecrets)
	}
}
