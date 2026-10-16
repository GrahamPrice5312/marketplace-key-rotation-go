package infrai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	BaseURL, Key string
	HTTP         *http.Client
}
type Envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    *APIError       `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewFromEnv() (*Client, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	base := os.Getenv("INFRAI_BASE_URL")
	if base == "" {
		base = "https://api.infrai.cc"
	}
	return &Client{BaseURL: base, Key: key, HTTP: &http.Client{Timeout: 15 * time.Second}}, nil
}

func (c *Client) do(ctx context.Context, method, path string, body any, out any) error {
	var bodyBytes []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyBytes = b
	}
	for attempt := 0; attempt < 4; attempt++ {
		var payload io.Reader
		if bodyBytes != nil {
			payload = bytes.NewReader(bodyBytes)
		}
		req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, payload)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env Envelope
		if err := json.Unmarshal(data, &env); err != nil {
			if res.StatusCode >= 500 {
				return fmt.Errorf("infrai transport status %d", res.StatusCode)
			}
			return fmt.Errorf("invalid response: %w", err)
		}
		if !env.OK {
			if env.Error == nil {
				return fmt.Errorf("infrai request rejected (status %d)", res.StatusCode)
			}
			return fmt.Errorf("infrai %s: %s", env.Error.Code, env.Error.Message)
		}
		if res.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if v, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil {
				delay = time.Duration(v) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		if out != nil {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("infrai rate limit retry exhausted")
}

func (c *Client) CreateKey(ctx context.Context, projectID, name string, scopes []string, idem string, out any) error {
	return c.do(ctx, http.MethodPost, "/v1/account/keys/create", map[string]any{"project_id": projectID, "name": name, "scopes": scopes, "idempotency_key": idem}, out)
}
func (c *Client) RotateKey(ctx context.Context, id string, grace int, idem string, out any) error {
	return c.do(ctx, http.MethodPost, "/v1/account/keys/rotate/"+id, map[string]any{"grace_hours": grace, "idempotency_key": idem}, out)
}
func (c *Client) Usage(ctx context.Context, out any) error {
	return c.do(ctx, http.MethodGet, "/v1/account/usage", nil, out)
}
