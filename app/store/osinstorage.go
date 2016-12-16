package store

import (
	"errors"
	"github.com/alexandre-normand/glukit/app/secrets"
	"github.com/alexandre-normand/osin"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"time"
)

var appSecrets = secrets.NewAppSecrets()

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

func NewOsinAppEngineStoreWithRequest(r *http.Request) *OsinAppEngineStore {
	c := appengine.NewContext(r)

	return NewOsinAppEngineStoreWithContext(c)
}

func NewOsinAppEngineStoreWithContext(c context.Context) *OsinAppEngineStore {
	s := &OsinAppEngineStore{}

	err := s.addClient(&osin.Client{
		Id:          appSecrets.GlukloaderClientId,
		Secret:      appSecrets.GlukloaderClientSecret,
		RedirectUri: "urn:ietf:wg:oauth:2.0:oob",
		UserData:    "",
	}, c)

	if err != nil {
		log.Warningf(c, "Failed to initialize oauth server: %v", err)
	}

	err = s.addClient(&osin.Client{
		Id:          appSecrets.GlukloaderShareEditionClientId,
		Secret:      appSecrets.GlukloaderShareEditionClientSecret,
		RedirectUri: "x-glukloader://oauth/callback",
		UserData:    "",
	}, c)

	if err != nil {
		log.Warningf(c, "Failed to initialize oauth server: %v", err)
	}

	err = s.addClient(&osin.Client{
		Id:          appSecrets.PostManClientId,
		Secret:      appSecrets.PostManClientSecret,
		RedirectUri: "https://www.getpostman.com/oauth2/callback",
		UserData:    "",
	}, c)

	if err != nil {
		log.Warningf(c, "Failed to initialize oauth server: %v", err)
	}

	err = s.addClient(&osin.Client{
		Id:          appSecrets.SimpleClientId,
		Secret:      appSecrets.SimpleClientSecret,
		RedirectUri: "http://localhost:14000/appauth/code",
		UserData:    "",
	}, c)

	if err != nil {
		log.Warningf(c, "Failed to initialize oauth server: %v", err)
	}

	err = s.addClient(&osin.Client{
		Id:          appSecrets.ChromadexClientId,
		Secret:      appSecrets.ChromadexClientSecret,
		RedirectUri: "https://aeapkfdflpgdigehfhjpgccjodakkjje.chromiumapp.org/provider_cb",
		UserData:    "",
	}, c)

	if err != nil {
		log.Warningf(c, "Failed to initialize oauth server: %v", err)
	}

	return s
}

func (s *OsinAppEngineStore) addClient(c *osin.Client, context context.Context) error {
	log.Debugf(context, "AddClient: %s...\n", c.Id)
	key := datastore.NewKey(context, "osin.client", c.Id, 0, nil)
	client := new(oClient)
	err := datastore.Get(context, key, client)

	if err == nil || err != datastore.ErrNoSuchEntity {
		log.Debugf(context, "Client [%s] already stored, skipping.\n", c.Id)
		return nil
	}

	_, err = datastore.Put(context, key, newInternalClient(c))

	if err != nil {
		log.Warningf(context, "Error storing client [%s]: %v", c.Id, err)
		return err
	}

	log.Debugf(context, "Stored new client [%s] successfully.\n", c.Id)
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

	return s.GetClientWithContext(id, context)
}

