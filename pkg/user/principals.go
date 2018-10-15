package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// StreamerInfo is a subfield of UserPrincipal
type StreamerInfo struct {
	StreamerBinaryURL string `json:"streamerBinaryUrl"`
	StreamerSocketURL string `json:"streamerSocketUrl"`
	Token             string `json:"token"`
	TokenTimestamp    string `json:"tokenTimestamp"`
	UserGroup         string `json:"userGroup"`
	AccessLevel       string `json:"accessLevel"`
	ACL               string `json:"acl"`
	AppID             string `json:"appId"`
}

// QuotesInfo is a subfield of UserPrincipal
type QuotesInfo struct {
	IsNyseDelayed   bool `json:"isNyseDelayed"`
	IsNasdaqDelayed bool `json:"isNasdaqDelayed"`
	IsOpraDelayed   bool `json:"isOpraDelayed"`
	IsAmexDelayed   bool `json:"isAmexDelayed"`
	IsCmeDelayed    bool `json:"isCmeDelayed"`
	IsIceDelayed    bool `json:"isIceDelayed"`
	isForexDelayed  bool `josn:"isForexDelayed"`
}

//SubKeys is the StreamerSubscriptionKeys subfield of UserPrincipal
type SubKeys struct {
	Keys []map[string]string `json:"keys"`
}

// AccountPreferences is a subfield of AccountInfo
type AccountPreferences struct {
	ExpressTrading                    bool   `json:"expressTrading"`
	DirectOptionsRouting              bool   `json:"directOptionsRouting"`
	DirectEquityRouting               bool   `json:"directEquityRouting"`
	DefaultEquityOrderLegInstructions string `json:"defaultEquityOrderLegInstructions" `
	DefaultEquityOrderType            string `json:"defaultEquityOrderType"`
	DefaultEquityOrderPriceLinkType   string `json:"defaultEquityOrderPriceLinkType"`
	DefaultEquityOrderDuration        string `json:"defaultEquityOrderDuration"`
	DefaultEquityOrderMarketSession   string `json:"defaultEquityOrderMarketSession"`
	DefaultEquityQuantity             int    `json:"defaultEquityQuantity"`
	MutalFundTaxLotMethod             string `json:"mutalFundTaxLotMethod"`
	OptionTaxLotMethod                string `json:"optionTaxLotMethod"`
	EquityTaxLotMethod                string `json:"equityTaxLotMethod"`
	DefaultAdvancedToolLaunch         string `json:"defaultAdvancedToolLaunch"`
	AuthTokenTimeout                  string `json:"authTokenTimeout"`
}

// AccountAuthorizations is a subfield of AccountInfo
type AccountAuthorizations struct {
	Apex               bool   `json:"apex"`
	LevelTwoQuotes     bool   `json:"levelTwoQuotes"`
	StockTrading       bool   `json:"stockTrading"`
	MarginTrading      bool   `json:"marginTrading"`
	StreamingNews      bool   `json:"streamingNews"`
	OptionTradingLevel string `json:"optionTradingLevel"`
	StreamerAccess     bool   `json:"streamerAccess"`
	AdvancedMargin     bool   `json:"advancedMargin"`
	ScottradeAccount   bool   `json:"scottradeAccount"`
}

// AccountInfo is the info for a particular account
type AccountInfo struct {
	AccountID         string                `json:"accountId"`
	Description       string                `json:"description"`
	DisplayName       string                `json:"displayName"`
	AccountCdDomainID string                `json:"accountCdDomainId"`
	Company           string                `json:"company"`
	Segment           string                `json:"segment"`
	SurrogateIDs      map[string]string     `json:"surrogateIds"`
	Preferences       AccountPreferences    `json:"preferences"`
	Authorizations    AccountAuthorizations `json:"authorizations"`
}

// Principals is the auth response
type Principals struct {
	AuthToken           string        `json:"authToken"`
	UserID              string        `json:"userId"`
	UserCdDomainID      string        `json:"userCdDomainId"`
	PrimaryAccountID    string        `json:"primaryAccountId"`
	LastLoginTime       string        `json:"lastLoginTime"`
	TokenExpirationTime string        `json:"tokenExpirationTime"`
	LoginTime           string        `json:"loginTime"`
	AccessLevel         string        `json:"accessLevel"`
	StalePassword       bool          `json:"stalePassword"`
	StreamerInfo        StreamerInfo  `json:"streamerInfo"`
	ProfessionalStatus  string        `json:"professionalStatus"`
	Quotes              QuotesInfo    `json:"quotes"`
	SubKeys             SubKeys       `json:"streamerSubscriptionKeys"`
	Accounts            []AccountInfo `json:"accounts"`
}

// GetPrincipals gets the user's principal information
func GetPrincipals(bearerToken string) (*Principals, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://api.tdameritrade.com/v1/userprincipals?fields=streamerSubscriptionKeys%2CstreamerConnectionInfo%2Cpreferences%2CsurrogateIds", nil)
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	res, err := client.Do(req)
	if err != nil {
		return &Principals{}, fmt.Errorf("Failed to GET principals from TD Server: %s", err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode > 300 {
		return &Principals{}, fmt.Errorf("Failed to get principals, status code: %d", res.StatusCode)
	}

	p := &Principals{}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(bodyBytes, p)
	if err != nil {
		return &Principals{}, fmt.Errorf("Failed to unmarshal principals: %s", err.Error())
	}

	return p, nil
}
