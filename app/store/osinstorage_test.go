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
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	osinStorage := NewOsinAppEngineStore(req)
	_, err = osinStorage.GetClient("***REMOVED***", req)

	if err != nil {
		t.Fatal(err)
	}
}

func TestAccessDataStorage(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	osinStorage := NewOsinAppEngineStore(req)
	client, err := osinStorage.GetClient("***REMOVED***", req)
	d := osin.AccessData{client, nil, nil, "token", "test", 0, "scope", "uri", time.Now(), nil}
	err = osinStorage.SaveAccess(&d, req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = osinStorage.LoadAccess("token", req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthorizeDataStorage(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	osinStorage := NewOsinAppEngineStore(nil)
	client, err := osinStorage.GetClient("***REMOVED***", req)
	d := osin.AuthorizeData{client, "code", 0, "scope", "uri", "state", time.Now(), nil}
	err = osinStorage.SaveAuthorize(&d, req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = osinStorage.LoadAuthorize("code", req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFullAccessDataStorage(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	osinStorage := NewOsinAppEngineStore(req)
	client, err := osinStorage.GetClient("***REMOVED***", req)
	d := osin.AuthorizeData{client, "code", 0, "scope", "uri", "state", time.Now(), nil}
	err = osinStorage.SaveAuthorize(&d, req)
	if err != nil {
		t.Fatal(err)
	}

	previousAccess := osin.AccessData{client, nil, nil, "token", "test", 0, "scope", "uri", time.Now(), nil}
	err = osinStorage.SaveAccess(&previousAccess, req)
	if err != nil {
		t.Fatal(err)
	}

	access := osin.AccessData{client, &d, &previousAccess, "token", "test", 0, "scope", "uri", time.Now(), nil}
	err = osinStorage.SaveAccess(&access, req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = osinStorage.LoadAccess("token", req)
	if err != nil {
		t.Fatal(err)
	}

}
