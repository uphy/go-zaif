package zaif

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type (
	PublicAPI struct {
		Retry int
	}
)

func NewPublicAPI() *PublicAPI {
	return &PublicAPI{10}
}

func (p *PublicAPI) getWithRetry(path string, v interface{}) error {
	for i := 0; i < p.Retry; i++ {
		err := p.get(path, v)
		if err != nil {
			return nil
		}
		if err == Err502 {
			continue
		}
		return err
	}
	return ErrTooManyRetry
}

func (p *PublicAPI) get(path string, v interface{}) error {
	resp, err := http.Get(baseURL + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 502 {
		return Err502
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
