package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

// ClientCredentialsGrant implements the OAuth 2 'Client Credentials' grant
func (service Service) ClientCredentialsGrant(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request if it was a POST:
	request := AuthRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Send the request to the datamanager and get grant information for the given credentials:
	grantInfo, err := service.DB.GetUserGrantsWithCredentials(request.ClientID, request.ClientSecret)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusUnauthorized)
		return
	}

	//	Get a token for the returned user information
	token, err := service.DB.GetNewToken(data.User{ID: grantInfo.ID}, 1*time.Hour)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusUnauthorized)
		return
	}

	//	Create our response and send information back:
	response := AuthResponse{
		TokenType:   "Bearer",
		ExpiresIn:   strconv.FormatFloat(token.Expires.Sub(time.Now()).Seconds(), 'f', 0, 64),
		AccessToken: token.ID,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
