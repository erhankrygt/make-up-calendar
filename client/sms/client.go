package smsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client interface {
	Send(ctx context.Context, messages []Message) bool
}

type client struct {
	Source   string `json:"source_addr"`
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
}

// NewClient creates and returns mail client
func NewClient(source, username, password, url string) Client {
	cli := &client{
		Source:   source,
		Username: username,
		Password: password,
		URL:      url,
	}

	return cli
}

type Message struct {
	PhoneNumber string `json:"dest"`
	Text        string `json:"msg"`
}

type request struct {
	client   client
	Messages []Message `json:"messages"`
}

func (c client) Send(_ context.Context, messages []Message) bool {
	r := request{
		client:   c,
		Messages: messages,
	}

	payload, err := json.Marshal(r)
	if err != nil {
		return false
	}

	httpClient := &http.Client{}

	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(payload))
	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Message sent successfully.")
		return true
	}

	return false
}
