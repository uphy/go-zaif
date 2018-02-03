package zaif

import "net/url"

type (
	Info struct {
		Info2
		TradeCount int64 `json:"trade_count"`
	}
)

func (c *PrivateAPI) Info() (*Info, error) {
	var r Info
	err := c.requestWithRetry("get_info", url.Values{}, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
