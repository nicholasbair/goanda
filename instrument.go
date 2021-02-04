package goanda

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/instrument-ep/

import (
	"strconv"
	"time"
)

type Candle struct {
	Open  float64 `json:"o,string"`
	Close float64 `json:"c,string"`
	Low   float64 `json:"l,string"`
	High  float64 `json:"h,string"`
}

type Candles struct {
	Complete bool      `json:"complete"`
	Volume   int       `json:"volume"`
	Time     time.Time `json:"time"`
	Mid      Candle    `json:"mid"`
}

type BidAskCandles struct {
	Candles []struct {
		Ask struct {
			C float64 `json:"c,string"`
			H float64 `json:"h,string"`
			L float64 `json:"l,string"`
			O float64 `json:"o,string"`
		} `json:"ask"`
		Bid struct {
			C float64 `json:"c,string"`
			H float64 `json:"h,string"`
			L float64 `json:"l,string"`
			O float64 `json:"o,string"`
		} `json:"bid"`
		Complete bool      `json:"complete"`
		Time     time.Time `json:"time"`
		Volume   int       `json:"volume"`
	} `json:"candles"`
}

type InstrumentHistory struct {
	Instrument  string    `json:"instrument"`
	Granularity string    `json:"granularity"`
	Candles     []Candles `json:"candles"`
}
type InstrumentBidAksHistory struct {
	Instrument  string          `json:"instrument"`
	Granularity string          `json:"granularity"`
	Candles     []BidAskCandles `json:"candles"`
}

type Bucket struct {
	Price             string `json:"price"`
	LongCountPercent  string `json:"longCountPercent"`
	ShortCountPercent string `json:"shortCountPercent"`
}

type BrokerBook struct {
	Instrument  string    `json:"instrument"`
	Time        time.Time `json:"time"`
	Price       string    `json:"price"`
	BucketWidth string    `json:"bucketWidth"`
	Buckets     []Bucket  `json:"buckets"`
}

