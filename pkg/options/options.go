package options

type Deliverable struct{}

type StrikeQuote struct {
	Ask                    float32
	AskSize                int32
	Bid                    float32
	BidSize                int32
	ClosePrice             float32
	DaysToExpiration       int
	Delta                  float32
	ExchangeName           string
	ExpirationDate         int64
	ExpirationType         string
	Gamma                  float32
	HighPrice              float32
	InTheMoney             bool
	IsIndexOption          bool
	Last                   float32
	LastSize               int32
	LastTradingDay         int64
	LowPrice               float32
	Mark                   float32
	MarkChange             float32
	MarkPercentChange      float32
	Mini                   bool
	Multiplier             float32
	NetChange              float32
	NonStandard            bool
	OpenInterest           float32
	OpenPrice              float32
	PercentChange          float32
	PutCall                string
	QuoteTimeInLong        int64
	Rho                    float32
	StrikePrice            float32
	Symbol                 string
	TheoreticalOptionValue float32
	TheoreticalVolatility  float32
	Theta                  float32
	TimeValue              float32
	TotalVolume            int32
	TradeTimeInLong        int64
	Vega                   float32
	Volatility             float32
}

type Response struct {
	CallExpDateMap   map[string]map[string][]StrikeQuote
	DaysToExpiration float32
	InterestRate     float32
	Interval         int64
	IsDelayed        bool
	IsIndex          bool
	PutExpDateMap    map[string]map[string][]StrikeQuote
	Status           string
	Strategy         string
	Symbol           string
	UnderlyingPrice  float32
	Volatility       float32
}
