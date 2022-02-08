package kraken_test

import (
	"errors"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/oliread/kraken"
	"github.com/shopspring/decimal"
)

func TestParseInvalidDataType(t *testing.T) {
	p := kraken.Parser{}

	type invalidDataType struct{}
	msg := invalidDataType{}

	if err := p.Parse(nil, &msg); !errors.Is(err, kraken.ErrParse) {
		t.Fatal(err)
	}
}

func TestParseNilPointer(t *testing.T) {
	p := kraken.Parser{}

	if err := p.Parse(nil, nil); !errors.Is(err, kraken.ErrParse) {
		t.Fatal(err)
	}
}

func TestParseTime(t *testing.T) {
	tcs := []struct {
		name     string
		input    []byte
		expected kraken.Time
		err      error
	}{
		{
			name: "ValidPayload",
			input: []byte(`
			{
				"error":[],
				"result":{
					"unixtime":1643584726,
					"rfc1123":"Sun, 30 Jan 22 23:18:46 +0000"
				}
			}
			`),
			expected: kraken.Time{
				Errors:    nil,
				Timestamp: time.Unix(1643584726, 0),
			},
		},
	}

	p := kraken.Parser{}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			msg := kraken.Time{}
			if err := p.Parse(tc.input, &msg); err != tc.err {
				t.Fatal(err)
			}

			if diff := deep.Equal(tc.expected, msg); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestParseSystemStatus(t *testing.T) {
	tcs := []struct {
		name     string
		input    []byte
		expected kraken.SystemStatus
		err      error
	}{
		{
			name: "ValidPayload",
			input: []byte(`
			{
				"error":[],
				"result":{
					"status":"online",
					"timestamp":"2022-01-31T00:44:35Z"
				}
			}
			`),
			expected: kraken.SystemStatus{
				Errors:    nil,
				Timestamp: time.Unix(1643589875, 0).UTC(),
				Status:    "online",
			},
		},
	}

	p := kraken.Parser{}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			msg := kraken.SystemStatus{}
			if err := p.Parse(tc.input, &msg); err != tc.err {
				t.Fatal(err)
			}

			if diff := deep.Equal(tc.expected, msg); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestParseAsset(t *testing.T) {
	tcs := []struct {
		name     string
		input    []byte
		expected kraken.Assets
		err      error
	}{
		{
			name: "ValidPayload",
			input: []byte(`
			{
				"error": [],
				"result": {
					"NANO": {
						"aclass": "currency",
						"altname": "NANO",
						"decimals": 10,
						"display_decimals": 5
					},
					"ZUSD": {
						"aclass": "currency",
						"altname": "USD",
						"decimals": 4,
						"display_decimals": 2
					},
					"XXBT": {
						"aclass": "currency",
						"altname": "XBT",
						"decimals": 10,
						"display_decimals": 5
					}
				}
			}
			`),
			expected: kraken.Assets{
				Assets: map[string]kraken.Asset{
					"NANO": {
						Name:             "NANO",
						Class:            "currency",
						AltName:          "NANO",
						Precision:        10,
						DisplayPrecision: 5,
					},
					"ZUSD": {
						Name:             "ZUSD",
						Class:            "currency",
						AltName:          "USD",
						Precision:        4,
						DisplayPrecision: 2,
					},
					"XXBT": {
						Name:             "XXBT",
						Class:            "currency",
						AltName:          "XBT",
						Precision:        10,
						DisplayPrecision: 5,
					},
				},
			},
			err: nil,
		},
	}

	p := kraken.Parser{}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			msg := kraken.Assets{}
			if err := p.Parse(tc.input, &msg); err != tc.err {
				t.Fatal(err)
			}

			if diff := deep.Equal(tc.expected, msg); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestParseAssetPairs(t *testing.T) {
	tcs := []struct {
		name     string
		input    []byte
		expected kraken.AssetPairs
		err      error
	}{
		{
			name: "ValidPayload",
			input: []byte(`
			{
				"error": [],
				"result": {
					"XXBTZUSD": {
						"altname": "XBTUSD",
						"wsname": "XBT/USD",
						"aclass_base": "currency",
						"base": "XXBT",
						"aclass_quote": "currency",
						"quote": "ZUSD",
						"lot": "unit",
						"pair_decimals": 1,
						"lot_decimals": 8,
						"lot_multiplier": 1,
						"leverage_buy": [
							2,
							3,
							4,
							5
						],
						"leverage_sell": [
							2,
							3,
							4,
							5
						],
						"fees": [
							[
								0,
								0.26
							],
							[
								50000,
								0.24
							],
							[
								100000,
								0.22
							],
							[
								250000,
								0.2
							],
							[
								500000,
								0.18
							],
							[
								1000000,
								0.16
							],
							[
								2500000,
								0.14
							],
							[
								5000000,
								0.12
							],
							[
								10000000,
								0.1
							]
						],
						"fees_maker": [
							[
								0,
								0.16
							],
							[
								50000,
								0.14
							],
							[
								100000,
								0.12
							],
							[
								250000,
								0.1
							],
							[
								500000,
								0.08
							],
							[
								1000000,
								0.06
							],
							[
								2500000,
								0.04
							],
							[
								5000000,
								0.02
							],
							[
								10000000,
								0
							]
						],
						"fee_volume_currency": "ZUSD",
						"margin_call": 80,
						"margin_stop": 40,
						"ordermin": 0.0001
					}
				}
			}
			`),
			expected: kraken.AssetPairs{
				Pairs: map[string]kraken.AssetPair{
					"XXBTZUSD": {
						AltName:         "XBTUSD",
						WebSocketName:   "XBT/USD",
						AssetClassBase:  "currency",
						Base:            "XXBT",
						AssetClassQuote: "currency",
						Quote:           "ZUSD",
						Lot:             "unit",
						PairPrecision:   1,
						LotPrecision:    8,
						LotMultiplier:   1,
						LeverageBuy:     []int{2, 3, 4, 5},
						LeverageSell:    []int{2, 3, 4, 5},
						FeesTaker: []kraken.Fee{
							{Volume: 0, Percentage: 0.26},
							{Volume: 50000, Percentage: 0.24},
							{Volume: 100000, Percentage: 0.22},
							{Volume: 250000, Percentage: 0.2},
							{Volume: 500000, Percentage: 0.18},
							{Volume: 1000000, Percentage: 0.16},
							{Volume: 2500000, Percentage: 0.14},
							{Volume: 5000000, Percentage: 0.12},
							{Volume: 10000000, Percentage: 0.1},
						},
						FeesMaker: []kraken.Fee{
							{Volume: 0, Percentage: 0.16},
							{Volume: 50000, Percentage: 0.14},
							{Volume: 100000, Percentage: 0.12},
							{Volume: 250000, Percentage: 0.1},
							{Volume: 500000, Percentage: 0.08},
							{Volume: 1000000, Percentage: 0.06},
							{Volume: 2500000, Percentage: 0.04},
							{Volume: 5000000, Percentage: 0.02},
							{Volume: 10000000, Percentage: 0},
						},
						FeeVolumeCurrency: "ZUSD",
						MarginCalls:       80,
						MarginStop:        40,
						OrderMin:          0.0001,
					},
				},
			},
			err: nil,
		},
	}

	p := kraken.Parser{}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			msg := kraken.AssetPairs{}
			if err := p.Parse(tc.input, &msg); err != tc.err {
				t.Fatal(err)
			}

			if diff := deep.Equal(tc.expected, msg); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestParseTicker(t *testing.T) {
	tcs := []struct {
		name     string
		input    []byte
		expected kraken.Tickers
		err      error
	}{
		{
			name: "ValidPayload",
			input: []byte(`
			{
				"error": [],
				"result": {
					"XXBTZUSD": {
						"a": [
							"38659.6",
							"1",
							"1.000"
						],
						"b": [
							"38658.7",
							"1",
							"1.000"
						],
						"c": [
							"38658.9",
							"0.021208"
						],
						"v": [
							"3150.86186124",
							"3404.34671"
						],
						"p": [
							"38609.60189",
							"38601.37073"
						],
						"t": [
							24864,
							27336
						],
						"l": [
							"38050.00000",
							"38050.00000"
						],
						"h": [
							"39290.00000",
							"39290.00000"
						],
						"o": "38512.00000"
					}
				}
			}
			`),
			expected: kraken.Tickers{
				Result: map[string]kraken.Ticker{
					"XXBTZUSD": {
						Pair: "XXBTZUSD",
						Ask: kraken.AskBid{
							Price:  decimal.New(386596, -1),
							Volume: decimal.New(1, 0),
						},
						Bid: kraken.AskBid{
							Price:  decimal.New(386587, -1),
							Volume: decimal.New(1, 0),
						},
						LastClose: kraken.Close{
							Price:  decimal.New(386589, -1),
							Volume: decimal.New(21208, -6),
						},
						VolumeToday:                           decimal.New(315086186124, -8),
						VolumeLast24Hours:                     decimal.New(340434671, -5),
						VolumeWeightedAveragePriceToday:       decimal.New(3860960189, -5),
						VolumeWeightedAveragePriceLast24Hours: decimal.New(3860137073, -5),
						NumberOfTradesToday:                   uint64(24864),
						NumberOfTradesLast24Hours:             uint64(27336),
						LowToday:                              decimal.New(3805000000, -5),
						LowLast24Hours:                        decimal.New(3805000000, -5),
						HighToday:                             decimal.New(3929000000, -5),
						HighLast24Hours:                       decimal.New(3929000000, -5),
						Open:                                  decimal.New(3851200000, -5),
					},
				},
			},
		},
	}

	p := kraken.Parser{}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			msg := kraken.Tickers{}
			if err := p.Parse(tc.input, &msg); err != tc.err {
				t.Fatal(err)
			}

			if diff := deep.Equal(tc.expected, msg); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestParseOHLC(t *testing.T) {
	tcs := []struct {
		name     string
		input    []byte
		expected kraken.OHLCs
		err      error
	}{
		{
			name: "ValidPayload",
			input: []byte(`
			{
				"error":[],
				"result":{
					"XXBTZUSD":[
						[
							1643714160,
							"38311.6",
							"38343.7",
							"38311.6",
							"38343.7",
							"38320.8",
							"0.40716249",
							11
						]
					],
					"last":1643757240
				}
			}
			`),
			expected: kraken.OHLCs{
				Result: map[string][]kraken.OHLC{
					"XXBTZUSD": {
						{
							Time:                       time.Unix(1643714160, 0).UTC(),
							Open:                       decimal.New(383116, -1),
							High:                       decimal.New(383437, -1),
							Low:                        decimal.New(383116, -1),
							Close:                      decimal.New(383437, -1),
							VolumeWeightedAveragePrice: decimal.New(383208, -1),
							Volume:                     decimal.New(40716249, -8),
							Count:                      11,
						},
					},
				},
				LastID: uint64(1643757240),
			},
		},
	}

	p := kraken.Parser{}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			msg := kraken.OHLCs{}
			if err := p.Parse(tc.input, &msg); err != tc.err {
				t.Fatal(err)
			}

			if diff := deep.Equal(tc.expected, msg); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestParseOrderBook(t *testing.T) {
	tcs := []struct {
		name     string
		input    []byte
		expected kraken.OrderBook
		err      error
	}{
		{
			name: "ValidPayload",
			input: []byte(`
			{
				"error": [],
				"result": {
					"XXBTZUSD": {
						"asks": [
							[
								37639.4,
								0.002,
								1643832845
							]
						],
						"bids": [
							[
								37639.3,
								3.488,
								1643832845
							]
						]
					}
				}
			}
			`),
			expected: kraken.OrderBook{
				Asks: map[string][]kraken.AskBid{
					"XXBTZUSD": {
						{
							Price:     decimal.New(376394, -1),
							Volume:    decimal.New(2, -3),
							Timestamp: time.Unix(1643832845, 0),
						},
					},
				},
				Bids: map[string][]kraken.AskBid{
					"XXBTZUSD": {
						{
							Price:     decimal.New(376393, -1),
							Volume:    decimal.New(3488, -3),
							Timestamp: time.Unix(1643832845, 0),
						},
					},
				},
			},
		},
	}

	p := kraken.Parser{}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			msg := kraken.OrderBook{}
			if err := p.Parse(tc.input, &msg); err != tc.err {
				t.Fatal(err)
			}

			if diff := deep.Equal(tc.expected, msg); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestParseRecentTrades(t *testing.T) {
	tcs := []struct {
		name     string
		input    []byte
		expected kraken.RecentTrades
		err      error
	}{
		{
			name: "ValidPayload",
			input: []byte(`
			{
				"error":[],
				"result":{
					"XXBTZUSD":[
						["42428.00000","0.00109505",1644189769.9122,"b","l",""],
						["42436.50000","0.00098631",1644189769.9134,"b","l",""]
					],
					"last": "1644191265969108820"
				}
			}
			`),
			expected: kraken.RecentTrades{
				Trades: map[string][]kraken.RecentTrade{
					"XXBTZUSD": {
						{
							Price:  decimal.New(42428, 0),
							Volume: decimal.New(109505, -8),
							Time:   time.Unix(1644189769, 0).UTC(),
							Action: kraken.OrderActionBuy,
							Type:   kraken.OrderTypeLimit,
						},
						{
							Price:  decimal.New(424365, -1),
							Volume: decimal.New(98631, -8),
							Time:   time.Unix(1644189769, 0).UTC(),
							Action: kraken.OrderActionBuy,
							Type:   kraken.OrderTypeLimit,
						},
					},
				},
				LastID: 1644191265969108820,
			},
		},
	}

	p := kraken.Parser{}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			msg := kraken.RecentTrades{}
			if err := p.Parse(tc.input, &msg); err != tc.err {
				t.Fatal(err)
			}

			if diff := deep.Equal(tc.expected, msg); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestParseRecentSpread(t *testing.T) {
	tcs := []struct {
		name     string
		input    []byte
		expected kraken.RecentSpreads
		err      error
	}{
		{
			name: "ValidPayload",
			input: []byte(`
			{
				"error":[],
				"result":{
					"XXBTZUSD":[
						[1644356229,"44223.30000","44225.10000"]
					],
					"last":1644356424
				}
			}
			`),
			expected: kraken.RecentSpreads{
				Spreads: map[string][]kraken.Spread{
					"XXBTZUSD": {
						{
							Timestamp: time.Unix(1644356229, 0),
							Bid:       decimal.New(442233, -1),
							Ask:       decimal.New(442251, -1),
						},
					},
				},
				LastID: 1644356424,
			},
		},
	}

	p := kraken.Parser{}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			msg := kraken.RecentSpreads{}
			if err := p.Parse(tc.input, &msg); err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(tc.expected, msg); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	input := []byte(`
	{
		"error":[
			"EGeneral:test error",
			"EAPI:test error",
			"EQuery:test error",
			"EOrder:test error",
			"ETrade:test error",
			"EFunding:test error",
			"EService:test error",
			"ESession:test error",
			"unknown test error"
		],
		"result":{
			"unixtime":1644358183,
			"rfc1123":"Tue,  8 Feb 22 22:09:43 +0000"
		}
	}
	`)

	output := []error{
		kraken.ErrGeneral,
		kraken.ErrAPI,
		kraken.ErrQuery,
		kraken.ErrOrder,
		kraken.ErrTrade,
		kraken.ErrFunding,
		kraken.ErrService,
		kraken.ErrSession,
		kraken.ErrAPIUnknown,
	}

	msg := kraken.Time{}
	p := kraken.Parser{}
	if err := p.Parse(input, &msg); err != nil {
		t.Fatal(err)
	}

	for i, err := range msg.Errors {
		if !errors.Is(msg.Errors[i], output[i]) {
			t.Fatalf("EXPECTED: %s\nACTUAL: %s", output[i], errors.Unwrap(err))
		}
	}
}
