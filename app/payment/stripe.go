// Payment package handles interaction with Stripe for processing donations
package payment

import (
	"appengine"
	"appengine/urlfetch"
	"appengine/user"
	"errors"
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
	if user == nil {
		return errors.New("User is nil, can't proceed with donation.")
	}

	client := urlfetch.Client(ctx)

	sc := &stripe.Client{}
	sc.Init(c.apiKey, client, nil)

	amountInCents, err := strconv.ParseUint(amountInCentsVal, 10, 32)
	if err != nil {
		return err
	}

	amountInDollars := float32(amountInCents) / 100.

	ctx.Debugf("Received donation request of [%.2f] with stripe token [%s] for user [%s]", amountInDollars, token, user)

	params := &stripe.ChargeParams{
		Token:    token,
		Amount:   amountInCents,
		Currency: stripe.USD,
		Desc:     fmt.Sprintf("Generous donation of $%.2f to Glukit.", amountInDollars),
		Meta: map[string]string{
			"email": user.Email,
		},
		Statement: "Donation",
		Email:     user.Email,
	}

	if charge, err := sc.Charges.Create(params); err != nil {
		return err
	} else {
		ctx.Infof("Charged donation of [%.2f] to [%s] successfully: [%v]", amountInDollars, user.Email, charge)
	}

	return nil
}
