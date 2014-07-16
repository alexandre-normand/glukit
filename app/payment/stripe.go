// Payment package handles interaction with Stripe for processing donations
package payment

import (
	"appengine"
	"appengine/urlfetch"
	"github.com/alexandre-normand/glukit/app/config"
	"github.com/cosn/stripe"
	"net/http"
)

type StripeClient struct {
	apiKey string
}

func NewStripeClient(appConfig *config.AppConfig) (c *StripeClient) {
	c = new(StripeClient)
	c.apiKey = appConfig.StripeKey
	return c
}

func (c *StripeClient) SubmitDonation(r *http.Request) {
	context := appengine.NewContext(r)
	client := urlfetch.Client(context)

	stripe := &stripe.Client{}
	stripe.Init(c.apiKey, client, nil)
}
