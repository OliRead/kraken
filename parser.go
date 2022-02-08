package kraken

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Parser handles parsing of response payloads from the Kraken API
// to a structured data type
type Parser struct{}

// Parse parse a payload
func (p *Parser) Parse(payload []byte, v interface{}) error {
	if v == nil {
		return fmt.Errorf("%w: cannot parse to nil pointer", ErrParse)
	}

	switch t := v.(type) {
	case *Time:
		return p.parsePublicTime(payload, t)
	case *SystemStatus:
		return p.parseSystemStatus(payload, t)
	case *Assets:
		return p.parseAssets(payload, t)
	case *AssetPairs:
		return p.parseAssetPairs(payload, t)
	case *Tickers:
		return p.parseTickers(payload, t)
	case *OHLCs:
		return p.parseOHLCs(payload, t)
	case *OrderBook:
		return p.parseOrderBook(payload, t)
	case *RecentTrades:
		return p.parseRecentTrades(payload, t)
	case *RecentSpreads:
		return p.parseRecentSpreads(payload, t)
	default:
		return fmt.Errorf("%w: unsupported type %s", ErrParse, reflect.TypeOf(v).String())
	}
}

func (p *Parser) parsePublicTime(payload []byte, parsed *Time) error {
	msg := responsePublicTime{}
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("%w:%s", ErrParse, err)
	}

	*parsed = Time{
		Errors:    p.parseErrors(msg.Errors),
		Timestamp: time.Unix(msg.Result.UnixTimestamp, 0),
	}

	return nil
}

func (p *Parser) parseSystemStatus(payload []byte, parsed *SystemStatus) error {
	msg := responseSystemStatus{}
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("%w:%s", ErrParse, err)
	}

	t, err := time.Parse(time.RFC3339, msg.Result.Timestamp)
	if err != nil {
		return fmt.Errorf("%w:%s", ErrParse, err)
	}

	*parsed = SystemStatus{
		Errors:    p.parseErrors(msg.Errors),
		Status:    msg.Result.Status,
		Timestamp: t.UTC(),
	}

	return nil
}

func (p *Parser) parseAssets(payload []byte, parsed *Assets) error {
	msg := responsePublicAssets{}
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("%w:%s", ErrParse, err)
	}

	assets := make(map[string]Asset)
	for name, asset := range msg.Result {
		assets[name] = Asset{
			Name:             name,
			Class:            asset.Class,
			AltName:          asset.AltName,
			Precision:        asset.Decimals,
			DisplayPrecision: asset.DisplayDecimals,
		}
	}

	*parsed = Assets{
		Errors: p.parseErrors(msg.Errors),
		Assets: assets,
	}

	return nil
}

func (p *Parser) parseAssetPairs(payload []byte, parsed *AssetPairs) error {
	msg := responsePublicAssetPairs{}
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("%w:%s", ErrParse, err)
	}

	pairs := make(map[string]AssetPair)
	for name, pair := range msg.Result {
		pairs[name] = AssetPair{
			AltName:           pair.AltName,
			WebSocketName:     pair.WSName,
			AssetClassBase:    pair.AClassBase,
			Base:              pair.Base,
			AssetClassQuote:   pair.AClassQuote,
			Quote:             pair.Quote,
			Lot:               pair.Lot,
			PairPrecision:     pair.PairDecimals,
			LotPrecision:      pair.LotDecimals,
			LotMultiplier:     pair.LotMultiplier,
			LeverageBuy:       pair.LeverageBuy,
			LeverageSell:      pair.LeverageSell,
			FeesTaker:         p.parseFees(pair.Fees),
			FeesMaker:         p.parseFees(pair.FeesMaker),
			FeeVolumeCurrency: pair.FeeVolumeCurrency,
			MarginCalls:       pair.MarginCalls,
			MarginStop:        pair.MarginStop,
			OrderMin:          pair.OrderMin,
		}
	}

	*parsed = AssetPairs{
		Errors: p.parseErrors(msg.Errors),
		Pairs:  pairs,
	}

	return nil
}

func (p *Parser) parseFees(fees [][]float32) []Fee {
	f := make([]Fee, len(fees))
	for i, fee := range fees {
		f[i] = Fee{
			Volume:     int(fee[0]),
			Percentage: fee[1],
		}
	}

	return f
}

