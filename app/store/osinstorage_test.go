package store_test

import (
	"appengine/aetest"
	. "github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/osin"
	"testing"
	"time"
)

func TestGetClient(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	t.Logf("Request is [%v] from context [%v]", c.Request(), c)

	osinStorage := NewOsinAppEngineStoreWithContext(c)
	_, err = osinStorage.GetClientWithContext("834681386231.mygluk.it", c)

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
	t.Logf("Request is [%v] from context [%v]", c.Request(), c)

	osinStorage := NewOsinAppEngineStoreWithContext(c)
	client, err := osinStorage.GetClientWithContext("834681386231.mygluk.it", c)
	d := osin.AccessData{client, nil, nil, "token", "test", 0, "scope", "uri", time.Now(), TEST_USER}
	err = osinStorage.SaveAccessWithContext(&d, c)
	if err != nil {
		t.Fatal(err)
	}

	_, err = osinStorage.LoadAccessWithContext("token", c)
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

	osinStorage := NewOsinAppEngineStoreWithContext(c)
	client, err := osinStorage.GetClientWithContext("834681386231.mygluk.it", c)
	d := osin.AuthorizeData{client, "code", 0, "scope", "uri", "state", time.Now(), TEST_USER}
	err = osinStorage.SaveAuthorizeWithContext(&d, c)
	if err != nil {
		t.Fatal(err)
	}

	_, err = osinStorage.LoadAuthorizeWithContext("code", c)
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

	osinStorage := NewOsinAppEngineStoreWithContext(c)
	client, err := osinStorage.GetClientWithContext("834681386231.mygluk.it", c)
	d := osin.AuthorizeData{client, "code", 0, "scope", "uri", "state", time.Now(), TEST_USER}
	err = osinStorage.SaveAuthorizeWithContext(&d, c)
	if err != nil {
		t.Fatal(err)
	}

	previousAccess := osin.AccessData{client, nil, nil, "token", "test", 0, "scope", "uri", time.Now(), TEST_USER}
	err = osinStorage.SaveAccessWithContext(&previousAccess, c)
	if err != nil {
		t.Fatal(err)
	}

	access := osin.AccessData{client, &d, &previousAccess, "token", "test", 0, "scope", "uri", time.Now(), TEST_USER}
	err = osinStorage.SaveAccessWithContext(&access, c)
	if err != nil {
		t.Fatal(err)
	}

	_, err = osinStorage.LoadAccessWithContext("token", c)
	if err != nil {
		t.Fatal(err)
	}

}
