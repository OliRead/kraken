package kraken

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	operationCount    *prometheus.CounterVec
	operationDuration *prometheus.HistogramVec
	errorCount        *prometheus.CounterVec
)

// InstrumentationClient handles prometheus metrics for calls to
// client functins
type InstrumentationClient struct {
	inner Client
}

// NewInstrumentationClient helper function for creating a new instrumenation
// client to add prometheus metrics
func NewInstrumentationClient(inner Client) InstrumentationClient {
	return InstrumentationClient{
		inner: inner,
	}
}

// Time handles prometheus metrics for client Time function
func (c *InstrumentationClient) Time(ctx context.Context) (Time, error) {
	timer := prometheus.NewTimer(
		operationDuration.WithLabelValues("Time"),
	)
	defer timer.ObserveDuration()

	operationCount.WithLabelValues("Time").Inc()

	v, err := c.inner.Time(ctx)
	if err != nil {
		errorCount.WithLabelValues("Time").Inc()
	}

	return v, err
}

// Status handles prometheus metrics for client Status function
func (c *InstrumentationClient) Status(ctx context.Context) (SystemStatus, error) {
	timer := prometheus.NewTimer(
		operationDuration.WithLabelValues("Status"),
	)
	defer timer.ObserveDuration()

	operationCount.WithLabelValues("Status").Inc()

	v, err := c.inner.Status(ctx)
	if err != nil {
		errorCount.WithLabelValues("Status").Inc()
	}

	return v, err
}

// Assets handles prometheus metrics for client Assets function
func (c *InstrumentationClient) Assets(ctx context.Context) (Assets, error) {
	timer := prometheus.NewTimer(
		operationDuration.WithLabelValues("Assets"),
	)
	defer timer.ObserveDuration()

	operationCount.WithLabelValues("Assets").Inc()

	v, err := c.inner.Assets(ctx)
	if err != nil {
		errorCount.WithLabelValues("Assets").Inc()
	}

	return v, err
}

// AssetPairs handles prometheus metrics for client AssetPairs function
func (c *InstrumentationClient) AssetPairs(ctx context.Context, info AssetPairInfo, pairs ...string) (AssetPairs, error) {
	timer := prometheus.NewTimer(
		operationDuration.WithLabelValues("AssetPairs"),
	)
	defer timer.ObserveDuration()

	operationCount.WithLabelValues("AssetPairs").Inc()

	v, err := c.inner.AssetPairs(ctx, info, pairs...)
	if err != nil {
		errorCount.WithLabelValues("AssetPairs").Inc()
	}

	return v, err
}

// OHLC handles prometheus metrics for client OHLC function
func (c *InstrumentationClient) OHLC(ctx context.Context, interval OHLCInterval, since *uint64, pairs ...string) (OHLCs, error) {
	timer := prometheus.NewTimer(
		operationDuration.WithLabelValues("OHLC"),
	)
	defer timer.ObserveDuration()

	operationCount.WithLabelValues("OHLC").Inc()

	v, err := c.inner.OHLC(ctx, interval, since, pairs...)
	if err != nil {
		errorCount.WithLabelValues("OHLC").Inc()
	}

	return v, err
}

// OrderBook handles prometheus metrics for client OrderBook function
func (c *InstrumentationClient) OrderBook(ctx context.Context, count uint, pairs ...string) (OrderBook, error) {
	timer := prometheus.NewTimer(
		operationDuration.WithLabelValues("OrderBook"),
	)
	defer timer.ObserveDuration()

	operationCount.WithLabelValues("OrderBook").Inc()

	v, err := c.inner.OrderBook(ctx, count, pairs...)
	if err != nil {
		errorCount.WithLabelValues("OrderBook").Inc()
	}

	return v, err
}

// RecentTrades handles prometheus metrics for client RecentTrades function
func (c *InstrumentationClient) RecentTrades(ctx context.Context, since *uint64, pairs ...string) (RecentTrades, error) {
	timer := prometheus.NewTimer(
		operationDuration.WithLabelValues("RecentTrades"),
	)
	defer timer.ObserveDuration()

	operationCount.WithLabelValues("RecentTrades").Inc()

	v, err := c.inner.RecentTrades(ctx, since, pairs...)
	if err != nil {
		errorCount.WithLabelValues("RecentTrades").Inc()
	}

	return v, err
}

// RecentSpreads handles prometheus metrics for client RecentSpreads function
func (c *InstrumentationClient) RecentSpreads(ctx context.Context, since *uint64, pairs ...string) (RecentSpreads, error) {
	timer := prometheus.NewTimer(
		operationDuration.WithLabelValues("RecentSpreads"),
	)
	defer timer.ObserveDuration()

	operationCount.WithLabelValues("RecentSpreads").Inc()

	v, err := c.inner.RecentSpreads(ctx, pairs, since)
	if err != nil {
		errorCount.WithLabelValues("RecentSpreads").Inc()
	}

	return v, err
}