func (s *OsinAppEngineStore) GetClientWithContext(id string, context context.Context) (*osin.Client, error) {
	log.Debugf(context, "GetClient: %s\n", id)
	key := datastore.NewKey(context, "osin.client", id, 0, nil)
	client := new(oClient)
	err := datastore.Get(context, key, client)

	if err != nil {
		log.Warningf(context, "Error looking up client by id [%s]: [%v]", id, err)
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
	return s.SaveAuthorizeWithContext(data, context)
}

func (s *OsinAppEngineStore) SaveAuthorizeWithContext(data *osin.AuthorizeData, context context.Context) error {
	log.Debugf(context, "SaveAuthorize: %s\n", data.Code)
	key := datastore.NewKey(context, "authorize.data", data.Code, 0, nil)
	_, err := datastore.Put(context, key, newInternalAuthorizeData(data))
	if err != nil {
		log.Warningf(context, "Error saving authorize data [%s]: [%v]", data.Code, err)
		return err
	}

	return nil
}

func (s *OsinAppEngineStore) LoadAuthorize(code string, r *http.Request) (*osin.AuthorizeData, error) {
	context := appengine.NewContext(r)
	return s.LoadAuthorizeWithContext(code, context)
}

func (s *OsinAppEngineStore) LoadAuthorizeWithContext(code string, context context.Context) (*osin.AuthorizeData, error) {
	log.Debugf(context, "LoadAuthorize: %s\n", code)
	key := datastore.NewKey(context, "authorize.data", code, 0, nil)
	authorizeData := new(oAuthorizeData)
	err := datastore.Get(context, key, authorizeData)
	if err != nil {
		log.Infof(context, "Authorization data not found for code [%s]: %v", code, err)
		return nil, errors.New("Authorize not found")
	}

	var c *osin.Client
	if authorizeData.ClientId != "" {
		c, err = s.GetClientWithContext(authorizeData.ClientId, context)
		if err != nil {
			log.Infof(context, "Authorization data can't load client with id [%s]: %v", authorizeData.ClientId, err)
			return nil, errors.New("Client for AuthorizeData not found")
		}
	}

	return newOsinAuthorizeData(authorizeData, c), nil
}

func (s *OsinAppEngineStore) RemoveAuthorize(code string, r *http.Request) error {
	context := appengine.NewContext(r)
	return s.RemoveAuthorizeWithContext(code, context)
}

func (s *OsinAppEngineStore) RemoveAuthorizeWithContext(code string, context context.Context) error {
	log.Debugf(context, "RemoveAuthorize: %s\n", code)
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
	return s.SaveAccessWithContext(data, context)
}

func (s *OsinAppEngineStore) SaveAccessWithContext(data *osin.AccessData, context context.Context) error {
	log.Debugf(context, "SaveAccess [%s]: [%v]\n", data.AccessToken, data)
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
	return s.LoadAccessWithContext(code, context)
}

func (s *OsinAppEngineStore) LoadAccessWithContext(code string, context context.Context) (*osin.AccessData, error) {
	log.Debugf(context, "LoadAccess: %s\n", code)
	key := datastore.NewKey(context, "access.data", code, 0, nil)
	accessData := new(oAccessData)
	err := datastore.Get(context, key, accessData)
	if err != nil {
		log.Infof(context, "Access data not found for code [%s]: %v", code, err)
		return nil, errors.New("Access data not found")
	}

	var c *osin.Client
	if accessData.ClientId != "" {
		c, err = s.GetClientWithContext(accessData.ClientId, context)
		if err != nil {
			log.Infof(context, "Access data can't load client with id [%s]: %v", accessData.ClientId, err)
			return nil, errors.New("Client for AccessData not found")
		}
	}

	var innerAuthData *osin.AuthorizeData
	if accessData.AuthorizeDataCode != "" {
		innerAuthData, err = s.LoadAuthorizeWithContext(accessData.AuthorizeDataCode, context)
		if err != nil {
			log.Infof(context, "Inner authorize data can't be loaded with code [%s]: %v", accessData.AuthorizeDataCode, err)
		}
	}

	var innerAccessData *osin.AccessData
	if accessData.AccessDataToken != "" {
		innerAccessData, err = s.LoadAccessWithContext(accessData.AccessDataToken, context)
		if err != nil {
			log.Infof(context, "Inner access data can't be loaded with token [%s]: %v", accessData.AccessDataToken, err)
		}
	}

	return newOsinAccessData(accessData, c, innerAuthData, innerAccessData), nil
}

func (s *OsinAppEngineStore) RemoveAccess(code string, r *http.Request) error {
	context := appengine.NewContext(r)
	return s.RemoveAccessWithContext(code, context)
}

func (s *OsinAppEngineStore) RemoveAccessWithContext(code string, context context.Context) error {
	log.Debugf(context, "RemoveAccess: %s\n", code)
	key := datastore.NewKey(context, "access.data", code, 0, nil)
	err := datastore.Delete(context, key)
	if err != nil {
		return err
	}

	return nil
}

func (s *OsinAppEngineStore) LoadRefresh(code string, r *http.Request) (*osin.AccessData, error) {
	context := appengine.NewContext(r)
	return s.LoadRefreshWithContext(code, context)
}

func (s *OsinAppEngineStore) LoadRefreshWithContext(code string, context context.Context) (*osin.AccessData, error) {
	log.Debugf(context, "LoadRefresh: %s\n", code)
	key := datastore.NewKey(context, "access.refresh", code, 0, nil)
	accessData := new(oAccessData)
	err := datastore.Get(context, key, accessData)
	if err != nil {
		log.Infof(context, "Refresh data not found for code [%s]: %v", code, err)
		errors.New("Refresh not found")
	}

	var c *osin.Client
	if accessData.ClientId != "" {
		c, err = s.GetClientWithContext(accessData.ClientId, context)
		if err != nil {
			log.Infof(context, "Access data can't load client with id [%s]: %v", accessData.ClientId, err)
			return nil, errors.New("Client for AccessData not found")
		}
	}

	var innerAuthData *osin.AuthorizeData
	if accessData.AuthorizeDataCode != "" {
		innerAuthData, err = s.LoadAuthorizeWithContext(accessData.AuthorizeDataCode, context)
		if err != nil {
			log.Infof(context, "Inner authorize data can't load client with code [%s]: %v", accessData.AuthorizeDataCode, err)
		}
	}

	var innerAccessData *osin.AccessData
	if accessData.AccessDataToken != "" {
		innerAccessData, err = s.LoadAccessWithContext(accessData.AccessDataToken, context)
		if err != nil {
			log.Infof(context, "Inner access data can't load client with token [%s]: %v", accessData.AccessDataToken, err)
		}
	}

	return newOsinAccessData(accessData, c, innerAuthData, innerAccessData), nil
}

func (s *OsinAppEngineStore) RemoveRefresh(code string, r *http.Request) error {
	context := appengine.NewContext(r)
	return s.RemoveRefreshWithContext(code, context)
}

func (s *OsinAppEngineStore) RemoveRefreshWithContext(code string, context context.Context) error {
	log.Debugf(context, "RemoveRefresh: %s\n", code)
	key := datastore.NewKey(context, "access.refresh", code, 0, nil)
	err := datastore.Delete(context, key)
	if err != nil {
		return err
	}

	return nil
}
