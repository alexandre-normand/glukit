package store

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"github.com/RangelReale/osin"
	"time"
)

type OsinAppEngineStore struct {
	c appengine.Context
}

type oClient struct {
	Id          string `datastore:"Id"`
	Secret      string `datastore:"Secret,noindex"`
	RedirectUri string `datastore:"RedirectUri,noindex"`
}

// Authorization data
type oAuthorizeData struct {
	Client      *oClient  `datastore:"-,noindex"`
	Code        string    `datastore:"Code"`
	ExpiresIn   int32     `datastore:"ExpiresIn"`
	Scope       string    `datastore:"Scope"`
	RedirectUri string    `datastore:"RedirectUri"`
	State       string    `datastore:"State"`
	CreatedAt   time.Time `datastore:"CreatedAt"`
}

// AccessData
type oAccessData struct {
	Client        *oClient        `datastore:"-,noindex"`
	AuthorizeData *oAuthorizeData `datastore:"-,noindex"`
	AccessData    *oAccessData    `datastore:"-,noindex"`
	AccessToken   string          `datastore:"accessToken"`
	RefreshToken  string          `datastore:"refreshToken"`
	ExpiresIn     int32           `datastore:"expiresIn,noindex"`
	Scope         string          `datastore:"scope,noindex"`
	RedirectUri   string          `datastore:"redirectUri,noindex"`
	CreatedAt     time.Time       `datastore:"createdAt,noindex"`
}

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
	_, err := datastore.Put(s.c, key, newInternalClient(c))

	if err != nil {
		return err
	}

	return nil
}

func newInternalClient(c *osin.Client) *oClient {
	if c == nil {
		return nil
	}
	return &oClient{c.Id, c.Secret, c.RedirectUri}
}

func newOsinClient(c *oClient) *osin.Client {
	if c == nil {
		return nil
	}
	return &osin.Client{c.Id, c.Secret, c.RedirectUri, nil}
}

func (d *oAuthorizeData) Load(c <-chan datastore.Property) error {
	if err := datastore.LoadStruct(d, c); err != nil {
		return err
	}

	return nil
}

func (d *oAuthorizeData) Save(c chan<- datastore.Property) error {
	return datastore.SaveStruct(d, c)
}

func (s *OsinAppEngineStore) GetClient(id string) (*osin.Client, error) {
	s.c.Debugf("GetClient: %s\n", id)
	key := datastore.NewKey(s.c, "osin.client", id, 0, nil)
	client := new(oClient)
	err := datastore.Get(s.c, key, client)

	if err != nil {
		s.c.Warningf("Error looking up client by id [%s]: [%v]", id, err)
		return nil, errors.New("Client not found")
	}
	osinClient := newOsinClient(client)
	return osinClient, nil
}

func newInternalAuthorizeData(d *osin.AuthorizeData) *oAuthorizeData {
	if d == nil {
		return nil
	}
	return &oAuthorizeData{newInternalClient(d.Client), d.Code, d.ExpiresIn, d.Scope, d.RedirectUri, d.State, d.CreatedAt}
}

func newOsinAuthorizeData(d *oAuthorizeData) *osin.AuthorizeData {
	if d == nil {
		return nil
	}
	return &osin.AuthorizeData{newOsinClient(d.Client), d.Code, d.ExpiresIn, d.Scope, d.RedirectUri, d.State, d.CreatedAt, nil}
}

func (s *OsinAppEngineStore) SaveAuthorize(data *osin.AuthorizeData) error {
	s.c.Debugf("SaveAuthorize: %s\n", data.Code)
	key := datastore.NewKey(s.c, "authorize.data", data.Code, 0, nil)
	_, err := datastore.Put(s.c, key, newInternalAuthorizeData(data))
	if err != nil {
		return err
	}

	return nil
}

func (s *OsinAppEngineStore) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	s.c.Debugf("LoadAuthorize: %s\n", code)
	key := datastore.NewKey(s.c, "authorize.data", code, 0, nil)
	authorizeData := new(oAuthorizeData)
	err := datastore.Get(s.c, key, authorizeData)
	if err != nil {
		s.c.Infof("Authorization data not found for code [%s]: %v", code, err)
		return nil, errors.New("Authorize not found")
	}
	return newOsinAuthorizeData(authorizeData), nil
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

func newInternalAccessData(d *osin.AccessData) *oAccessData {
	if d == nil {
		return nil
	}
	return &oAccessData{newInternalClient(d.Client), newInternalAuthorizeData(d.AuthorizeData), newInternalAccessData(d.AccessData), d.AccessToken, d.RefreshToken, d.ExpiresIn, d.Scope, d.RedirectUri, d.CreatedAt}
}

func newOsinAccessData(d *oAccessData) *osin.AccessData {
	if d == nil {
		return nil
	}
	return &osin.AccessData{newOsinClient(d.Client), newOsinAuthorizeData(d.AuthorizeData), newOsinAccessData(d.AccessData), d.AccessToken, d.RefreshToken, d.ExpiresIn, d.Scope, d.RedirectUri, d.CreatedAt, nil}
}

func (s *OsinAppEngineStore) SaveAccess(data *osin.AccessData) error {
	s.c.Debugf("SaveAccess: %s\n", data.AccessToken)
	key := datastore.NewKey(s.c, "access.data", data.AccessToken, 0, nil)
	internalAccessData := newInternalAccessData(data)
	_, err := datastore.Put(s.c, key, internalAccessData)
	if err != nil {
		return err
	}

	if data.RefreshToken != "" {
		key = datastore.NewKey(s.c, "access.refresh", data.RefreshToken, 0, nil)
		_, err := datastore.Put(s.c, key, internalAccessData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *OsinAppEngineStore) LoadAccess(code string) (*osin.AccessData, error) {
	s.c.Debugf("LoadAccess: %s\n", code)
	key := datastore.NewKey(s.c, "access.data", code, 0, nil)
	accessData := new(oAccessData)
	err := datastore.Get(s.c, key, accessData)
	if err != nil {
		s.c.Infof("Access data not found for code [%s]: %v", code, err)
		return nil, errors.New("Access data not found")
	}

	return newOsinAccessData(accessData), nil
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
	accessData := new(oAccessData)
	_, err := datastore.Put(s.c, key, accessData)
	if err != nil {
		s.c.Infof("Refresh data not found for code [%s]: %v", code, err)
		errors.New("Refresh not found")
	}
	return newOsinAccessData(accessData), nil
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
