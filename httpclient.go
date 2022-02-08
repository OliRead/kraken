package kraken

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// HTTPClient used to interact with the Kraken API and return parsed responses
type HTTPClient struct {
	httpClient *http.Client
	parser     Parser
	dryRun     bool
	secret     string
	baseURL    string
}

// NewHTTPClient helper function for creating a new Kraken HTTPClient
func NewHTTPClient(opts ...HTTPClientOption) (*HTTPClient, error) {
	c := HTTPClient{
		httpClient: http.DefaultClient,
		baseURL:    "https://api.kraken.com/0",
		parser:     Parser{},
	}

	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// Time query the Kraken /public/time endpoint and return a parsed response
func (c *HTTPClient) Time(ctx context.Context) (Time, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/public/time", c.baseURL), nil)
	if err != nil {
		return Time{}, err
	}

	res, err := c.execute(req)
	if err != nil {
		return Time{}, err
	}

	defer res.Body.Close()
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Time{}, err
	}

	msg := Time{}
	if err := c.parser.Parse(payload, &msg); err != nil {
		return Time{}, err
	}

	return msg, err
}

// Status query the Kraken /public/SystemStatus endpoint and return a
// parsed response
func (c *HTTPClient) Status(ctx context.Context) (SystemStatus, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/public/SystemStatus", c.baseURL), nil)
	if err != nil {
		return SystemStatus{}, err
	}

	res, err := c.execute(req)
	if err != nil {
		return SystemStatus{}, err
	}

	defer res.Body.Close()
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return SystemStatus{}, err
	}

	msg := SystemStatus{}
	if err := c.parser.Parse(payload, &msg); err != nil {
		return SystemStatus{}, err
	}

	return msg, err
}

// Assets query the Kraken /public/Assets endpoint and return a parsed response
func (c *HTTPClient) Assets(ctx context.Context) (Assets, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/public/Assets", c.baseURL), nil)
	if err != nil {
		return Assets{}, err
	}

	res, err := c.execute(req)
	if err != nil {
		return Assets{}, err
	}

	defer res.Body.Close()
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Assets{}, err
	}

	msg := Assets{}
	if err := c.parser.Parse(payload, &msg); err != nil {
		return Assets{}, err
	}

	return msg, err
}

// AssetPairs query the Kraken /public/AssetPairs endpoint and return a parsed
// response
func (c *HTTPClient) AssetPairs(ctx context.Context, info AssetPairInfo, pairs ...string) (AssetPairs, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/public/AssetPairs", c.baseURL), nil)
	if err != nil {
		return AssetPairs{}, err
	}

	query := req.URL.Query()
	query["info"] = []string{string(info)}
	if len(pairs) != 0 {
		query["pairs"] = []string{strings.Join(pairs, ",")}
	}
	req.URL.RawQuery = query.Encode()

	res, err := c.execute(req)
	if err != nil {
		return AssetPairs{}, err
	}

	defer res.Body.Close()
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return AssetPairs{}, err
	}

	msg := AssetPairs{}
	if err := c.parser.Parse(payload, &msg); err != nil {
		return AssetPairs{}, err
	}

	return msg, err
}

// OHLC query the Kraken /public/OHLC endpoint and return a parsed
// response
func (c *HTTPClient) OHLC(ctx context.Context, interval OHLCInterval, since *uint64, pairs ...string) (OHLCs, error) {
	if len(pairs) == 0 {
		return OHLCs{}, fmt.Errorf("pairs are required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/public/OHLC", c.baseURL), nil)
	if err != nil {
		return OHLCs{}, err
	}

	query := req.URL.Query()
	query["pairs"] = []string{strings.Join(pairs, ",")}
	query["interval"] = []string{strconv.Itoa(int(interval))}

	if since != nil {
		query["since"] = []string{strconv.FormatUint(*since, 10)}
	}
	req.URL.RawQuery = query.Encode()

	res, err := c.execute(req)
	if err != nil {
		return OHLCs{}, err
	}

	defer res.Body.Close()
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return OHLCs{}, err
	}

	msg := OHLCs{}
	if err := c.parser.Parse(payload, &msg); err != nil {
		return OHLCs{}, err
	}

	return msg, err
}