type InstrumentPricing struct {
	Time   time.Time `json:"time"`
	Prices []struct {
		Type string    `json:"type"`
		Time time.Time `json:"time"`
		Bids []struct {
			Price     float64 `json:"price,string"`
			Liquidity int     `json:"liquidity"`
		} `json:"bids"`
		Asks []struct {
			Price     float64 `json:"price,string"`
			Liquidity int     `json:"liquidity"`
		} `json:"asks"`
		CloseoutBid    float64 `json:"closeoutBid,string"`
		CloseoutAsk    float64 `json:"closeoutAsk,string"`
		Status         string  `json:"status"`
		Tradeable      bool    `json:"tradeable"`
		UnitsAvailable struct {
			Default struct {
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"default"`
			OpenOnly struct {
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"openOnly"`
			ReduceFirst struct {
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"reduceFirst"`
			ReduceOnly struct {
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"reduceOnly"`
		} `json:"unitsAvailable"`
		QuoteHomeConversionFactors struct {
			PositiveUnits string `json:"positiveUnits"`
			NegativeUnits string `json:"negativeUnits"`
		} `json:"quoteHomeConversionFactors"`
		Instrument string `json:"instrument"`
	} `json:"prices"`
}

func (c *OandaConnection) GetCandles(instrument string, count string, granularity string) (InstrumentHistory, error) {
	endpoint := "/instruments/" + instrument + "/candles?count=" + count + "&granularity=" + granularity
	candles, err := c.Request(endpoint)
	if err != nil {
		return InstrumentHistory{}, err
	}

	data := InstrumentHistory{}
	unmarshalJson(candles, &data)

	return data, nil
}

/*
	Gets candles by to and from time.

	param: instrument  string Symbol to query.
    param: count       int64  The number of candlesticks to return in the response. Can't use when to and from is set.
    param: granularity string The granularity of the candlesticks to fetch. i.e., S5, S15, M1, M15, M30, H1, D, W, M.
    param: to          string The start of the time range to fetch candlesticks for, represented in Unix representation.
    param: from        string The end of the time range to fetch candlesticks for, represented in Unix representation.
    param: smooth      bool   A smoothed candlestick uses the previous candle’s close price as its open price, while an un-smoothed candlestick uses the first price from its time range as its open price..

    return: InstrumentHistory, error
*/
func (c *OandaConnection) GetCandlesByTime(instrument string, granularity string, from string, to string, smooth bool) (InstrumentHistory, error) {
	endpoint := "/instruments/" + instrument + "/candles?" + "&granularity=" + granularity +
		"&from=" + from + "&to=" + to + "&smooth=" + strconv.FormatBool(smooth)
	candles, err := c.Request(endpoint)
	if err != nil {
		return InstrumentHistory{}, err
	}

	data := InstrumentHistory{}
	unmarshalJson(candles, &data)
	return data, nil
}

/*
	Gets candles by to and from time.

	param: instrument  string Symbol to query.
    param: count       int64  The number of candlesticks to return in the response. Can't use when to and from is set.
    param: granularity string The granularity of the candlesticks to fetch. i.e., S5, S15, M1, M15, M30, H1, D, W, M.
    param: to          string The start of the time range to fetch candlesticks for, represented in Unix representation.
    param: from        string The end of the time range to fetch candlesticks for, represented in Unix representation.
    param: smooth      bool   A smoothed candlestick uses the previous candle’s close price as its open price, while an un-smoothed candlestick uses the first price from its time range as its open price..

    return: InstrumentBidAksHistory, error
*/
func (c *OandaConnection) GetBidAsksCandlesByTime(instrument string, granularity string, from string, to string, smooth bool) (InstrumentBidAksHistory, error) {
	endpoint := "/instruments/" + instrument + "/candles?" + "&granularity=" + granularity + "&price=BA" +
		"&from=" + from + "&to=" + to + "&smooth=" + strconv.FormatBool(smooth)
	candles, err := c.Request(endpoint)
	if err != nil {
		return InstrumentBidAksHistory{}, err
	}

	data := InstrumentBidAksHistory{}
	unmarshalJson(candles, &data)
	return data, nil
}

func (c *OandaConnection) GetBidAskCandles(instrument string, count string, granularity string) (BidAskCandles, error) {
	endpoint := "/instruments/" + instrument + "/candles?count=" + count + "&granularity=" + granularity + "&price=BA"
	candles, err := c.Request(endpoint)
	if err != nil {
		return BidAskCandles{}, err
	}

	data := BidAskCandles{}
	unmarshalJson(candles, &data)

	return data, nil
}

func (c *OandaConnection) OrderBook(instrument string) (BrokerBook, error) {
	endpoint := "/instruments/" + instrument + "/orderBook"
	orderbook, err := c.Request(endpoint)
	if err != nil {
		return BrokerBook{}, err
	}

	data := BrokerBook{}
	unmarshalJson(orderbook, &data)

	return data, nil
}

func (c *OandaConnection) PositionBook(instrument string) (BrokerBook, error) {
	endpoint := "/instruments/" + instrument + "/positionBook"
	orderbook, err := c.Request(endpoint)
	if err != nil {
		return BrokerBook{}, err
	}

	data := BrokerBook{}
	unmarshalJson(orderbook, &data)

	return data, nil
}

func (c *OandaConnection) GetInstrumentPrice(instrument string) (InstrumentPricing, error) {
	endpoint := "/accounts/" + c.accountID + "/pricing?instruments=" + instrument
	pricing, err := c.Request(endpoint)
	if err != nil {
		return InstrumentPricing{}, err
	}

	data := InstrumentPricing{}
	unmarshalJson(pricing, &data)

	return data, nil
}

func (i *InstrumentHistory) ExtractClosed() []float64 {
	var list []float64
	for _, candle := range i.Candles {
		list = append(list, candle.Mid.Close)
	}

	return list
}

func (i *InstrumentHistory) ExtractOpen() []float64 {
	var list []float64
	for _, candle := range i.Candles {
		list = append(list, candle.Mid.Open)
	}

	return list
}

func (i *InstrumentHistory) ExtractHigh() []float64 {
	var list []float64
	for _, candle := range i.Candles {
		list = append(list, candle.Mid.High)
	}

	return list
}

func (i *InstrumentHistory) ExtractLow() []float64 {
	var list []float64
	for _, candle := range i.Candles {
		list = append(list, candle.Mid.Low)
	}

	return list
}

func (i *InstrumentHistory) ExtractVol() []int {
	var list []int
	for _, candle := range i.Candles {
		list = append(list, candle.Volume)
	}

	return list
}