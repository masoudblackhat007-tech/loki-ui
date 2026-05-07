package loki

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client: کلاینت ساده Loki
type Client struct {
	baseURL string
	client  *http.Client
}

// baseURL مثل: http://localhost:3100
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Ready checks whether Loki is reachable and ready.
func (c *Client) Ready(ctx context.Context) error {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("parse baseURL: %w", err)
	}

	u.Path = "/ready"
	u.RawQuery = ""

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("build ready request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do ready request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("loki ready returned status %d", resp.StatusCode)
	}

	return nil
}

// QueryRange: کال به /loki/api/v1/query_range و تبدیل به []LogEntry
func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, limit int) ([]LogEntry, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse baseURL: %w", err)
	}

	u.Path = "/loki/api/v1/query_range"

	q := u.Query()
	q.Set("query", query)
	q.Set("start", fmt.Sprintf("%d", start.UnixNano()))
	q.Set("end", fmt.Sprintf("%d", end.UnixNano()))
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("direction", "BACKWARD") // جدیدترین لاگ‌ها اول
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("loki returned status %d", resp.StatusCode)
	}

	var qr QueryRangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&qr); err != nil {
		return nil, fmt.Errorf("decode loki response: %w", err)
	}

	if qr.Status != "success" {
		return nil, fmt.Errorf("loki status is %q", qr.Status)
	}

	entries := make([]LogEntry, 0, 256)

	for _, stream := range qr.Data.Result {
		labels := stream.Stream

		for _, pair := range stream.Values {
			if len(pair) != 2 {
				continue
			}

			tsStr := pair[0]
			line := pair[1]

			tsInt, err := strconv.ParseInt(tsStr, 10, 64)
			if err != nil {
				continue
			}
			ts := time.Unix(0, tsInt).UTC()

			entry := LogEntry{
				Timestamp: ts,
				Raw:       line,
				Labels:    labels,
				Parsed:    parseLaravelLog(line),
			}

			entries = append(entries, entry)
		}
	}

	return entries, nil
}

// Parse JSON لاگ لاراول
func parseLaravelLog(raw string) *LaravelLog {
	var l LaravelLog
	if err := json.Unmarshal([]byte(raw), &l); err != nil {
		return nil
	}
	return &l
}