// OrderBook query the Kraken /public/OrderBook endpoint and return a parsed
// response
func (c *HTTPClient) OrderBook(ctx context.Context, count uint, pairs ...string) (OrderBook, error) {
	if len(pairs) == 0 {
		return OrderBook{}, fmt.Errorf("pairs are required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/public/OrderBook", c.baseURL), nil)
	if err != nil {
		return OrderBook{}, err
	}

	query := req.URL.Query()
	query["pairs"] = []string{strings.Join(pairs, ",")}
	query["count"] = []string{strconv.FormatUint(uint64(count), 10)}
	req.URL.RawQuery = query.Encode()

	res, err := c.execute(req)
	if err != nil {
		return OrderBook{}, err
	}

	defer res.Body.Close()
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return OrderBook{}, err
	}

	msg := OrderBook{}
	if err := c.parser.Parse(payload, &msg); err != nil {
		return OrderBook{}, err
	}

	return msg, err
}

// RecentTrades query the Kraken /public/Trades endpoint and return a parsed
// response
func (c *HTTPClient) RecentTrades(ctx context.Context, since *uint64, pairs ...string) (RecentTrades, error) {
	if len(pairs) == 0 {
		return RecentTrades{}, fmt.Errorf("pairs are required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/public/Trades", c.baseURL), nil)
	if err != nil {
		return RecentTrades{}, err
	}

	query := req.URL.Query()
	query["pairs"] = []string{strings.Join(pairs, ",")}

	if since != nil {
		query["since"] = []string{strconv.FormatUint(*since, 10)}
	}

	res, err := c.execute(req)
	if err != nil {
		return RecentTrades{}, err
	}

	defer res.Body.Close()
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return RecentTrades{}, err
	}

	msg := RecentTrades{}
	if err := c.parser.Parse(payload, &msg); err != nil {
		return RecentTrades{}, err
	}

	return msg, err
}

// RecentSpreads query the Kraken /public/Spread endpoint and return a parsed
// response
func (c *HTTPClient) RecentSpreads(ctx context.Context, since *uint64, pairs ...string) (RecentSpreads, error) {
	if len(pairs) == 0 {
		return RecentSpreads{}, fmt.Errorf("pairs are required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/public/Spread", c.baseURL), nil)
	if err != nil {
		return RecentSpreads{}, err
	}
	query := req.URL.Query()
	query["pairs"] = []string{strings.Join(pairs, ",")}

	if since != nil {
		query["since"] = []string{strconv.FormatUint(*since, 10)}
	}

	res, err := c.execute(req)
	if err != nil {
		return RecentSpreads{}, err
	}

	defer res.Body.Close()
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return RecentSpreads{}, err
	}

	msg := RecentSpreads{}
	if err := c.parser.Parse(payload, &msg); err != nil {
		return RecentSpreads{}, err
	}

	return msg, err
}

func (c *HTTPClient) signature(path string, query url.Values) (string, error) {
	decodedSecret, err := base64.StdEncoding.DecodeString(c.secret)
	if err != nil {
		return "", err
	}

	sha := sha256.New()
	if _, err := sha.Write([]byte(query.Get("nonce") + query.Encode())); err != nil {
		return "", err
	}
	shaSum := sha.Sum(nil)

	mac := hmac.New(sha512.New, decodedSecret)
	if _, err := mac.Write(append([]byte(path), shaSum...)); err != nil {
		return "", err
	}
	macSum := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(macSum), nil
}

func (c *HTTPClient) execute(req *http.Request) (*http.Response, error) {
	if c.dryRun {
		return nil, ErrDryRun
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrNetwork, err)
	}

	return res, nil
}
