package fields

// All possible fields for Option
const (
	OptionSymbol                 = 0
	OptionDescription            = 1
	OptionBidPrice               = 2
	OptionAskPrice               = 3
	OptionLastPrice              = 4
	OptionHighPrice              = 5
	OptionLowPrice               = 6
	OptionClosePrice             = 7
	OptionTotalVolume            = 8
	OptionOpenInterest           = 9
	OptionVolatility             = 10
	OptionQuoteTime              = 11
	OptionTradeTime              = 12
	OptionMoneyIntrinsicValue    = 13
	OptionQuoteDay               = 14
	OptionTradeDay               = 15
	OptionExpirationYear         = 16
	OptionMultiplier             = 17
	OptionDigits                 = 18
	OptionOpenPrice              = 19
	OptionBidSize                = 20
	OptionAskSize                = 21
	OptionLastSize               = 22
	OptionNetChange              = 23
	OptionStrikePrice            = 24
	OptionContractType           = 25
	OptionUnderlying             = 26
	OptionExpirationMonth        = 27
	OptionDeliverables           = 28
	OptionTimeValue              = 29
	OptionExpirationDay          = 30
	OptionDaystoExpiration       = 31
	OptionDelta                  = 32
	OptionGamma                  = 33
	OptionTheta                  = 34
	OptionVega                   = 35
	OptionRho                    = 36
	OptionSecurityStatus         = 37
	OptionTheoreticalOptionValue = 38
	OptionUnderlyingPrice        = 39
	OptionUVExpirationType       = 40
	OptionMark                   = 41
)

// OptionFields struct for unmarshalling option response
type OptionFields struct {
	Symbol                 string
	Description            string
	BidPrice               float64
	AskPrice               float64
	LastPrice              float64
	HighPrice              float64
	LowPrice               float64
	ClosePrice             float64
	TotalVolume            int64
	OpenInterest           int
	Volatility             float64
	QuoteTime              int64
	TradeTime              int64
	MoneyIntrinsicValue    float64
	QuoteDay               int
	TradeDay               int
	ExpirationYear         int
	Multiplier             float64
	Digits                 int
	OpenPrice              float64
	BidSize                float64
	AskSize                float64
	LastSize               float64
	NetChange              float64
	StrikePrice            float64
	ContractType           string
	Underlying             string
	ExpirationMonth        int
	Deliverables           string
	TimeValue              float64
	ExpirationDay          int
	DaystoExpiration       int
	Delta                  float64
	Gamma                  float64
	Theta                  float64
	Vega                   float64
	Rho                    float64
	SecurityStatus         string
	TheoreticalOptionValue float64
	UnderlyingPrice        float64
	UVExpirationType       string
	Mark                   float64
}
