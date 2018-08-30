package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/danesparza/authserver/data"
)

// AuthRequest is an OAuth2 based request.  For more information on the
// various grant types that can use this request object:
// https://alexbilbie.com/guide-to-oauth-2-grants/
type AuthRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
	UserName     string `json:"username"`
	Password     string `json:"password"`
	CSRFToken    string `json:"state"`
	RedirectURI  string `json:"redirect_uri"`
	ResponseType string `json:"response_type"`
	Code         string `json:"code"`
}

// AuthResponse is an OAuth2 based response
type AuthResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// HelloWorld emits a hello world
func HelloWorld(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello, world - service")
}

// ClientCredentialsGrant implements the OAuth 2 'Client Credentials' grant --
// see https://alexbilbie.com/guide-to-oauth-2-grants/ for more information
func (service Service) ClientCredentialsGrant(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Get the authorization header:
	authHeader := req.Header.Get("Authorization")

	//	If the basic auth header wasn't supplied, return an error
	if basicHeaderValid(authHeader) != true {
		sendErrorResponse(rw, fmt.Errorf("HTTP basic auth credentials not supplied"), http.StatusUnauthorized)
		return
	}

	//	Get just the credentials from basic auth information:
	clientid, clientsecret := getCredentialsFromAuthHeader(authHeader)

	//	Decode the request using ParseForm:
	err := req.ParseForm()
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	/*
		log.Println("Parsed grant type: ", req.PostForm["grant_type"])
		log.Println("Parsed scopes: ", req.PostForm["scope"])
	*/

	//	Send the request to the datamanager and get grant information for the given credentials:
	scopeUser, err := service.DB.GetUserScopesWithCredentials(clientid, clientsecret)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusUnauthorized)
		return
	}

	//	Get a token for the returned user information
	token, err := service.DB.GetNewToken(data.User{ID: scopeUser.ID}, 1*time.Hour)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusUnauthorized)
		return
	}

	//	Create our response and send information back:
	encodedToken := base64.StdEncoding.EncodeToString([]byte(token.ID))
	response := AuthResponse{
		TokenType:   "Bearer",
		ExpiresIn:   strconv.FormatFloat(token.Expires.Sub(time.Now()).Seconds(), 'f', 0, 64),
		AccessToken: encodedToken,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// ScopesForUserID gets the scope information for the userID passed in the url
// @Summary gets the scope information
// @Description gets the scope information for the userID passed in the url
// @ID scopes-for-user-id
// @Accept  json
// @Produce  json
// @Security OAuth2Application
// @Success 200 {object} api.AuthResponse
// @Failure 401 {object} api.ErrorResponse
// @Router /oauth/authorize [get]
func (service Service) ScopesForUserID(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Get the authorization header:
	authHeader := req.Header.Get("Authorization")

	//	If the auth header wasn't supplied, return an error
	if authHeaderValid(authHeader) != true {
		sendErrorResponse(rw, fmt.Errorf("Bearer token was not supplied"), http.StatusUnauthorized)
		return
	}

	//	Get just the bearer token itself:
	token := getTokenFromAuthHeader(authHeader)

	//	Send the request to the datamanager and get scope information for the given credentials:
	response, err := service.DB.GetScopesForToken(token)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusUnauthorized)
		return
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// authHeaderValid returns true if the passed header value is a valid
// for a "bearer token" authorization field -- otherwise return false
func authHeaderValid(header string) bool {
	retval := true

	//	If we don't have at least x number characters,
	//	it must not include the prefix text 'Bearer '
	if len(header) < len("Bearer ") {
		return false
	}

	//	If the first part of the string isn't 'Bearer ' then it's not a bearer token...
	if strings.EqualFold(header[:len("Bearer ")], "Bearer ") != true {
		return false
	}

	return retval
}

// basicHeaderValid returns true if the passed header value is a valid
// for a "http basic authentication" authorization field -- otherwise return false
func basicHeaderValid(header string) bool {
	retval := true

	//	If we don't have at least x number characters,
	//	it must not include the prefix text 'Basic '
	if len(header) < len("Basic ") {
		return false
	}

	//	If the first part of the string isn't 'Basic ' then it's not a basic auth header...
	if strings.EqualFold(header[:len("Basic ")], "Basic ") != true {
		return false
	}

	return retval
}

// getTokenFromAuthHeader returns the token itself from the Authorization header
func getTokenFromAuthHeader(header string) string {
	retval := ""

	//	If we don't have at least x number characters,
	//	it must not include the prefix text 'Bearer '
	if len(header) < len("Bearer ") {
		return ""
	}

	//	If the first part of the string isn't 'Bearer ' then it's not a bearer token...
	if strings.EqualFold(header[:len("Bearer ")], "Bearer ") != true {
		return ""
	}

	//	Get the token and decode it
	encodedToken := header[len("Bearer "):]
	tokenBytes, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		return ""
	}

	//	Change the type to string
	retval = string(tokenBytes)

	return retval
}

// getCredentialsFromAuthHeader returns the username/password from the Authorization header
func getCredentialsFromAuthHeader(header string) (string, string) {
	username := ""
	password := ""

	//	If we don't have at least x number characters,
	//	it must not include the prefix text 'Basic '
	if len(header) < len("Basic ") {
		return "", ""
	}

	//	If the first part of the string isn't 'Basic ' then it's not a basic auth string...
	if strings.EqualFold(header[:len("Basic ")], "Basic ") != true {
		return "", ""
	}

	//	Get the credentials and decode them
	encodedCredentials := header[len("Basic "):]
	credentialBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return "", ""
	}

	//	Change the type to string
	credentials := strings.Split(string(credentialBytes), ":")

	username = credentials[0]
	password = credentials[1]

	return username, password
}
