package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// StreamerCredentials to pass to streamer request
type StreamerCredentials struct {
	AccountID         string `json:"userid"`
	Token             string `json:"token"`
	Company           string `json:"company"`
	Segment           string `json:"segment"`
	AccountCdDomainID string `json:"cddomain"`
	UserGroup         string `json:"usergroup"`
	AccessLevel       string `json:"accesslevel"`
	Authorized        string `json:"authorized"`
	Timestamp         int64  `json:"timestamp"`
	AppID             string `json:"appid"`
	ACL               string `json:"acl"`
}

// CredsToQueryString converts struct to query string
func (sc *StreamerCredentials) CredsToQueryString() (string, error) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		panic(err)
	}
	q := req.URL.Query()
	q.Add("test", "value")

	fmt.Printf("")
	return "not implemented", nil
}

// AccessTokenRequest specifies necessary info for request
type AccessTokenRequest struct {
	AccessType  string
	Code        string
	ClientID    string
	RedirectURI string
}

// AccessTokenResponse specifies AT response format
type AccessTokenResponse struct {
	AccessToken        string `json:"access_token"`
	RefreshToken       string `json:"refresh_token"`
	TokenType          string `json:"token_type"`
	ExpiresIn          int64  `json:"expires_in"`
	Scope              string `json:"scope"`
	RefreshTokenExpiry int64  `json:"refresh_token_expires_in"`
}

// AccessToken function POSTs to the TD Ameritrade server to get the new tokens
func AccessToken(req *AccessTokenRequest) (*AccessTokenResponse, error) {
	atr := &AccessTokenResponse{}

	// Must be ordered as specified in the docs! Go's url.Values will
	// sort them into alphabetical order.
	ordered := fmt.Sprintf("grant_type=authorization_code&refresh_token=&access_type=%s&code=%s&client_id=%s&redirect_uri=%s", url.QueryEscape(req.AccessType), url.QueryEscape(req.Code), url.QueryEscape(req.ClientID), url.QueryEscape(req.RedirectURI))

	client := &http.Client{}
	r, _ := http.NewRequest("POST", "https://api.tdameritrade.com/v1/oauth2/token", strings.NewReader(ordered))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	fmt.Printf("%+v", r)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		errBody, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("%s", string(errBody))
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(bodyBytes, atr)
	if err != nil {
		return nil, err
	}

	return atr, nil
}

// RefreshToken function refreshes access token timeout
// TODO: No longer apart of the struct - need to pass in res values as well
func RefreshToken(req *AccessTokenRequest) (*AccessTokenResponse, error) {
	res := &AccessTokenResponse{}

	// Must be ordered as specified in the docs! Go's url.Values will
	// sort them into alphabetical order.
	ordered := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&access_type=%s&code=%s&client_id=%s&redirect_uri=%s",
		url.QueryEscape(req.AccessType),
		url.QueryEscape(res.RefreshToken),
		url.QueryEscape(res.AccessToken),
		url.QueryEscape(req.ClientID),
		url.QueryEscape(req.RedirectURI))

	client := &http.Client{}
	nr, _ := http.NewRequest("POST", "https://api.tdameritrade.com/v1/oauth2/token", strings.NewReader(ordered))
	nr.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	r, err := client.Do(nr)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)

	fmt.Printf("%+v", string(bodyBytes))

	err = json.Unmarshal(bodyBytes, res)
	if err != nil {
		panic(err)
	}

	return res, nil
}
