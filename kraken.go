package kraken

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// Client handles requesting data from a Kraken API and parsing it
// to the relative data type
type Client interface {
	Time(ctx context.Context) (Time, error)
	Status(ctx context.Context) (SystemStatus, error)
	Assets(ctx context.Context) (Assets, error)
	AssetPairs(ctx context.Context, info AssetPairInfo, pairs ...string) (AssetPairs, error)
	OHLC(ctx context.Context, interval OHLCInterval, since *uint64, pairs ...string) (OHLCs, error)
	OrderBook(ctx context.Context, count uint, pairs ...string) (OrderBook, error)
	RecentTrades(ctx context.Context, since *uint64, pairs ...string) (RecentTrades, error)
	RecentSpreads(ctx context.Context, pairs []string, since *uint64) (RecentSpreads, error)
}

// Time a parsed response from the "/public/Time" API endpoint
type Time struct {
	Errors    []error
	Timestamp time.Time
}

// SystemStatus a parsed response from the "/public/SystemStatus" API endpoint
type SystemStatus struct {
	Errors    []error
	Status    string
	Timestamp time.Time
}

// Assets a parsed response from the "/public/Assets" API endpoint
type Assets struct {
	Errors []error
	Assets map[string]Asset
}

// Asset a single parsed asset from the "/public/Assets" API endpoint
type Asset struct {
	Name             string
	Class            string
	AltName          string
	Precision        int
	DisplayPrecision int
}

// AssetPairs a parsed response from the "/public/AssetPairs" API endpoint
type AssetPairs struct {
	Errors []error
	Pairs  map[string]AssetPair
}

// AssetPair a single parsed asset pair from the "/public/AssetPairs" API endpoint
type AssetPair struct {
	AltName           string
	WebSocketName     string
	AssetClassBase    string
	Base              string
	AssetClassQuote   string
	Quote             string
	Lot               string
	PairPrecision     int
	LotPrecision      int
	LotMultiplier     int
	LeverageBuy       []int
	LeverageSell      []int
	FeesTaker         []Fee
	FeesMaker         []Fee
	FeeVolumeCurrency string
	MarginCalls       int
	MarginStop        int
	OrderMin          float32
}

// Fee a single parsed fee from the from the "/public/AssetPairs" API endpoint
type Fee struct {
	Volume     int
	Percentage float32
}

// Tickers a parsed response from the "/public/Ticker" API endpoint
type Tickers struct {
	Errors []error
	Result map[string]Ticker
}

// Ticker a single parsed ticker from the "/public/Ticker" API endpoint
type Ticker struct {
	Pair                                  string
	Ask                                   AskBid
	Bid                                   AskBid
	LastClose                             Close
	VolumeToday                           decimal.Decimal
	VolumeLast24Hours                     decimal.Decimal
	VolumeWeightedAveragePriceToday       decimal.Decimal
	VolumeWeightedAveragePriceLast24Hours decimal.Decimal
	NumberOfTradesToday                   uint64
	NumberOfTradesLast24Hours             uint64
	LowToday                              decimal.Decimal
	LowLast24Hours                        decimal.Decimal
	HighToday                             decimal.Decimal
	HighLast24Hours                       decimal.Decimal
	Open                                  decimal.Decimal
}

// AskBid a single parsed ask bid value from the the "/public/Ticker" API endpoint
type AskBid struct {
	Price     decimal.Decimal
	Volume    decimal.Decimal
	Timestamp time.Time
}

// Close a single parsed Close value from the "/public/Ticker" API endpoint
type Close struct {
	Price  decimal.Decimal
	Volume decimal.Decimal
}

// OHLCs a parsed response from the "/public/OHLC" API endpoint
type OHLCs struct {
	Errors []error
	Result map[string][]OHLC
	LastID uint64
}

// OHLC a single parsed OHLC value from the "/public/OHLC" API endpoint
type OHLC struct {
	Time                       time.Time
	Open                       decimal.Decimal
	High                       decimal.Decimal
	Low                        decimal.Decimal
	Close                      decimal.Decimal
	Volume                     decimal.Decimal
	VolumeWeightedAveragePrice decimal.Decimal
	Count                      uint64
}

