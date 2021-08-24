// Copyright (c) 2021 deadc0de6

package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Webhook alert file struct
type Webhook struct {
	url     string
	headers map[string]string
	options map[string]string
}

// Notify notifies
func (a *Webhook) Notify(content string) error {
	data := map[string]string{"alert": content}
	jData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", a.url, bytes.NewBuffer(jData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range a.headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", "checkah")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("return status: %s", resp.Status)
	}

	return nil
}

// GetOptions returns this alert options
func (a *Webhook) GetOptions() map[string]string {
	return a.options
}

// GetDescription returns a description for this alert
func (a *Webhook) GetDescription() string {
	return fmt.Sprintf("alert to webhook \"%s\"", a.url)
}

// NewAlertWebhook creates a new file alert instance
func NewAlertWebhook(options map[string]string) (*Webhook, error) {
	url, ok := options["url"]
	if !ok {
		return nil, fmt.Errorf("\"url\" option required")
	}

	// get headers
	headers := make(map[string]string)
	for i := 0; i < 10; i++ {
		headerName := fmt.Sprintf("header%d", i)
		valueName := fmt.Sprintf("value%d", i)
		h, ok := options[headerName]
		if !ok {
			break
		}
		v, ok := options[valueName]
		if !ok {
			break
		}
		headers[h] = v
	}

	a := &Webhook{
		url:     url,
		headers: headers,
		options: options,
	}
	return a, nil
}
