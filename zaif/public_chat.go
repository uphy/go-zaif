package zaif

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"strconv"

	"golang.org/x/net/websocket"
)

type (
	ChatClient struct {
		ws *websocket.Conn
	}
	ChatPacket struct {
		Code  ChatCode
		Value interface{}
	}
	ChatCode    int
	JoinMessage struct {
		Status   string
		Hash     string
		Nickname string
		Icon     string
	}
	ChatMessage struct {
		Likes           int
		Mona            Amount
		XEM             Amount
		BTC             Amount
		Status          string
		MessageHTMLFull string `json:"message_html_full"`
		MessageHTML     string `json:"message_html"`
		MessageHTMLLink string `json:"message_html_link"`
		Message         string
		NS              string
		Timestamp       int64
		ID              string
		Hash            string
		Icon            string
		Nickname        string
	}
	ChatHistory []ChatMessage
	ChatChannel string
)

const (
	ChatCodeConfig     ChatCode    = 0
	ChatCodeUnknown    ChatCode    = 40
	ChatCodeNormal     ChatCode    = 42
	ChatCodePing       ChatCode    = 2
	ChatCodePong       ChatCode    = 3
	ChatDebug          bool        = false
	ChatChannelDefault ChatChannel = "/"
	ChatChannelBotUser ChatChannel = "/botuser"
)

var (
	ErrChatUnexpectedCode = errors.New("unexpected code")
)

func (p *PublicAPI) Chat(channel ChatChannel) (<-chan ChatMessage, <-chan error) {
	err := make(chan error, 10)
	msg := make(chan ChatMessage, 100)

	client, e := NewChatClient()
	if e != nil {
		err <- e
	}
	if e := client.Join(channel); e != nil {
		err <- e
	}
	m, e := client.History()
	if e != nil {
		err <- e
	}
	go func() {
		defer close(msg)
		defer close(err)
		for _, v := range m {
			msg <- v
		}
		client.Start(msg, err)
	}()
	return msg, err
}

func NewChatClient() (*ChatClient, error) {
	ws, err := websocket.Dial("wss://chat2.zaif.jp:8080/socket.io/?EIO=3&transport=websocket", "", "https://ws.zaif.jp:8888")
	if err != nil {
		return nil, err
	}
	client := &ChatClient{ws}

	// receive config
	m, err := client.receive()
	if err != nil {
		return nil, err
	}
	if m.Code != ChatCodeConfig {
		return nil, ErrChatUnexpectedCode
	}
	// receive unknown
	m, err = client.receive()
	if err != nil {
		return nil, err
	}
	if m.Code != ChatCodeUnknown {
		return nil, ErrChatUnexpectedCode
	}
	// ping-pong
	go func() {
		ticker := time.NewTicker(time.Second * 9)
	l:
		for {
			select {
			case _, ok := <-ticker.C:
				if ok {
					if e := client.send(ChatCodePing, nil); e != nil {
						log.Println("ping failure")
						return
					}
				} else {
					break l
				}
			}
		}
	}()
	return client, nil
}

func (c *ChatClient) Join(channel ChatChannel) error {
	if channel == "" {
		channel = "/"
	}
	// join
	if err := c.send(ChatCodeNormal, []string{"join_channel", string(channel)}); err != nil {
		return err
	}
	_, err := c.receive()
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatClient) History() (ChatHistory, error) {
	if e := c.send(ChatCodeNormal, []string{"history", "/"}); e != nil {
		return nil, e
	}
	m, e := c.receive()
	if e != nil {
		return nil, e
	}
	var h ChatHistory
	if err := c.cast(m.Value.([]interface{})[1], &h); err != nil {
		return nil, err
	}
	return h, nil
}

func (c *ChatClient) cast(v interface{}, as interface{}) error {
	d, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, as)
}

func (c *ChatClient) Start(msg chan<- ChatMessage, err chan<- error) (<-chan ChatMessage, <-chan error) {
	for {
		m, e := c.receive()
		if e != nil {
			err <- e
		} else {
			if m.Code == ChatCodePong {
				continue
			}
		}
		switch m.Code {
		case ChatCodeNormal:
			v := m.Value.([]interface{})
			messageType := v[0]
			value := v[1]
			switch messageType {
			case "say":
				var message ChatMessage
				if e := c.cast(value, &message); e != nil {
					err <- e
				} else {
					msg <- message
				}
			case "join":
				// ignore
			case "change":
				// not supported yet
			}
		case ChatCodePong:
		}
	}
}

func (c *ChatClient) send(code ChatCode, value interface{}) error {
	m := ChatPacket{code, value}
	msg, err := m.message()
	if err != nil {
		return err
	}
	if ChatDebug {
		log.Println("SEND: " + msg)
	}
	return websocket.Message.Send(c.ws, msg)
}

func (c *ChatClient) receive() (*ChatPacket, error) {
	var v string
	if err := websocket.Message.Receive(c.ws, &v); err != nil {
		return nil, err
	}
	if ChatDebug {
		log.Println("RECV: " + v)
	}
	code := 0
	i := 0
	for {
		if i >= len(v) {
			return &ChatPacket{ChatCode(code), nil}, nil
		}
		n, err := strconv.ParseInt(string(v[i]), 10, 0)
		if err != nil {
			break
		}
		code *= 10
		code += int(n)
		i++
	}
	var j interface{}
	if err := json.Unmarshal([]byte(v[i:]), &j); err != nil {
		return nil, err
	}
	return &ChatPacket{ChatCode(code), j}, nil
}

func (c *ChatPacket) message() (string, error) {
	data, err := json.Marshal(c.Value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d%s", c.Code, string(data)), nil
}

func (m *ChatMessage) Time() time.Time {
	return time.Unix(m.Timestamp, 0)
}
