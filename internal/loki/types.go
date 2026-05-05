package loki

import "time"

// ساختار ریسپانس Loki برای /loki/api/v1/query_range
type QueryRangeResponse struct {
	Status string      `json:"status"`
	Data   QueryResult `json:"data"`
}

type QueryResult struct {
	ResultType string         `json:"resultType"`
	Result     []StreamResult `json:"result"`
}

type StreamResult struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"` // [ [timestamp_ns, line], ... ]
}

// JSON لاگ لاراول (همون چیزی که از laravel-json-YYYY-MM-DD.log میاد)
type LaravelLog struct {
	Message   string         `json:"message"`
	Context   map[string]any `json:"context"`
	Level     int            `json:"level"`
	LevelName string         `json:"level_name"`
	Channel   string         `json:"channel"`
	Datetime  string         `json:"datetime"`
	Extra     map[string]any `json:"extra"`
}

// چیزی که داخل Go باهاش کار می‌کنیم
type LogEntry struct {
	Timestamp time.Time
	Raw       string
	Labels    map[string]string
	Parsed    *LaravelLog
}
