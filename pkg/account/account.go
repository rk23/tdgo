package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
AssetType: 'EQUITY' or 'OPTION' or 'INDEX' or 'MUTUAL_FUND' or 'CASH_EQUIVALENT' or 'FIXED_INCOME' or 'CURRENCY'
ComplexOrderStrategyType: 'NONE' or 'COVERED' or 'VERTICAL' or 'BACK_RATIO' or 'CALENDAR' or 'DIAGONAL' or 'STRADDLE' or 'STRANGLE' or 'COLLAR_SYNTHETIC' or 'BUTTERFLY' or 'CONDOR' or 'IRON_CONDOR' or 'VERTICAL_ROLL' or 'COLLAR_WITH_STOCK' or 'DOUBLE_DIAGONAL' or 'UNBALANCED_BUTTERFLY' or 'UNBALANCED_CONDOR' or 'UNBALANCED_IRON_CONDOR' or 'UNBALANCED_VERTICAL_ROLL' or 'CUSTOM'
CurrencyType: 'USD' or 'CAD' or 'EUR' or 'JPY'
Execution Activity Type: 'EXECUTION' or 'ORDER_ACTION'
Execution Type: 'FILL'
Option Type: 'VANILLA' or 'BINARY' or 'BARRIER'
Order Duration: 'DAY' or 'GOOD_TILL_CANCEL' or 'FILL_OR_KILL'
OrderLegType: 'EQUITY' or 'OPTION' or 'INDEX' or 'MUTUAL_FUND' or 'CASH_EQUIVALENT' or 'FIXED_INCOME' or 'CURRENCY'
Order Leg Instruction: 'BUY' or 'SELL' or 'BUY_TO_COVER' or 'SELL_SHORT' or 'BUY_TO_OPEN' or 'BUY_TO_CLOSE' or 'SELL_TO_OPEN' or 'SELL_TO_CLOSE' or 'EXCHANGE'
Order Leg PositionEffect: 'OPENING' or 'CLOSING' or 'AUTOMATIC'
Order Leg QuantityType: 'ALL_SHARES' or 'DOLLARS' or 'SHARES'
Order Type: 'MARKET' or 'LIMIT' or 'STOP' or 'STOP_LIMIT' or 'TRAILING_STOP' or 'MARKET_ON_CLOSE' or 'EXERCISE' or 'TRAILING_STOP_LIMIT' or 'NET_DEBIT' or 'NET_CREDIT' or 'NET_ZERO'
Order Session: 'NORMAL' or 'AM' or 'PM' or 'SEAMLESS'
Order Status: 'AWAITING_PARENT_ORDER' or 'AWAITING_CONDITION' or 'AWAITING_MANUAL_REVIEW' or 'ACCEPTED' or 'AWAITING_UR_OUT' or 'PENDING_ACTIVATION' or 'QUEUED' or 'WORKING' or 'REJECTED' or 'PENDING_CANCEL' or 'CANCELED' or 'PENDING_REPLACE' or 'REPLACED' or 'FILLED' or 'EXPIRED'
OrderStrategyType: 'SINGLE' or 'OCO' or 'TRIGGER'
PriceLinkBasis: 'MANUAL' or 'BASE' or 'TRIGGER' or 'LAST' or 'BID' or 'ASK' or 'ASK_BID' or 'MARK' or 'AVERAGE'
PriceLinkType: 'VALUE' or 'PERCENT' or 'TICK'
PutCall: 'PUT' or 'CALL'
RequestedDestination: 'INET' or 'ECN_ARCA' or 'CBOE' or 'AMEX' or 'PHLX' or 'ISE' or 'BOX' or 'NYSE' or 'NASDAQ' or 'BATS' or 'C2' or 'AUTO'
SpecialInstructions: 'ALL_OR_NONE' or 'DO_NOT_REDUCE' or 'ALL_OR_NONE_DO_NOT_REDUCE'
StopPriceLinkBasis: 'MANUAL' or 'BASE' or 'TRIGGER' or 'LAST' or 'BID' or 'ASK' or 'ASK_BID' or 'MARK' or 'AVERAGE'
StopPriceLinkType: 'VALUE' or 'PERCENT' or 'TICK'
StopType: 'STANDARD' or 'BID' or 'ASK' or 'LAST' or 'MARK'
TaxLotMethod: 'FIFO' or 'LIFO' or 'HIGH_COST' or 'LOW_COST' or 'AVERAGE_COST' or 'SPECIFIC_LOT'
Mutual Fund Type: 'NOT_APPLICABLE' or 'OPEN_END_NON_TAXABLE' or 'OPEN_END_TAXABLE' or 'NO_LOAD_NON_TAXABLE' or 'NO_LOAD_TAXABLE'
*/

