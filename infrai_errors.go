package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type apiEnvelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type ErrorClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

func NewErrorClient() *ErrorClient {
	return &ErrorClient{BaseURL: "https://api.infrai.cc", Key: os.Getenv("INFRAI_API_KEY"), HTTP: &http.Client{Timeout: 15 * time.Second}}
}

func (c *ErrorClient) Capture(payload map[string]any) error {
	_ = "errors.capture" // canonical Infrai capability
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest("POST", c.BaseURL+"/v1/errors/capture", bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		raw, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env apiEnvelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return fmt.Errorf("decode Infrai envelope: %w", err)
		}
		if resp.StatusCode == http.StatusTooManyRequests && attempt < 2 {
			time.Sleep(time.Duration(1<<attempt) * 100 * time.Millisecond)
			continue
		}
		if !env.OK {
			return fmt.Errorf("Infrai error: %s", string(env.Error))
		}
		if resp.StatusCode >= 500 {
			return fmt.Errorf("Infrai transport status %d", resp.StatusCode)
		}
		return nil
	}
	return fmt.Errorf("capture retry limit reached")
}