// OrderBook a parsed response from the "/public/OrderBook" API endpoint
type OrderBook struct {
	Errors []error
	Asks   map[string][]AskBid
	Bids   map[string][]AskBid
}

// RecentTrades a parsed response from the "/public/Trades" API endpoint
type RecentTrades struct {
	Errors []error
	Trades map[string][]RecentTrade
	LastID uint64
}

// RecentTrade a single parsed trade value from the "/public/Trade" API endpoint
type RecentTrade struct {
	Price         decimal.Decimal
	Volume        decimal.Decimal
	Time          time.Time
	Action        OrderAction
	Type          OrderType
	Miscellaneous string
}

// RecentSpreads a parsed respones from the "/public/Spread" API endpoint
type RecentSpreads struct {
	Errors  []error
	Spreads map[string][]Spread
	LastID  uint64
}

// Spread a single parsed spread value from the "/public/Spread" API endpoint
type Spread struct {
	Timestamp time.Time
	Bid       decimal.Decimal
	Ask       decimal.Decimal
}

// OrderAction an action of a trade, either buy or sell
type OrderAction byte

// String return a string value of the order action
func (t OrderAction) String() string {
	switch t {
	case OrderActionBuy:
		return "buy"
	case OrderActionSell:
		return "sell"
	default:
		return "unknown"
	}
}

const (
	// OrderActionBuy enum representing a buy order action
	OrderActionBuy = iota
	// OrderActionSell enum representing a sell order action
	OrderActionSell
	// OrderActionUnknown enum representing an unknown order action
	OrderActionUnknown
)

// OrderType a type of trade, either market or limit
type OrderType byte

// String return a string value of the order type
func (t OrderType) String() string {
	switch t {
	case OrderTypeMarket:
		return "market"
	case OrderTypeLimit:
		return "limit"
	default:
		return "unknown"
	}
}

const (
	// OrderTypeMarket enum representing a market order
	OrderTypeMarket = iota
	// OrderTypeLimit enum representing a limit order
	OrderTypeLimit
	// OrderTypeUnknown enum representing an unknown order action
	OrderTypeUnknown
)

// AssetPairInfo info values used in asset pair queries
type AssetPairInfo string

const (
	// AssetPairInfoInfo "info" value used in asset pair queries
	AssetPairInfoInfo = AssetPairInfo("info")
	// AssetPairInfoLeverage "info" value used in asset pair queries
	AssetPairInfoLeverage = AssetPairInfo("leverage")
	// AssetPairInfoFees "info" value used in asset pair queries
	AssetPairInfoFees = AssetPairInfo("fees")
	// AssetPairInfoMargin "info" value used in asset pair queries
	AssetPairInfoMargin = AssetPairInfo("margin")
)

// OHLCInterval interval value in OHLC queries
type OHLCInterval int

const (
	// OHLCIntervalMinute interval values in OHLC queries
	OHLCIntervalMinute = OHLCInterval(1)
	// OHLCInterval5Minutes interval values in OHLC queries
	OHLCInterval5Minutes = OHLCInterval(5)
	// OHLCInterval15Minutes interval values in OHLC queries
	OHLCInterval15Minutes = OHLCInterval(15)
	// OHLCInterval30Minutes interval values in OHLC queries
	OHLCInterval30Minutes = OHLCInterval(30)
	// OHLCIntervalHour interval values in OHLC queries
	OHLCIntervalHour = OHLCInterval(60)
	// OHLCInterval4Hour interval values in OHLC queries
	OHLCInterval4Hour = OHLCInterval(240)
	// OHLCIntervalDaily interval values in OHLC queries
	OHLCIntervalDaily = OHLCInterval(1440)
	// OHLCIntervalWeekly interval values in OHLC queries
	OHLCIntervalWeekly = OHLCInterval(10080)
	// OHLCInterval15Days interval values in OHLC queries
	OHLCInterval15Days = OHLCInterval(21600)
)
