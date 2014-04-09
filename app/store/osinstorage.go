package store

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"github.com/RangelReale/osin"
)

type OsinAppEngineStore struct {
	c appengine.Context
}

type oclient struct {
	Id          string `datastore:"Id"`
	Secret      string `datastore:"Secret,noindex"`
	RedirectUri string `datastore:"RedirectUri,noindex"`
}

// TODO: create structs
func NewOsinAppEngineStore(c appengine.Context) *OsinAppEngineStore {
	r := &OsinAppEngineStore{c}

	err := r.addClient(&osin.Client{
		Id:          "834681386231.mygluk.it",
		Secret:      "xEh2sZvNRvYnK9his1S_sdd2MlUc",
		RedirectUri: "urn:ietf:wg:oauth:2.0:oob",
	})

	if err != nil {
		c.Warningf("Failed to initialize oauth server: %v", err)
	}
	return r
}

func (s *OsinAppEngineStore) addClient(c *osin.Client) error {
	s.c.Debugf("AddClient: %s\n", c.Id)
	key := datastore.NewKey(s.c, "osin.client", c.Id, 0, nil)
	oclient := oclient{c.Id, c.Secret, c.RedirectUri}
	_, err := datastore.Put(s.c, key, &oclient)

	if err != nil {
		return err
	}

	return nil
}

func (s *OsinAppEngineStore) GetClient(id string) (*osin.Client, error) {
	s.c.Debugf("GetClient: %s\n", id)
	key := datastore.NewKey(s.c, "osin.client", id, 0, nil)
	client := new(oclient)
	err := datastore.Get(s.c, key, client)

	if err != nil {
		s.c.Warningf("Error looking up client by id [%s]: [%v]", id, err)
		return nil, errors.New("Client not found")
	}
	osinClient := osin.Client{client.Id, client.Secret, client.RedirectUri, nil}
	return &osinClient, nil
}

func (s *OsinAppEngineStore) SaveAuthorize(data *osin.AuthorizeData) error {
	s.c.Debugf("SaveAuthorize: %s\n", data.Code)
	key := datastore.NewKey(s.c, "authorize.data", data.Code, 0, nil)
	_, err := datastore.Put(s.c, key, data)
	if err != nil {
		return err
	}

	return nil
}

func (s *OsinAppEngineStore) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	s.c.Debugf("LoadAuthorize: %s\n", code)
	key := datastore.NewKey(s.c, "authorize.data", code, 0, nil)
	authorizeData := new(osin.AuthorizeData)
	err := datastore.Get(s.c, key, authorizeData)
	if err != nil {
		s.c.Infof("Authorization data not found for code [%s]: %v", code, err)
		return nil, errors.New("Authorize not found")
	}
	return authorizeData, nil
}

func (s *OsinAppEngineStore) RemoveAuthorize(code string) error {
	s.c.Debugf("RemoveAuthorize: %s\n", code)
	key := datastore.NewKey(s.c, "authorize.data", code, 0, nil)
	err := datastore.Delete(s.c, key)
	if err != nil {
		return err
	}

	return nil
}

func (s *OsinAppEngineStore) SaveAccess(data *osin.AccessData) error {
	s.c.Debugf("SaveAccess: %s\n", data.AccessToken)
	key := datastore.NewKey(s.c, "access.data", data.AccessToken, 0, nil)
	_, err := datastore.Put(s.c, key, data)
	if err != nil {
		return err
	}

	if data.RefreshToken != "" {
		key = datastore.NewKey(s.c, "access.refresh", data.RefreshToken, 0, nil)
		_, err := datastore.Put(s.c, key, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *OsinAppEngineStore) LoadAccess(code string) (*osin.AccessData, error) {
	s.c.Debugf("LoadAccess: %s\n", code)
	key := datastore.NewKey(s.c, "access.data", code, 0, nil)
	accessData := new(osin.AccessData)
	err := datastore.Get(s.c, key, accessData)
	if err != nil {
		s.c.Infof("Access data not found for code [%s]: %v", code, err)
		return nil, errors.New("Access data not found")
	}

	return accessData, nil
}

func (s *OsinAppEngineStore) RemoveAccess(code string) error {
	s.c.Debugf("RemoveAccess: %s\n", code)
	key := datastore.NewKey(s.c, "access.data", code, 0, nil)
	err := datastore.Delete(s.c, key)
	if err != nil {
		return err
	}

	return nil
}

func (s *OsinAppEngineStore) LoadRefresh(code string) (*osin.AccessData, error) {
	s.c.Debugf("LoadRefresh: %s\n", code)
	key := datastore.NewKey(s.c, "access.refresh", code, 0, nil)
	accessData := new(osin.AccessData)
	_, err := datastore.Put(s.c, key, accessData)
	if err != nil {
		s.c.Infof("Refresh data not found for code [%s]: %v", code, err)
		errors.New("Refresh not found")
	}
	return accessData, nil
}

func (s *OsinAppEngineStore) RemoveRefresh(code string) error {
	s.c.Debugf("RemoveRefresh: %s\n", code)
	key := datastore.NewKey(s.c, "access.refresh", code, 0, nil)
	err := datastore.Delete(s.c, key)
	if err != nil {
		return err
	}

	return nil
}
