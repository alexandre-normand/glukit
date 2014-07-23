package payment_test

import (
	"appengine/aetest"
	"appengine/urlfetch"
	"appengine/user"
	"github.com/alexandre-normand/glukit/app/config"
	. "github.com/alexandre-normand/glukit/app/payment"
	"github.com/cosn/stripe"
	"strconv"
	"strings"
	"testing"
)

func TestDonationWithNoUser(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	appConfig := config.NewAppConfig()
	client := NewStripeClient(appConfig)

	err = client.SubmitDonation(c, "testtoken", "100")
	if err == nil {
		t.Error("TestDonationWithNoUser should have failed")
	}

	if !strings.Contains(err.Error(), "User is nil") {
		t.Errorf("TestDonationWithNoUser should have failed with a \"User is nil\" error but failed with [%s]", err.Error())
	}
}

func TestDonationWithInvalidAmount(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	c.Login(&user.User{Email: "test@glukit.com", AuthDomain: "glukit.com", Admin: false})
	appConfig := config.NewAppConfig()
	client := NewStripeClient(appConfig)

	err = client.SubmitDonation(c, "invalidToken", "invalidVal")
	if err == nil {
		t.Error("TestDonationWithInvalidAmount should have failed")
	}

	if _, ok := err.(*strconv.NumError); !ok {
		t.Errorf("TestDonationWithInvalidAmount should have failed with a NumError but failed with [%v]", err)
	}
}

func TestDonationWithInvalidToken(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	c.Login(&user.User{Email: "test@glukit.com", AuthDomain: "glukit.com", Admin: false})
	appConfig := config.NewAppConfig()
	client := NewStripeClient(appConfig)

	err = client.SubmitDonation(c, "invalidToken", "100")
	if err == nil {
		t.Error("TestDonationWithInvalidToken should have failed")
	}

	if !strings.Contains(err.Error(), "invalidToken") {
		t.Errorf("TestDonationWithInvalidToken should have failed with \"invalidToken\" but failed with [%s]", err.Error())
	}
}

func TestValidDonation(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	c.Login(&user.User{Email: "test@glukit.com", AuthDomain: "glukit.com", Admin: false})
	appConfig := config.NewAppConfig()
	client := NewStripeClient(appConfig)

	httpClient := urlfetch.Client(c)
	sc := &stripe.Client{}
	sc.Init(appConfig.StripeKey, httpClient, nil)
	tokenParams := &stripe.TokenParams{Card: &stripe.CardParams{Number: "4242 4242 4242 4242", CVC: "1234", Month: "12", Year: "2030"}}
	token, err := sc.Tokens.Create(tokenParams)
	if err != nil {
		t.Errorf("TestValidDonation needs a token to test but couldn't get one from stripe: [%v]", err)
	}

	err = client.SubmitDonation(c, token.Id, "100")

	if err != nil {
		t.Errorf("TestValidDonation should have succeeded but failed with error: [%v]", err)
	}
}
