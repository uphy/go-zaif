package zaif

import "net/url"

type (
	Funds   map[string]Amount
	Deposit map[string]Amount
	Rights  struct {
		Info         int `json:"info"`
		Trade        int `json:"trade"`
		Withdraw     int `json:"withdraw"`
		PersonalInfo int `json:"personal_info"`
	}
	Info2 struct {
		Funds      Funds   `json:"funds"`
		Deposit    Deposit `json:"deposit"`
		Rights     Rights  `json:"rights"`
		OpenOrders int     `json:"open_orders"`
		ServerTime int64   `json:"server_time"`
	}
)

func (c *PrivateAPI) Info2() (*Info2, error) {
	var r Info2
	err := c.requestWithRetry("get_info2", url.Values{}, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
