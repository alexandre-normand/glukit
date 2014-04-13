package store

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"github.com/alexandre-normand/osin"
	"net/http"
	"time"
)

type OsinAppEngineStore struct {
}

type oClient struct {
	Id          string `datastore:"Id"`
	Secret      string `datastore:"Secret,noindex"`
	RedirectUri string `datastore:"RedirectUri,noindex"`
	UserData    string `datastore:"UserData,noindex"`
}

// Authorization data
type oAuthorizeData struct {
	ClientId    string    `datastore:"ClientId,noindex"`
	Code        string    `datastore:"Code"`
	ExpiresIn   int32     `datastore:"ExpiresIn"`
	Scope       string    `datastore:"Scope"`
	RedirectUri string    `datastore:"RedirectUri"`
	State       string    `datastore:"State"`
	CreatedAt   time.Time `datastore:"CreatedAt"`
	UserData    string    `datastore:"UserData,noindex"`
}

// AccessData
type oAccessData struct {
	ClientId          string    `datastore:"ClientId,noindex"`
	AuthorizeDataCode string    `datastore:"AuthorizeDataCode,noindex"`
	AccessDataToken   string    `datastore:"AccessDataToken,noindex"`
	AccessToken       string    `datastore:"AccessToken"`
	RefreshToken      string    `datastore:"RefreshToken"`
	ExpiresIn         int32     `datastore:"ExpiresIn,noindex"`
	Scope             string    `datastore:"Scope,noindex"`
	RedirectUri       string    `datastore:"RedirectUri,noindex"`
	CreatedAt         time.Time `datastore:"CreatedAt,noindex"`
	UserData          string    `datastore:"UserData,noindex"`
}

func NewOsinAppEngineStore(r *http.Request) *OsinAppEngineStore {
	c := appengine.NewContext(r)
	s := &OsinAppEngineStore{}

	err := s.addClient(&osin.Client{
		Id:          "834681386231.mygluk.it",
		Secret:      "xEh2sZvNRvYnK9his1S_sdd2MlUc",
		RedirectUri: "urn:ietf:wg:oauth:2.0:oob",
	}, r)

	if err != nil {
		c.Warningf("Failed to initialize oauth server: %v", err)
	}
	return s
}

func (s *OsinAppEngineStore) addClient(c *osin.Client, r *http.Request) error {
	context := appengine.NewContext(r)
	context.Debugf("AddClient: %s\n", c.Id)
	key := datastore.NewKey(context, "osin.client", c.Id, 0, nil)
	_, err := datastore.Put(context, key, newInternalClient(c))

	if err != nil {
		return err
	}

	return nil
}

func newInternalClient(c *osin.Client) *oClient {
	if c == nil {
		return nil
	}
	return &oClient{c.Id, c.Secret, c.RedirectUri, c.UserData.(string)}
}

func newOsinClient(c *oClient) *osin.Client {
	if c == nil {
		return nil
	}
	return &osin.Client{c.Id, c.Secret, c.RedirectUri, c.UserData}
}

func (s *OsinAppEngineStore) GetClient(id string, r *http.Request) (*osin.Client, error) {
	context := appengine.NewContext(r)
	context.Debugf("GetClient: %s\n", id)
	key := datastore.NewKey(context, "osin.client", id, 0, nil)
	client := new(oClient)
	err := datastore.Get(context, key, client)

	if err != nil {
		context.Warningf("Error looking up client by id [%s]: [%v]", id, err)
		return nil, errors.New("Client not found")
	}
	osinClient := newOsinClient(client)
	return osinClient, nil
}

func newInternalAuthorizeData(d *osin.AuthorizeData) *oAuthorizeData {
	if d == nil {
		return nil
	}

	clientId := ""
	if client := d.Client; client != nil {
		clientId = client.Id
	}

	return &oAuthorizeData{clientId, d.Code, d.ExpiresIn, d.Scope, d.RedirectUri, d.State, d.CreatedAt, d.UserData.(string)}
}

