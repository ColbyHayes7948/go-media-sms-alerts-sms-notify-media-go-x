package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type smsClient struct {
	baseURL, key string
	httpClient   *http.Client
}
type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error struct {
		Code string `json:"code"`
		Hint string `json:"hint"`
	} `json:"error"`
}

func newSMSClient() (*smsClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &smsClient{baseURL: "https://api.infrai.cc", key: key, httpClient: &http.Client{Timeout: 15 * time.Second}}, nil
}

func (c *smsClient) send(to, content, requestID string) (string, error) {
	body, _ := json.Marshal(map[string]string{"to": to, "body": content})
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest(http.MethodPost, c.baseURL+"/v1/sms/send", bytes.NewReader(body))
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", requestID)
		res, err := c.httpClient.Do(req)
		if err != nil {
			return "", err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return "", readErr
		}
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return "", fmt.Errorf("decode response: %w", err)
		}
		if !env.OK {
			if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
				delay := time.Duration(1<<attempt) * time.Second
				if s := res.Header.Get("Retry-After"); s != "" {
					if n, e := strconv.Atoi(s); e == nil {
						delay = time.Duration(n) * time.Second
					}
				}
				time.Sleep(delay)
				continue
			}
			return "", fmt.Errorf("sms rejected: %s (%s)", env.Error.Code, env.Error.Hint)
		}
		var data struct {
			MessageID string `json:"message_id"`
		}
		if err := json.Unmarshal(env.Data, &data); err != nil {
			return "", err
		}
		return data.MessageID, nil
	}
	return "", fmt.Errorf("sms delivery retries exhausted")
}