type ExecutionLeg struct {
	LegID             int     `json:"legId"`
	Quantity          float32 `json:"quantity"`
	MismarkedQuantity float32 `json:"mismarkedQuantity"`
	Price             float32 `json:"price"`
	Time              string  `json:"time"`
}

type Execution struct {
	ActivityType           string         `json:"activityType"`
	ExecutionType          string         `json:"executionType"`
	Quantity               float32        `json:"quantity"`
	OrderRemainingQuantity float32        `json:"orderRemainingQuantity"`
	ExecutionLegs          []ExecutionLeg `json:"executionLegs"`
}

type OptionDeliverable struct {
	Symbol           string `json:"symbol"`
	DeliverableUnits int    `json:"deliverableUnits"`
	CurrencyType     string `json:"currencyType"`
	AssetType        string `json:"assetType"`
}

type Instrument struct {
	AssetType   string `json:"assetType"`
	CUSIP       string `json:"cusip,omitempty"`
	Symbol      string `json:"symbol"`
	Description string `json:"description,omitempty"`
}

type Option struct {
	Instrument
	Type               string              `json:"type"`
	PutCall            string              `json:"putCall"`
	UnderlyingSymbol   string              `json:"underlyingSymbol"`
	OptionMultiplier   int                 `json:"optionMultiplier`
	OptionDeliverables []OptionDeliverable `json:"optionDeliverables"`
}

type MutualFund struct {
	Instrument
	Type string `json:"type"`
}

type CashEquivalent struct {
	Instrument
	Type string `json:"type"`
}

type Equity struct {
	Instrument
}

type FixedIncome struct {
	Instrument
	MaturityDate string  `json:"maturityDate"`
	VariableRate float32 `json:"variableRate"`
	Factor       int     `json:"factor"`
}

type OrderCancelTime struct {
	Date        string `json:"date,omitempty"`
	ShortFormat bool   `json:"shortFormat,omitempty"`
}

type OrderLeg struct {
	OrderLegType   string     `json:"orderLegType,omitempty"`
	LegID          int        `json:"legId,omitempty"`
	Instrument     Instrument `json:"instrument"`
	Instruction    string     `json:"instruction"`
	PositionEffect string     `json:"positionEffect,omitempty"`
	Quantity       float32    `json:"quantity"`
	QuantityType   string     `json:"quantityType,omitempty"`
}

