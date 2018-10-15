package streamer

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/rk23/tdgo/pkg/user"
)

type lc struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type response struct {
	Service   string `json:"service"`
	Command   string `json:"command"`
	RequestID int    `json:"requestid"`
	Timestamp int64  `json:"timestmap"`
	Content   lc     `json:"content"`
}

type responses struct {
	Responses []response `json:"responses"`
}

type reqparams struct {
	Credential string `json:"credential"`
	Token      string `json:"token"`
	Version    string `json:"version"`
}

type request struct {
	Service    string    `json:"service"`
	Command    string    `json:"command"`
	RequestID  int       `json:"requestid"`
	Account    string    `json:"account"`
	Source     string    `json:"source"`
	Parameters reqparams `json:"parameters"`
}

type requests struct {
	Requests []request `json:"requests"`
}

func encode(m map[string]string) string {
	q := ""
	for k, v := range m {
		q += url.QueryEscape(k) + url.QueryEscape("=") + url.QueryEscape(v) + url.QueryEscape("&")
	}
	return q
}

// LoginRequest returns a login request
func LoginRequest(user *user.Principals) ([]byte, error) {

	// ISO-8601 != RFC3339
	t, err := time.Parse("2006-01-02T15:04:05+0000", user.StreamerInfo.TokenTimestamp)
	if err != nil {
		fmt.Println("err: " + err.Error())
		return nil, err
	}
	ms := strconv.FormatInt(t.Unix()*1000, 10)

	credentials := map[string]string{
		"userid":      user.Accounts[0].AccountID,
		"token":       user.StreamerInfo.Token,
		"company":     user.Accounts[0].Company,
		"segment":     user.Accounts[0].Segment,
		"cddomain":    user.Accounts[0].AccountCdDomainID,
		"usergroup":   user.StreamerInfo.UserGroup,
		"accesslevel": user.StreamerInfo.AccessLevel,
		"authorized":  "Y",
		"timestamp":   ms,
		"appid":       user.StreamerInfo.AppID,
		"acl":         user.StreamerInfo.ACL,
	}

	r := &requests{[]request{
		request{
			Service:   "ADMIN",
			Command:   "LOGIN",
			RequestID: 0,
			Account:   user.Accounts[0].AccountID,
			Source:    user.StreamerInfo.AppID,
			Parameters: reqparams{
				Credential: encode(credentials),
				Token:      user.StreamerInfo.Token,
				Version:    "1.0",
			},
		}}}

	req, err := json.Marshal(r)

	return req, err
}

// LogoutRequest does a thing
func LogoutRequest(user *user.Principals) ([]byte, error) {
	r := &requests{[]request{
		request{
			Service:    "ADMIN",
			Command:    "LOGOUT",
			RequestID:  1,
			Account:    user.Accounts[0].AccountID,
			Source:     user.StreamerInfo.AppID,
			Parameters: reqparams{},
		}}}

	req, err := json.Marshal(r)

	return req, err
}
