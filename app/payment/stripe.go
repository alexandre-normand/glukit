// Payment package handles interaction with Stripe for processing donations
package payment

import (
	"appengine"
	"appengine/urlfetch"
	"appengine/user"
	"fmt"
	"github.com/alexandre-normand/glukit/app/config"
	"github.com/cosn/stripe"
	"strconv"
)

const (
	STRIPE_TOKEN    = "stripeToken"
	DONATION_AMOUNT = "donation-amount"
)

type StripeClient struct {
	apiKey string
}

func NewStripeClient(appConfig *config.AppConfig) (c *StripeClient) {
	c = new(StripeClient)
	c.apiKey = appConfig.StripeKey
	return c
}

func (c *StripeClient) SubmitDonation(ctx appengine.Context, token string, amountInCentsVal string) error {
	user := user.Current(ctx)
	email := ""
	if user != nil {
		email = user.Email
	}

	client := urlfetch.Client(ctx)

	sc := &stripe.Client{}
	sc.Init(c.apiKey, client, nil)

	amountInCents, err := strconv.ParseUint(amountInCentsVal, 10, 32)
	if err != nil {
		return err
	}

	amountInDollars := float32(amountInCents) / 100.

	ctx.Debugf("Received donation request of [%.2f] with stripe token [%s] for user [%s]", amountInDollars, token, email)

	params := &stripe.ChargeParams{
		Token:    token,
		Amount:   amountInCents,
		Currency: stripe.USD,
		Desc:     fmt.Sprintf("Generous donation of $%.2f to Glukit.", amountInDollars),
		Meta: map[string]string{
			"email": email,
		},
		Statement: "Donation",
		Email:     email,
	}

	if charge, err := sc.Charges.Create(params); err != nil {
		return err
	} else {
		ctx.Infof("Charged donation of [%.2f] to [%s] successfully: [%v]", amountInDollars, email, charge)
	}

	return nil
}
