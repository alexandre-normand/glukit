package glukit

import (
	"appengine"
	"appengine/user"
	"fmt"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/osin"
	"net/http"
	"strings"
)

var server *osin.Server

const (
	TOKEN_ROUTE     = "token"
	AUTHORIZE_ROUTE = "authorize"
)

type oauthAuthenticatedHandler struct {
	authenticatedHandler http.Handler
}

func initOauthProvider(writer http.ResponseWriter, request *http.Request) {
	sconfig := osin.NewServerConfig()
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN}
	// 30 days
	sconfig.AccessExpiration = 60 * 60 * 24 * 30
	sconfig.AllowGetAccessRequest = true
	server = osin.NewServer(sconfig, store.NewOsinAppEngineStoreWithRequest(request))
	muxRouter.Get(AUTHORIZE_ROUTE).HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c := appengine.NewContext(req)
		user := user.Current(c)
		resp := server.NewResponse()
		c.Debugf("Processing authorization request: %v", req)
		req.ParseForm()
		req.SetBasicAuth(req.Form.Get("client_id"), req.Form.Get("client_secret"))
		if ar := server.HandleAuthorizeRequest(resp, req); ar != nil {
			// Nothing to do since the page is already login restricted by gae app configuration
			ar.Authorized = true
			ar.UserData = user.Email
			server.FinishAuthorizeRequest(resp, req, ar)

			if resp.URL == "urn:ietf:wg:oauth:2.0:oob" {
				resp.Type = osin.DATA
			}
		}
		if resp.IsError && resp.InternalError != nil {
			c.Debugf("ERROR: %s\n", resp.InternalError)
		}
		osin.OutputJSON(resp, w, req)
	})

	// Access token endpoint
	muxRouter.Get(TOKEN_ROUTE).HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c := appengine.NewContext(req)
		c.Debugf("Processing token request: %v", req)
		user := user.Current(c)
		resp := server.NewResponse()
		req.ParseForm()
		req.SetBasicAuth(req.Form.Get("client_id"), req.Form.Get("client_secret"))
		if ar := server.HandleAccessRequest(resp, req); ar != nil {
			c.Debugf("Processing token request with form: %v", req.Form)
			ar.Authorized = true
			ar.UserData = user.Email
			server.FinishAccessRequest(resp, req, ar)
		}
		osin.OutputJSON(resp, w, req)
	})
	context := appengine.NewContext(request)
	context.Debugf("Oauth server loaded: [%v]", server)
}

func (handler *oauthAuthenticatedHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := appengine.NewContext(request)
	request.ParseForm()
	c.Debugf("Checking authentication for request [%v]...", request)

	ret := server.NewResponse()
	authorizationValue := request.Header.Get("Authorization")
	if authorizationValue == "" {
		ret.SetError(osin.E_INVALID_REQUEST, "Empty authorization")
		osin.OutputJSON(ret, writer, request)
		return
	}

	accessCode := strings.TrimPrefix(authorizationValue, "Bearer ")
	if accessCode == "" {
		ret.SetError(osin.E_INVALID_REQUEST, "Empty authorization value")
		osin.OutputJSON(ret, writer, request)
		return
	}

	var err error

	// load access data
	accessData, err := server.Storage.LoadAccess(accessCode, request)
	if err != nil {
		ret.SetError(osin.E_INVALID_REQUEST, fmt.Sprintf("Error loading access data for code [%s]: [%v]", accessCode, err))
		ret.InternalError = err
		osin.OutputJSON(ret, writer, request)
		return
	}
	if accessData.Client == nil {
		ret.SetError(osin.E_UNAUTHORIZED_CLIENT, "")
		osin.OutputJSON(ret, writer, request)
		return
	}
	if accessData.Client.RedirectUri == "" {
		ret.SetError(osin.E_UNAUTHORIZED_CLIENT, "")
		osin.OutputJSON(ret, writer, request)
		return
	}
	if accessData.IsExpired() {
		ret.SetError(osin.E_INVALID_GRANT, "")
		osin.OutputJSON(ret, writer, request)
		return
	}

	if ret.IsError {
		osin.OutputJSON(ret, writer, request)
		return
	}

	handler.authenticatedHandler.ServeHTTP(writer, request)
}

func newOauthAuthenticationHandler(next http.Handler) *oauthAuthenticatedHandler {
	return &oauthAuthenticatedHandler{next}
}