type Order struct {
	Session                  string           `json:"session"`
	Duration                 string           `json:"duration"`
	OrderType                string           `json:"orderType"`
	CancelTime               *OrderCancelTime `json:"cancelTime,omitempty"`
	ComplexOrderStrategyType string           `json:"complexOrderStrategyType"`
	Quantity                 float32          `json:"quantity,omitempty"`
	FilledQuantity           float32          `json:"filledQuantity,omitempty"`
	RemainingQuantity        float32          `json:"remainingQuantity,omitempty"`
	RequestedDestination     string           `json:"requestedDestination,omitempty"`
	DestinationLinkName      string           `json:"destinationLinkName,omitempty"`
	ReleaseTime              string           `json:"releaseTime,omitempty"`
	StopPrice                float32          `json:"stopPrice,omitempty"`
	StopPriceLinkBasis       string           `json:"stopPriceLinkBasis,omitempty"`
	StopPriceLinkType        string           `json:"stopPriceLinkType,omitempty"`
	StopPriceOffset          float32          `json:"stopPriceOffset,omitempty"`
	StopType                 string           `json:"stopType,omitempty"`
	PriceLinkBasis           string           `json:"priceLinkBasis,omitempty"`
	PriceLinkType            string           `json:"priceLinkType,omitempty"`
	Price                    float32          `json:"price,omitempty"`
	TaxLotMethod             string           `json:"taxLotMethod,omitempty"`
	OrderLegCollection       []OrderLeg       `json:"orderLegCollection"`
	ActivationPrice          float32          `json:"activationPrice,omitempty"`
	SpecialInstructions      string           `json:"specialInstruction,omitempty"`
	OrderStrategyType        string           `json:"orderStrategyType"`
	OrderID                  int              `json:"orderId,omitempty"`
	Cancelable               bool             `json:"cancelable,omitempty"`
	Editable                 bool             `json:"editable,omitempty"`
	Status                   string           `json:"status,omitempty"`
	EnteredTime              string           `json:"enteredTime,omitempty"`
	CloseTime                string           `json:"closeTime,omitempty"`
	Tag                      string           `json:"tag,omitempty"`
	AccountID                int              `json:"accountId,omitempty"`
	OrderActivityCollection  []Execution      `json:"orderActivityCollection,omitempty"`
	ReplacingOrderCollection []Execution      `json:"replacingOrderCollection,omitempty"`
	ChildOrderStrategies     []Execution      `json:"childOrderStrategies,omitempty"`
	StatusDescription        string           `json:"statusDescription,omitempty"`
}

type GetOrderRequest struct {
	AccountID   string
	BearerToken string
	OrderID     string
}

type GetOrdersRequest struct {
	GetOrderRequest
	MaxResults      *int
	FromEnteredTime *string
	ToEnteredTime   *string
	Status          *string
}

func GetOrder(r GetOrderRequest) (*Order, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://api.tdameritrade.com/v1/accounts/"+r.AccountID+"/orders/"+r.OrderID, nil)
	req.Header.Add("Authorization", "Bearer "+r.BearerToken)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 300 {
		return nil, fmt.Errorf("Failed to get orders for account: %s, status code: %d", r.AccountID, res.StatusCode)
	}

	o := &Order{}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(bodyBytes, o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// GetOrders returns orders from the specified account on the same given day.
// Date format is yyyy-MM-dd
func GetOrders(r GetOrdersRequest) (*[]Order, error) {
	client := &http.Client{}
	var url string

	// Must be ordered as specified in the docs! Go's url.Values will
	// sort them into alphabetical order.
	ordered := "?"
	if r.MaxResults != nil {
		ordered += fmt.Sprintf("maxResults=%d", *r.MaxResults)
	}
	if r.FromEnteredTime != nil {
		ordered += "fromEnteredTime=" + *r.FromEnteredTime + "&"
	}
	if r.ToEnteredTime != nil {
		ordered += "toEnteredTime=" + *r.ToEnteredTime
	}
	if r.Status != nil {
		ordered += "status=" + *r.Status
	}

	if r.AccountID != "" {
		url = "https://api.tdameritrade.com/v1/accounts/" + r.AccountID + "/orders" + ordered
	} else {
		url = "https://api.tdameritrade.com/v1/accounts/orders" + ordered
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+r.BearerToken)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 300 {
		return nil, fmt.Errorf("Failed to get orders for account: %s, status code: %d", r.AccountID, res.StatusCode)
	}

	o := &[]Order{}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(bodyBytes))
	err = json.Unmarshal(bodyBytes, o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

type PlaceOrderRequest struct {
	Order
	BearerToken string
	AccountID   string
}

func PlaceOrder(r PlaceOrderRequest) error {
	client := &http.Client{}

	body, err := json.Marshal(r.Order)
	if err != nil {
		panic(err)
	}
	req, _ := http.NewRequest("POST", "https://api.tdameritrade.com/v1/accounts/"+r.AccountID+"/orders", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+r.BearerToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		return fmt.Errorf("Failed to place order(s) for account: " + r.AccountID + ", reason: " + string(bodyBytes) + ", status: " + res.Status)
	}

	return nil
}
