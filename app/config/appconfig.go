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
	appConfig.GoogleClientId = "414109645872-g5og4q7pmua0na6sod0jtnvt16mdl4fh.apps.googleusercontent.com"
	appConfig.GoogleClientSecret = "U3KV6G8sYqxa-qtjoxRnk6tX"
	appConfig.Host = "localhost:8080"
	appConfig.StripeKey = "sk_test_4PYk89tUHopayPe2fctjjtuh"
	appConfig.StripePublishableKey = "pk_test_4PYkRwABg9hyrcISOgOgdfJY"

	return appConfig
}

// newProdAppConfig returns the AppConfig for the production environment
func newProdAppConfig() *AppConfig {
	appConfig := new(AppConfig)
	appConfig.GoogleClientId = "414109645872-actrhemoam4scv3b7dqco3b5ai5fov6o.apps.googleusercontent.com"
	appConfig.GoogleClientSecret = "_gTyIGfjBTO7PmiQ6l8jEEE8"
	appConfig.Host = "www.mygluk.it"
	appConfig.StripeKey = "sk_live_4PYk4mJU3wMOncwt0UY5xE2G"
	appConfig.StripePublishableKey = "pk_live_4PYkBwbMIrlFhtf9O5h108Xh"

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
