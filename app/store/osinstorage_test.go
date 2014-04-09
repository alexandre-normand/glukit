package store_test

import (
	"appengine/aetest"
	. "github.com/alexandre-normand/glukit/app/store"
	"testing"
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
