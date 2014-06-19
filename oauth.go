package glukit

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"fmt"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"github.com/alexandre-normand/glukit/lib/goauth2/oauth"
	"github.com/alexandre-normand/osin"
	"html/template"
	"net/http"
	"strings"
	"time"
)

var server *osin.Server

const (
	TOKEN_ROUTE     = "token"
	AUTHORIZE_ROUTE = "authorize"
)

type oauthAuthenticatedHandler struct {
	authenticatedHandler http.Handler
}

// Some variables that are used during rendering of oauth templates
type OauthRenderVariables struct {
	Code  string
	State string
}

var authorizeLocalAppTemplate = template.Must(template.ParseFiles("view/templates/oauthorize.html"))

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
		req.ParseForm()
		req.SetBasicAuth(req.Form.Get("client_id"), req.Form.Get("client_secret"))
		c.Debugf("Processing authorization request: %v and form [%v]", req, req.PostForm)
		if ar := server.HandleAuthorizeRequest(resp, req); ar != nil {
			// Nothing to do since the page is already login restricted by gae app configuration
			ar.Authorized = true
			ar.UserData = user.Email

			_, _, _, err := store.GetUserData(c, user.Email)
			if err == datastore.ErrNoSuchEntity {
				c.Debugf("Creating GlukitUser on first oauth access for [%s]: ", user.Email)
				// If the user doesn't exist already, create it
				glukitUser := model.GlukitUser{user.Email, "", "", time.Now(),
					model.DIABETES_TYPE_1, "", util.GLUKIT_EPOCH_TIME, apimodel.UNDEFINED_GLUCOSE_READ, oauth.Token{"", "", util.GLUKIT_EPOCH_TIME}, "",
					model.UNDEFINED_SCORE, model.UNDEFINED_SCORE, false, "", time.Now()}
				_, err = store.StoreUserProfile(c, time.Now(), glukitUser)
				if err != nil {
					resp.SetError(osin.E_SERVER_ERROR, fmt.Sprintf("Fail to initialize user for email [%s]: [%v]", user.Email, err))
					resp.StatusCode = 500
					osin.OutputJSON(resp, writer, request)
					return
				}
			} else if err != nil {
				resp.SetError(osin.E_SERVER_ERROR, fmt.Sprintf("Unable to find user for email [%s]: [%v]", user.Email, err))
				resp.StatusCode = 500
				osin.OutputJSON(resp, writer, request)
				return
			}

			if err != nil {
				util.Propagate(err)
			}

			server.FinishAuthorizeRequest(resp, req, ar)

			if resp.URL == "urn:ietf:wg:oauth:2.0:oob" {
				// Render a page with the title including the code
				data := resp.Output
				renderVariables := &OauthRenderVariables{Code: data["code"].(string), State: data["state"].(string)}

				if err := authorizeLocalAppTemplate.Execute(w, renderVariables); err != nil {
					c.Criticalf("Error executing template [%s]", authorizeLocalAppTemplate.Name())
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		}
		if resp.IsError && resp.InternalError != nil {
			c.Debugf("ERROR: %s\n", resp.InternalError)
		}
		c.Debugf("Writing response: %v", resp.Output)
		osin.OutputJSON(resp, w, req)
	})

	// Access token endpoint
	muxRouter.Get(TOKEN_ROUTE).HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c := appengine.NewContext(req)
		resp := server.NewResponse()
		req.ParseForm()
		req.SetBasicAuth(req.Form.Get("client_id"), req.Form.Get("client_secret"))
		c.Debugf("Processing token request: %v with form [%v]", req, req.PostForm)
		if ar := server.HandleAccessRequest(resp, req); ar != nil {
			c.Debugf("Retrieved authorize data [%v]", ar)
			ar.Authorized = true
			server.FinishAccessRequest(resp, req, ar)
		}
		c.Debugf("Writing response: %v", resp.Output)
		osin.OutputJSON(resp, w, req)
	})
	context := appengine.NewContext(request)
	context.Debugf("Oauth server loaded: [%v]", server)
}

func (handler *oauthAuthenticatedHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := appengine.NewContext(request)
	request.ParseForm()
	c.Debugf("Checking authentication for request [%s]...", request.RequestURI)

	ret := server.NewResponse()
	authorizationValue := request.Header.Get("Authorization")
	if authorizationValue == "" {
		ret.SetError(osin.E_INVALID_REQUEST, "Empty authorization")
		ret.StatusCode = 403
		osin.OutputJSON(ret, writer, request)
		return
	}

	accessCode := strings.TrimPrefix(authorizationValue, "Bearer ")
	if accessCode == "" {
		ret.SetError(osin.E_INVALID_REQUEST, "Empty authorization value")
		ret.StatusCode = 403
		osin.OutputJSON(ret, writer, request)
		return
	}

	var err error

	// load access data
	accessData, err := server.Storage.LoadAccess(accessCode, request)
	if err != nil {
		ret.SetError(osin.E_INVALID_REQUEST, fmt.Sprintf("Error loading access data for code [%s]: [%v]", accessCode, err))
		ret.StatusCode = 403
		osin.OutputJSON(ret, writer, request)
		return
	}
	if accessData.Client == nil {
		ret.SetError(osin.E_UNAUTHORIZED_CLIENT, "")
		ret.StatusCode = 403
		osin.OutputJSON(ret, writer, request)
		return
	}
	if accessData.Client.RedirectUri == "" {
		ret.SetError(osin.E_UNAUTHORIZED_CLIENT, "")
		ret.StatusCode = 403
		osin.OutputJSON(ret, writer, request)
		return
	}
	if accessData.IsExpired() {
		ret.SetError(osin.E_INVALID_GRANT, "")
		ret.StatusCode = 403
		osin.OutputJSON(ret, writer, request)
		return
	}

	if ret.IsError {
		ret.StatusCode = 403
		osin.OutputJSON(ret, writer, request)
		return
	}

	handler.authenticatedHandler.ServeHTTP(writer, request)
}

func newOauthAuthenticationHandler(next http.Handler) *oauthAuthenticatedHandler {
	return &oauthAuthenticatedHandler{next}
}