func (p *Parser) parseTickers(payload []byte, parsed *Tickers) error {
	msg := responsePublicTicker{}
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("%w:%s", ErrParse, err)
	}

	tickers := make(map[string]Ticker, len(msg.Result))
	for pair, ticker := range msg.Result {
		t, err := p.parseTicker(pair, ticker)
		if err != nil {
			return err
		}

		tickers[pair] = t
	}

	*parsed = Tickers{
		Errors: p.parseErrors(msg.Errors),
		Result: tickers,
	}

	return nil
}

func (p *Parser) parseTicker(pair string, ticker responsePublicTickerInformation) (Ticker, error) {
	ask, err := p.parseAskBid(ticker.Ask[0], ticker.Ask[2], nil)
	if err != nil {
		return Ticker{}, err
	}

	bid, err := p.parseAskBid(ticker.Bid[0], ticker.Bid[2], nil)
	if err != nil {
		return Ticker{}, err
	}

	close, err := p.parseClose(ticker.LastClose[0], ticker.LastClose[1])
	if err != nil {
		return Ticker{}, err
	}

	volumeToday, err := decimal.NewFromString(ticker.Volume[0])
	if err != nil {
		return Ticker{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	volumeLast24Hours, err := decimal.NewFromString(ticker.Volume[1])
	if err != nil {
		return Ticker{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	volumeWeightedAveragePriceToday, err := decimal.NewFromString(ticker.VolumeWeightedAveragePrice[0])
	if err != nil {
		return Ticker{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	volumeWeightedAveragePriceLast24Hours, err := decimal.NewFromString(ticker.VolumeWeightedAveragePrice[1])
	if err != nil {
		return Ticker{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	lowToday, err := decimal.NewFromString(ticker.Low[0])
	if err != nil {
		return Ticker{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	lowLast24Hours, err := decimal.NewFromString(ticker.Low[1])
	if err != nil {
		return Ticker{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	highToday, err := decimal.NewFromString(ticker.High[0])
	if err != nil {
		return Ticker{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	highLast24Hours, err := decimal.NewFromString(ticker.High[1])
	if err != nil {
		return Ticker{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	open, err := decimal.NewFromString(ticker.Open)
	if err != nil {
		return Ticker{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	return Ticker{
		Pair:                                  pair,
		Ask:                                   ask,
		Bid:                                   bid,
		LastClose:                             close,
		VolumeToday:                           volumeToday,
		VolumeLast24Hours:                     volumeLast24Hours,
		VolumeWeightedAveragePriceToday:       volumeWeightedAveragePriceToday,
		VolumeWeightedAveragePriceLast24Hours: volumeWeightedAveragePriceLast24Hours,
		NumberOfTradesToday:                   ticker.NumberOfTrades[0],
		NumberOfTradesLast24Hours:             ticker.NumberOfTrades[1],
		LowToday:                              lowToday,
		LowLast24Hours:                        lowLast24Hours,
		HighToday:                             highToday,
		HighLast24Hours:                       highLast24Hours,
		Open:                                  open,
	}, nil
}

func (p *Parser) parseAskBid(price, volume string, timestamp *int64) (AskBid, error) {
	priceDecimal, err := decimal.NewFromString(price)
	if err != nil {
		return AskBid{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	volumeDecimal, err := decimal.NewFromString(volume)
	if err != nil {
		return AskBid{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	if timestamp == nil {
		return AskBid{
			Price:  priceDecimal,
			Volume: volumeDecimal,
		}, nil
	}

	return AskBid{
		Price:     priceDecimal,
		Volume:    volumeDecimal,
		Timestamp: time.Unix(*timestamp, 0).UTC(),
	}, nil
}

func (p *Parser) parseClose(price, volume string) (Close, error) {
	priceDecimal, err := decimal.NewFromString(price)
	if err != nil {
		return Close{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	volumeDecimal, err := decimal.NewFromString(volume)
	if err != nil {
		return Close{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	return Close{
		Price:  priceDecimal,
		Volume: volumeDecimal,
	}, nil
}

func (p *Parser) parseOHLCs(payload []byte, parsed *OHLCs) error {
	msg := responsePublicOHLC{}
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("%w:%s", ErrParse, err)
	}

	ohlcs := make(map[string][]OHLC)

	for k, v := range msg.Result {
		if k == "last" {
			last := new(big.Rat)
			last.SetFloat64(v.(float64))
			parsed.LastID = last.Num().Uint64()
			continue
		}

		pairOHLCs := []OHLC{}
		for _, ohlcValue := range v.([]interface{}) {
			ohlc, err := p.parseOHLC(ohlcValue.([]interface{}))
			if err != nil {
				return err
			}

			pairOHLCs = append(pairOHLCs, ohlc)
		}

		ohlcs[k] = pairOHLCs
	}

	parsed.Result = ohlcs

	return nil
}

func (p *Parser) parseOHLC(v []interface{}) (OHLC, error) {
	open, err := decimal.NewFromString(v[1].(string))
	if err != nil {
		return OHLC{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	high, err := decimal.NewFromString(v[2].(string))
	if err != nil {
		return OHLC{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	low, err := decimal.NewFromString(v[3].(string))
	if err != nil {
		return OHLC{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	close, err := decimal.NewFromString(v[4].(string))
	if err != nil {
		return OHLC{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	volumeWeightedAveragePrice, err := decimal.NewFromString(v[5].(string))
	if err != nil {
		return OHLC{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	volume, err := decimal.NewFromString(v[6].(string))
	if err != nil {
		return OHLC{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	return OHLC{
		Time:                       time.Unix((&big.Rat{}).SetFloat64(v[0].(float64)).Num().Int64(), 0).UTC(),
		Open:                       open,
		High:                       high,
		Low:                        low,
		Close:                      close,
		VolumeWeightedAveragePrice: volumeWeightedAveragePrice,
		Volume:                     volume,
		Count:                      (&big.Rat{}).SetFloat64(v[7].(float64)).Num().Uint64(),
	}, nil
}

func (p *Parser) parseOrderBook(payload []byte, parsed *OrderBook) error {
	msg := responsePublicOrderBook{}
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("%w:%s", ErrParse, err)
	}

	pairAsks := make(map[string][]AskBid)
	pairBids := make(map[string][]AskBid)

	for pair, askbids := range msg.Result {
		asks := []AskBid{}
		for _, ask := range askbids.Asks {
			price := decimal.NewFromFloat(ask[0].(float64))
			volume := decimal.NewFromFloat(ask[1].(float64))
			timestamp := decimal.NewFromFloat(ask[2].(float64)).IntPart()

			a := AskBid{
				Price:     price,
				Volume:    volume,
				Timestamp: time.Unix(timestamp, 0),
			}

			asks = append(asks, a)
		}
		pairAsks[pair] = asks

		bids := []AskBid{}
		for _, bid := range askbids.Bids {
			price := decimal.NewFromFloat(bid[0].(float64))
			volume := decimal.NewFromFloat(bid[1].(float64))
			timestamp := decimal.NewFromFloat(bid[2].(float64)).IntPart()

			b := AskBid{
				Price:     price,
				Volume:    volume,
				Timestamp: time.Unix(timestamp, 0),
			}

			bids = append(bids, b)
		}
		pairBids[pair] = bids
	}

	*parsed = OrderBook{
		Errors: p.parseErrors(msg.Error),
		Asks:   pairAsks,
		Bids:   pairBids,
	}

	return nil
}

func (p *Parser) parseRecentTrades(payload []byte, parsed *RecentTrades) error {
	msg := responsePublicRecentTrades{}
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("%w:%s", ErrParse, err)
	}

	trades := make(map[string][]RecentTrade)

	for k, v := range msg.Result {
		if k == "last" {
			lastID, err := strconv.ParseUint(v.(string), 10, 64)
			if err != nil {
				return fmt.Errorf("%w:%s", ErrParse, err)
			}

			parsed.LastID = lastID
			continue
		}

		pairTrades := []RecentTrade{}
		for _, tradeValue := range v.([]interface{}) {
			trade, err := p.parseRecentTrade(tradeValue.([]interface{}))
			if err != nil {
				return err
			}

			pairTrades = append(pairTrades, trade)
		}

		trades[k] = pairTrades
	}

	parsed.Trades = trades
	parsed.Errors = p.parseErrors(msg.Error)

	return nil
}

func (p *Parser) parseRecentTrade(v []interface{}) (RecentTrade, error) {
	price, err := decimal.NewFromString(v[0].(string))
	if err != nil {
		return RecentTrade{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	volume, err := decimal.NewFromString(v[1].(string))
	if err != nil {
		return RecentTrade{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	orderTime := decimal.NewFromFloat(v[2].(float64))
	if err != nil {
		return RecentTrade{}, fmt.Errorf("%w:%s", ErrParse, err)
	}

	orderAction := v[3].(string)
	orderType := v[4].(string)
	misc := v[5].(string)

	trade := RecentTrade{
		Price:  price,
		Volume: volume,
		// TODO get microseconds working properly
		Time:          time.Unix(orderTime.IntPart(), 0),
		Miscellaneous: misc,
	}

	switch orderAction {
	case "b":
		trade.Action = OrderActionBuy
	case "s":
		trade.Action = OrderActionSell
	default:
		trade.Action = OrderActionUnknown
	}

	switch orderType {
	case "l":
		trade.Type = OrderTypeLimit
	case "m":
		trade.Type = OrderTypeMarket
	default:
		trade.Type = OrderTypeUnknown
	}

	return trade, nil
}

func (p *Parser) parseRecentSpreads(payload []byte, parsed *RecentSpreads) error {
	msg := responsePublicRecentSpreads{}
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("%w: %s", ErrParse, err)
	}

	spreads := make(map[string][]Spread)
	for k, v := range msg.Result {
		if k == "last" {
			last := new(big.Rat)
			last.SetFloat64(v.(float64))
			parsed.LastID = last.Num().Uint64()
			continue
		}

		pairSpreads := []Spread{}
		for _, spreadValue := range v.([]interface{}) {
			spread, err := p.parseRecentSpread(spreadValue.([]interface{}))
			if err != nil {
				return err
			}

			pairSpreads = append(pairSpreads, spread)
		}

		spreads[k] = pairSpreads
	}

	parsed.Spreads = spreads
	parsed.Errors = p.parseErrors(msg.Error)

	return nil
}

func (p *Parser) parseRecentSpread(v []interface{}) (Spread, error) {
	timestamp := decimal.NewFromFloat(v[0].(float64)).IntPart()

	bid, err := decimal.NewFromString(v[1].(string))
	if err != nil {
		return Spread{}, fmt.Errorf("%w: %s", ErrParse, err)
	}

	ask, err := decimal.NewFromString(v[2].(string))
	if err != nil {
		return Spread{}, fmt.Errorf("%w: %s", ErrParse, err)
	}

	return Spread{
		Timestamp: time.Unix(timestamp, 0),
		Bid:       bid,
		Ask:       ask,
	}, nil
}

func (p *Parser) parseErrors(errStrings []string) []error {
	if len(errStrings) == 0 {
		return nil
	}

	errs := make([]error, len(errStrings))
	for i, errString := range errStrings {
		errParts := strings.SplitN(errString, ":", 2)

		switch errParts[0] {
		case "EGeneral":
			errs[i] = fmt.Errorf("%w:%s", ErrGeneral, errParts[1])
		case "EAPI":
			errs[i] = fmt.Errorf("%w:%s", ErrAPI, errParts[1])
		case "EQuery":
			errs[i] = fmt.Errorf("%w:%s", ErrQuery, errParts[1])
		case "EOrder":
			errs[i] = fmt.Errorf("%w:%s", ErrOrder, errParts[1])
		case "ETrade":
			errs[i] = fmt.Errorf("%w:%s", ErrTrade, errParts[1])
		case "EFunding":
			errs[i] = fmt.Errorf("%w:%s", ErrFunding, errParts[1])
		case "EService":
			errs[i] = fmt.Errorf("%w:%s", ErrService, errParts[1])
		case "ESession":
			errs[i] = fmt.Errorf("%w:%s", ErrSession, errParts[1])
		default:
			errs[i] = fmt.Errorf("%w:%s", ErrAPIUnknown, errString)
		}
	}

	return errs
}