func newOsinAuthorizeData(d *oAuthorizeData, c *osin.Client) *osin.AuthorizeData {
	if d == nil {
		return nil
	}
	return &osin.AuthorizeData{c, d.Code, d.ExpiresIn, d.Scope, d.RedirectUri, d.State, d.CreatedAt, d.UserData}
}

func (s *OsinAppEngineStore) SaveAuthorize(data *osin.AuthorizeData, r *http.Request) error {
	context := appengine.NewContext(r)
	context.Debugf("SaveAuthorize: %s\n", data.Code)
	key := datastore.NewKey(context, "authorize.data", data.Code, 0, nil)
	_, err := datastore.Put(context, key, newInternalAuthorizeData(data))
	if err != nil {
		context.Warningf("Error saving authorize data [%s]: [%v]", data.Code, err)
		return err
	}

	return nil
}

func (s *OsinAppEngineStore) LoadAuthorize(code string, r *http.Request) (*osin.AuthorizeData, error) {
	context := appengine.NewContext(r)
	context.Debugf("LoadAuthorize: %s\n", code)
	key := datastore.NewKey(context, "authorize.data", code, 0, nil)
	authorizeData := new(oAuthorizeData)
	err := datastore.Get(context, key, authorizeData)
	if err != nil {
		context.Infof("Authorization data not found for code [%s]: %v", code, err)
		return nil, errors.New("Authorize not found")
	}

	var c *osin.Client
	if authorizeData.ClientId != "" {
		c, err = s.GetClient(authorizeData.ClientId, r)
		if err != nil {
			context.Infof("Authorization data can't load client with id [%s]: %v", authorizeData.ClientId, err)
			return nil, errors.New("Client for AuthorizeData not found")
		}
	}

	return newOsinAuthorizeData(authorizeData, c), nil
}

func (s *OsinAppEngineStore) RemoveAuthorize(code string, r *http.Request) error {
	context := appengine.NewContext(r)
	context.Debugf("RemoveAuthorize: %s\n", code)
	key := datastore.NewKey(context, "authorize.data", code, 0, nil)
	err := datastore.Delete(context, key)
	if err != nil {
		return err
	}

	return nil
}

func newInternalAccessData(d *osin.AccessData) *oAccessData {
	if d == nil {
		return nil
	}

	clientId := ""
	if client := d.Client; client != nil {
		clientId = client.Id
	}

	authCode := ""
	if authorizeData := d.AuthorizeData; authorizeData != nil {
		authCode = authorizeData.Code
	}

	accessToken := ""
	if accessData := d.AccessData; accessData != nil {
		accessToken = accessData.AccessToken
	}
	return &oAccessData{clientId, authCode, accessToken, d.AccessToken, d.RefreshToken, d.ExpiresIn, d.Scope, d.RedirectUri, d.CreatedAt, d.UserData.(string)}
}

func newOsinAccessData(d *oAccessData, c *osin.Client, authData *osin.AuthorizeData, accessData *osin.AccessData) *osin.AccessData {
	if d == nil {
		return nil
	}
	return &osin.AccessData{c, authData, accessData, d.AccessToken, d.RefreshToken, d.ExpiresIn, d.Scope, d.RedirectUri, d.CreatedAt, d.UserData}
}

