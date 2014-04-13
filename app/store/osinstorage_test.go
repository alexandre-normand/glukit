package store_test

import (
	"appengine/aetest"
	. "github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/osin"
	"net/http"
	"testing"
	"time"
)

func TestGetClient(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	osinStorage := NewOsinAppEngineStore(c.Request().(*http.Request))
	_, err = osinStorage.GetClient("***REMOVED***", c.Request().(*http.Request))

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

	osinStorage := NewOsinAppEngineStore(c.Request().(*http.Request))
	client, err := osinStorage.GetClient("***REMOVED***", c.Request().(*http.Request))
	d := osin.AccessData{client, nil, nil, "token", "test", 0, "scope", "uri", time.Now(), nil}
	err = osinStorage.SaveAccess(&d, c.Request().(*http.Request))
	if err != nil {
		t.Fatal(err)
	}

	_, err = osinStorage.LoadAccess("token", c.Request().(*http.Request))
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

	osinStorage := NewOsinAppEngineStore(nil)
	client, err := osinStorage.GetClient("***REMOVED***", c.Request().(*http.Request))
	d := osin.AuthorizeData{client, "code", 0, "scope", "uri", "state", time.Now(), nil}
	err = osinStorage.SaveAuthorize(&d, c.Request().(*http.Request))
	if err != nil {
		t.Fatal(err)
	}

	_, err = osinStorage.LoadAuthorize("code", c.Request().(*http.Request))
	if err != nil {
		t.Fatal(err)
	}
}

func TestFullAccessDataStorage(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	osinStorage := NewOsinAppEngineStore(c.Request().(*http.Request))
	client, err := osinStorage.GetClient("***REMOVED***", c.Request().(*http.Request))
	d := osin.AuthorizeData{client, "code", 0, "scope", "uri", "state", time.Now(), nil}
	err = osinStorage.SaveAuthorize(&d, c.Request().(*http.Request))
	if err != nil {
		t.Fatal(err)
	}

	previousAccess := osin.AccessData{client, nil, nil, "token", "test", 0, "scope", "uri", time.Now(), nil}
	err = osinStorage.SaveAccess(&previousAccess, c.Request().(*http.Request))
	if err != nil {
		t.Fatal(err)
	}

	access := osin.AccessData{client, &d, &previousAccess, "token", "test", 0, "scope", "uri", time.Now(), nil}
	err = osinStorage.SaveAccess(&access, c.Request().(*http.Request))
	if err != nil {
		t.Fatal(err)
	}

	_, err = osinStorage.LoadAccess("token", c.Request().(*http.Request))
	if err != nil {
		t.Fatal(err)
	}

}
