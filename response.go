package kraken

type responsePublicTime struct {
	Errors []string                 `json:"error"`
	Result responsePublicTimeResult `json:"result"`
}

type responsePublicTimeResult struct {
	UnixTimestamp int64  `json:"unixtime"`
	RFC1123       string `json:"rfc1123"`
}

type responseSystemStatus struct {
	Errors []string                   `json:"error"`
	Result responseSystemStatusResult `json:"result"`
}

type responseSystemStatusResult struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type responsePublicAssets struct {
	Errors []string                                   `json:"error"`
	Result map[string]responsePublicAssetsResultAsset `json:"result"`
}

type responsePublicAssetsResultAsset struct {
	Class           string `json:"aclass"`
	AltName         string `json:"altname"`
	Decimals        int    `json:"decimals"`
	DisplayDecimals int    `json:"display_decimals"`
}

type responsePublicAssetPairs struct {
	Errors []string
	Result map[string]responsePublicAssetPairResultPair
}

type responsePublicAssetPairResultPair struct {
	AltName           string      `json:"altname"`
	WSName            string      `json:"wsname"`
	AClassBase        string      `json:"aclass_base"`
	Base              string      `json:"base"`
	AClassQuote       string      `json:"aclass_quote"`
	Quote             string      `json:"quote"`
	Lot               string      `json:"lot"`
	PairDecimals      int         `json:"pair_decimals"`
	LotDecimals       int         `json:"lot_decimals"`
	LotMultiplier     int         `json:"lot_multiplier"`
	LeverageBuy       []int       `json:"leverage_buy"`
	LeverageSell      []int       `json:"leverage_sell"`
	Fees              [][]float32 `json:"fees"`
	FeesMaker         [][]float32 `json:"fees_maker"`
	FeeVolumeCurrency string      `json:"fee_volume_currency"`
	MarginCalls       int         `json:"margin_call"`
	MarginStop        int         `json:"margin_stop"`
	OrderMin          float32     `json:"ordermin"`
}

type responsePublicTicker struct {
	Errors []string
	Result map[string]responsePublicTickerInformation
}

type responsePublicTickerInformation struct {
	Ask                        []string `json:"a"`
	Bid                        []string `json:"b"`
	LastClose                  []string `json:"c"`
	Volume                     []string `json:"v"`
	VolumeWeightedAveragePrice []string `json:"p"`
	NumberOfTrades             []uint64 `json:"t"`
	Low                        []string `json:"l"`
	High                       []string `json:"h"`
	Open                       string   `json:"o"`
}

type responsePublicOHLC struct {
	Errors []string
	Result map[string]interface{}
}

type responsePublicOHLCValue struct {
	Timestamp uint64
	Open      string
	High      string
	Low       string
	Close     string
	Volume    string
	Count     uint64
}

type responsePublicOrderBook struct {
	Result map[string]responsePublicOrderBookResultAskBid `json:"result"`
	Error  []string                                       `json:"error"`
}

type responsePublicOrderBookResultAskBid struct {
	Asks [][]interface{} `json:"asks"`
	Bids [][]interface{} `json:"bids"`
}

type responsePublicRecentTrades struct {
	Error  []string               `json:"error"`
	Result map[string]interface{} `json:"result"`
}

type responsePublicRecentSpreads struct {
	Error  []string               `json:"error"`
	Result map[string]interface{} `json:"result"`
}
