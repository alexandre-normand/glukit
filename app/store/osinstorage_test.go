package store_test

import (
	"appengine/aetest"
	"github.com/RangelReale/osin"
	. "github.com/alexandre-normand/glukit/app/store"
	"testing"
	"time"
)

func TestGetClient(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	osinStorage := NewOsinAppEngineStore(c)
	_, err = osinStorage.GetClient("***REMOVED***")

	if err != nil {
		t.Fatal(err)
	}
}

func TestAccessDataStorage(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	osinStorage := NewOsinAppEngineStore(c)
	client, err := osinStorage.GetClient("***REMOVED***")
	d := osin.AccessData{client, nil, nil, "token", "test", 0, "scope", "uri", time.Now(), nil}
	err = osinStorage.SaveAccess(&d)
	if err != nil {
		t.Fatal(err)
	}

	_, err = osinStorage.LoadAccess("token")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthorizeDataStorage(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	osinStorage := NewOsinAppEngineStore(c)
	client, err := osinStorage.GetClient("***REMOVED***")
	d := osin.AuthorizeData{client, "code", 0, "scope", "uri", "state", time.Now(), nil}
	err = osinStorage.SaveAuthorize(&d)
	if err != nil {
		t.Fatal(err)
	}

	_, err = osinStorage.LoadAuthorize("code")
	if err != nil {
		t.Fatal(err)
	}
}