func (s *OsinAppEngineStore) SaveAccess(data *osin.AccessData, r *http.Request) error {
	context := appengine.NewContext(r)
	context.Debugf("SaveAccess: %s\n", data.AccessToken)
	key := datastore.NewKey(context, "access.data", data.AccessToken, 0, nil)
	internalAccessData := newInternalAccessData(data)
	_, err := datastore.Put(context, key, internalAccessData)
	if err != nil {
		return err
	}

	if data.RefreshToken != "" {
		key = datastore.NewKey(context, "access.refresh", data.RefreshToken, 0, nil)
		_, err := datastore.Put(context, key, internalAccessData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *OsinAppEngineStore) LoadAccess(code string, r *http.Request) (*osin.AccessData, error) {
	context := appengine.NewContext(r)
	context.Debugf("LoadAccess: %s\n", code)
	key := datastore.NewKey(context, "access.data", code, 0, nil)
	accessData := new(oAccessData)
	err := datastore.Get(context, key, accessData)
	if err != nil {
		context.Infof("Access data not found for code [%s]: %v", code, err)
		return nil, errors.New("Access data not found")
	}

	var c *osin.Client
	if accessData.ClientId != "" {
		c, err = s.GetClient(accessData.ClientId, r)
		if err != nil {
			context.Infof("Access data can't load client with id [%s]: %v", accessData.ClientId, err)
			return nil, errors.New("Client for AccessData not found")
		}
	}

	var innerAuthData *osin.AuthorizeData
	if accessData.AuthorizeDataCode != "" {
		innerAuthData, err = s.LoadAuthorize(accessData.AuthorizeDataCode, r)
		if err != nil {
			context.Infof("Inner authorize data can't load client with code [%s]: %v", accessData.AuthorizeDataCode, err)
			return nil, errors.New("Inner AuthorizeData for AccessData not found")
		}
	}

	var innerAccessData *osin.AccessData
	if accessData.AccessDataToken != "" {
		innerAccessData, err = s.LoadAccess(accessData.AccessDataToken, r)
		if err != nil {
			context.Infof("Inner access data can't load client with token [%s]: %v", accessData.AccessDataToken, err)
			return nil, errors.New("Inner AccessData for AccessData not found")
		}
	}

	return newOsinAccessData(accessData, c, innerAuthData, innerAccessData), nil
}

func (s *OsinAppEngineStore) RemoveAccess(code string, r *http.Request) error {
	context := appengine.NewContext(r)
	context.Debugf("RemoveAccess: %s\n", code)
	key := datastore.NewKey(context, "access.data", code, 0, nil)
	err := datastore.Delete(context, key)
	if err != nil {
		return err
	}

	return nil
}

func (s *OsinAppEngineStore) LoadRefresh(code string, r *http.Request) (*osin.AccessData, error) {
	context := appengine.NewContext(r)
	context.Debugf("LoadRefresh: %s\n", code)
	key := datastore.NewKey(context, "access.refresh", code, 0, nil)
	accessData := new(oAccessData)
	_, err := datastore.Put(context, key, accessData)
	if err != nil {
		context.Infof("Refresh data not found for code [%s]: %v", code, err)
		errors.New("Refresh not found")
	}

	var c *osin.Client
	if accessData.ClientId != "" {
		c, err = s.GetClient(accessData.ClientId, r)
		if err != nil {
			context.Infof("Access data can't load client with id [%s]: %v", accessData.ClientId, err)
			return nil, errors.New("Client for AccessData not found")
		}
	}

	var innerAuthData *osin.AuthorizeData
	if accessData.AuthorizeDataCode != "" {
		innerAuthData, err = s.LoadAuthorize(accessData.AuthorizeDataCode, r)
		if err != nil {
			context.Infof("Inner authorize data can't load client with code [%s]: %v", accessData.AuthorizeDataCode, err)
			return nil, errors.New("Inner AuthorizeData for AccessData not found")
		}
	}

	var innerAccessData *osin.AccessData
	if accessData.AccessDataToken != "" {
		innerAccessData, err = s.LoadAccess(accessData.AccessDataToken, r)
		if err != nil {
			context.Infof("Inner access data can't load client with token [%s]: %v", accessData.AccessDataToken, err)
			return nil, errors.New("Inner AccessData for AccessData not found")
		}
	}

	return newOsinAccessData(accessData, c, innerAuthData, innerAccessData), nil
}

func (s *OsinAppEngineStore) RemoveRefresh(code string, r *http.Request) error {
	context := appengine.NewContext(r)
	context.Debugf("RemoveRefresh: %s\n", code)
	key := datastore.NewKey(context, "access.refresh", code, 0, nil)
	err := datastore.Delete(context, key)
	if err != nil {
		return err
	}

	return nil
}
