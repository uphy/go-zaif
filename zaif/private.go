package zaif

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type (
	PrivateAPI struct {
		key       string
		secret    string
		nonce     int64
		Test      bool
		publicAPI *PublicAPI
		Retry     int
	}
)

func NewPrivateApi(key, secret string) *PrivateAPI {
	return &PrivateAPI{key, secret, time.Now().Unix(), false, NewPublicAPI(), 10}
}

func (z *PrivateAPI) requestWithRetry(method string, params url.Values, v interface{}) error {
	for i := 0; i < z.Retry; i++ {
		err := z.request(method, params, v)
		if err == Err502 {
			continue
		}
		return err
	}
	return ErrTooManyRetry
}

func (z *PrivateAPI) request(method string, params url.Values, v interface{}) error {
	type Data struct {
		Success int `json:"success"`
		Return  interface{}
		Error   string `json:"error"`
	}
	uri := "https://api.zaif.jp/tapi"
	params.Add("method", method)
	params.Add("nonce", strconv.FormatInt(z.nonce, 10))
	if z.Test {
		fmt.Println(params.Encode())
		return nil
	}
	z.nonce++

	encodedParams := params.Encode()
	req, _ := http.NewRequest("POST", uri, strings.NewReader(encodedParams))

	hash := hmac.New(sha512.New, []byte(z.secret))
	hash.Write([]byte(encodedParams))
	signature := hex.EncodeToString(hash.Sum(nil))

	req.Header.Add("Key", z.key)
	req.Header.Add("Sign", signature)
	client := new(http.Client)
	resp, err := client.Do(req)
	if resp.StatusCode == 502 {
		return Err502
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	data := Data{
		Return: v,
	}
	if err := json.Unmarshal(byteArray, &data); err != nil {
		return err
	}
	if data.Success != 1 {
		return errors.New("request failed: " + data.Error)
	}
	return nil
}
