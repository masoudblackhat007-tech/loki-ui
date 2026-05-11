package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	stdhtml "html"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"loki-ui/internal/loki"
)

type Handler struct {
	lokiClient *loki.Client
	tmpl       *template.Template
	loc        *time.Location
}

type DocsPageData struct {
	Title             string
	BodyHTML          template.HTML
	Page              int
	Total             int
	PrevPage          int
	NextPage          int
	HasPrev           bool
	HasNext           bool
	CurrentFile       string
	Lang              string
	LangLabel         string
	TextDir           string
	OppositeLang      string
	OppositeLangLabel string
}

// برای رندر سروری جدول /logs
type LogView struct {
	Time        string
	Level       string
	Service     string
	Route       string
	Method      string
	RequestID   string
	LogType     string
	Message     string
	RawJSON     string
	Status      int
	DurationMs  int
	ErrorCode   string
	SafeMessage string
}

// داده‌ای که به template می‌دیم
type LogsPageData struct {
	Logs      []LogView
	Total     int
	Range     string
	Limit     int
	Service   string
	Level     string
	Text      string
	RequestID string
}

// DTO برای خروجی JSON /api/logs و صفحه‌ی جزئیات
type LogDTO struct {
	Timestamp  int64  `json:"ts"`   // UTC nanoseconds (timezone-independent)
	Time       string `json:"time"` // formatted in UI timezone
	Level      string `json:"level"`
	Service    string `json:"service"`
	Route      string `json:"route"`
	Method     string `json:"method"`
	RequestID  string `json:"request_id"`
	LogType    string `json:"log_type"`
	Status     int    `json:"status"`
	DurationMs int    `json:"duration_ms"`
	Message    string `json:"message"`

	// ریسپانس وب‌سرویس (از http_request)
	ResponseSuccess   *bool          `json:"response_success,omitempty"`
	ResponseMessage   string         `json:"response_message,omitempty"`
	ResponseErrorCode string         `json:"response_error_code,omitempty"`
	ResponseBody      map[string]any `json:"response_body,omitempty"`

	// اطلاعات خطا
	ErrorCode   string            `json:"error_code,omitempty"`
	SafeMessage string            `json:"safe_message,omitempty"`
	Context     map[string]any    `json:"context"`
	Raw         string            `json:"raw"`
	Labels      map[string]string `json:"labels,omitempty"`
	Extra       map[string]any    `json:"extra,omitempty"`

	// فقط برای صفحه HTML جزئیات
	RawPretty     string            `json:"-"`
	RequestJSON   string            `json:"-"`
	ResponseJSON  string            `json:"-"`
	QueryJSON     string            `json:"-"`
	DBQueries     []DBQueryDTO      `json:"db_queries,omitempty"`
	DBQueriesJSON string            `json:"-"`
	Upstreams     []UpstreamCallDTO `json:"upstreams,omitempty"`
	UpstreamsJSON string            `json:"-"`
}

// DTO برای db_query
type DBQueryDTO struct {
	Connection string        `json:"connection"`
	SQL        string        `json:"sql"`
	Bindings   []interface{} `json:"bindings"`
	TimeMs     float64       `json:"time_ms"`
}

// DTO برای callهای HTTP خارجی (upstream)
type UpstreamCallDTO struct {
	URL             string         `json:"url"`
	Method          string         `json:"method"`
	Status          int            `json:"status"`
	DurationMs      int            `json:"duration_ms"`
	RequestHeaders  map[string]any `json:"request_headers,omitempty"`
	RequestBody     any            `json:"request_body,omitempty"`
	ResponseHeaders map[string]any `json:"response_headers,omitempty"`
	ResponseBody    any            `json:"response_body,omitempty"`
}

func NewHandler() *Handler {
	lokiURL := os.Getenv("LOKI_URL")
	if lokiURL == "" {
		panic("LOKI_URL is required")
	}

	// timezone ثابت برای کل UI (پیش‌فرض: Asia/Tehran)
	tz := os.Getenv("UI_TIMEZONE")
	if tz == "" {
		tz = "Asia/Tehran"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic("invalid UI_TIMEZONE: " + err.Error())
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"prettyJSON": func(v any) string {
			return prettyJSON(v, "")
		},
	}

	tmpl := template.Must(
		template.New("").
			Funcs(funcMap).
			ParseFiles(
				"templates/logs.tmpl",
				"templates/log_detail.tmpl",
				"templates/docs.tmpl",
			),
	)

	return &Handler{
		lokiClient: loki.NewClient(lokiURL),
		tmpl:       tmpl,
		loc:        loc,
	}
}

func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.lokiClient.Ready(ctx); err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if r.Method == http.MethodHead {
		return
	}

	_, _ = w.Write([]byte("ready\n"))
}

// ============================================================================
// /logs – صفحه‌ی HTML
// ============================================================================

func (h *Handler) LogsPage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	service := q.Get("service")
	level := q.Get("level")
	text := q.Get("text")
	reqID := q.Get("request_id")

	rangeStr := q.Get("range")
	if rangeStr == "" {
		rangeStr = "1h"
	}
	dur, err := time.ParseDuration(rangeStr)
	if err != nil {
		http.Error(w, "invalid range", http.StatusBadRequest)
		return
	}

	limitStr := q.Get("limit")
	if limitStr == "" {
		limitStr = "200"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 2000 {
		limit = 200
	}

	end := time.Now()
	start := end.Add(-dur)

	query := buildLogQL(service, level, text, reqID)

	entries, err := h.lokiClient.QueryRange(r.Context(), query, start, end, limit)
	if err != nil {
		http.Error(w, "loki error: "+err.Error(), http.StatusBadGateway)
		return
	}

	views := make([]LogView, 0, len(entries))

	for _, e := range entries {
		tt := e.Timestamp.In(h.loc)

		v := LogView{
			Time:    tt.Format("2006-01-02 15:04:05"),
			RawJSON: e.Raw,
		}

		if e.Parsed != nil && e.Parsed.Context != nil {
			ctx := e.Parsed.Context
			httpCtx := asMap(ctx["http"])
			v.Level = e.Parsed.LevelName
			v.Message = e.Parsed.Message

			if s, ok := ctx["service"].(string); ok {
				v.Service = s
			}
			if rt, ok := firstString(ctx, httpCtx, "route"); ok {
				v.Route = rt
			} else if p, ok := firstString(ctx, httpCtx, "path"); ok {
				v.Route = p
			}
			if m, ok := firstString(ctx, httpCtx, "method"); ok {
				v.Method = m
			}
			if id, ok := ctx["request_id"].(string); ok {
				v.RequestID = id
			}
			if lt, ok := ctx["log_type"].(string); ok {
				v.LogType = lt
			}
			if sc, ok := firstNumber(ctx, httpCtx, "status_code"); ok {
				v.Status = int(sc)
			}
			if d, ok := firstNumber(ctx, httpCtx, "duration_ms"); ok {
				v.DurationMs = int(d)
			}
			if ec, ok := ctx["error_code"].(string); ok && ec != "" {
				v.ErrorCode = ec
			}

			// safe_message مستقیم از لاگ، اگر نبود می‌افتیم روی response_message
			if sm, ok := ctx["safe_message"].(string); ok && sm != "" {
				v.SafeMessage = sm
			}
			if v.SafeMessage == "" {
				if rm, ok := ctx["response_message"].(string); ok && rm != "" {
					v.SafeMessage = rm
				}
			}
			// اگر error_code نداریم، از response_error_code بردار
			if v.ErrorCode == "" {
				if rec, ok := ctx["response_error_code"].(string); ok && rec != "" {
					v.ErrorCode = rec
				}
			}
		}

		if v.Service == "" {
			if s, ok := e.Labels["service_name"]; ok && s != "" {
				v.Service = s
			}
		}

		views = append(views, v)
	}

	data := LogsPageData{
		Logs:      views,
		Total:     len(views),
		Range:     rangeStr,
		Limit:     limit,
		Service:   service,
		Level:     level,
		Text:      text,
		RequestID: reqID,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "logs", data); err != nil {
		http.Error(w, "render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// ============================================================================
// /api/logs – JSON برای UI
// ============================================================================

func (h *Handler) LogsAPI(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// فیلترهای پایه برای LogQL
	service := q.Get("service")
	level := q.Get("level")
	text := q.Get("text")
	reqID := q.Get("request_id")

	// فیلترهای پیشرفته سمت Go
	routeFilter := q.Get("route")
	methodFilter := strings.ToUpper(q.Get("method"))
	logTypeFilter := q.Get("log_type")

	statusFilterStr := q.Get("status")
	minDurStr := q.Get("min_duration")
	maxDurStr := q.Get("max_duration")

	var statusFilter int
	var hasStatusFilter bool
	if statusFilterStr != "" {
		if v, err := strconv.Atoi(statusFilterStr); err == nil {
			statusFilter = v
			hasStatusFilter = true
		}
	}

	var minDur, maxDur int
	var hasMin, hasMax bool
	if minDurStr != "" {
		if v, err := strconv.Atoi(minDurStr); err == nil {
			minDur = v
			hasMin = true
		}
	}
	if maxDurStr != "" {
		if v, err := strconv.Atoi(maxDurStr); err == nil {
			maxDur = v
			hasMax = true
		}
	}

	rangeStr := q.Get("range")
	if rangeStr == "" {
		rangeStr = "1h"
	}
	dur, err := time.ParseDuration(rangeStr)
	if err != nil {
		http.Error(w, "invalid range", http.StatusBadRequest)
		return
	}

	limitStr := q.Get("limit")
	if limitStr == "" {
		limitStr = "500"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 2000 {
		limit = 500
	}

	end := time.Now()
	start := end.Add(-dur)

	// فقط لاگ‌های http_request
	query := buildLogQL(service, level, text, reqID) + ` |= "log_type\":\"http_request\""`

	entries, err := h.lokiClient.QueryRange(r.Context(), query, start, end, limit)
	if err != nil {
		http.Error(w, "loki error: "+err.Error(), http.StatusBadGateway)
		return
	}

	logs := make([]LogDTO, 0, len(entries))
	lowerText := strings.ToLower(text)

	for _, e := range entries {
		tt := e.Timestamp.In(h.loc)

		dto := LogDTO{
			Timestamp: e.Timestamp.UnixNano(),           // ثابت و مستقل از TZ
			Time:      tt.Format("2006-01-02 15:04:05"), // نمایش Tehran
			Context:   map[string]any{},
			Raw:       e.Raw,
			Labels:    e.Labels,
		}

		var ctx map[string]any

		if e.Parsed != nil {
			dto.Level = e.Parsed.LevelName
			dto.Message = e.Parsed.Message

			if e.Parsed.Context != nil {
				ctx = e.Parsed.Context
				httpCtx := asMap(ctx["http"])

				// فیلدهای عمومی
				if s, ok := ctx["service"].(string); ok {
					dto.Service = s
				}
				if rt, ok := firstString(ctx, httpCtx, "route"); ok {
					dto.Route = rt
				} else if p, ok := firstString(ctx, httpCtx, "path"); ok {
					dto.Route = p
				}
				if m, ok := firstString(ctx, httpCtx, "method"); ok {
					dto.Method = m
				}
				if id, ok := ctx["request_id"].(string); ok {
					dto.RequestID = id
				}
				if lt, ok := ctx["log_type"].(string); ok {
					dto.LogType = lt
				}
				if sc, ok := firstNumber(ctx, httpCtx, "status_code"); ok {
					dto.Status = int(sc)
				}
				if d, ok := firstNumber(ctx, httpCtx, "duration_ms"); ok {
					dto.DurationMs = int(d)
				}
				if ec, ok := ctx["error_code"].(string); ok && ec != "" {
					dto.ErrorCode = ec
				}
				if sm, ok := ctx["safe_message"].(string); ok && sm != "" {
					dto.SafeMessage = sm
				}

				// فیلدهای ریسپانس وب‌سرویس
				if v, ok := ctx["response_success"].(bool); ok {
					dto.ResponseSuccess = &v
				}
				if v, ok := ctx["response_message"].(string); ok {
					dto.ResponseMessage = v
				}
				if v, ok := ctx["response_error_code"].(string); ok {
					dto.ResponseErrorCode = v
				}
				if rb, ok := ctx["response_body"].(map[string]any); ok {
					dto.ResponseBody = rb
				}

				// fallback برای SafeMessage / ErrorCode از response_*
				if dto.SafeMessage == "" && dto.ResponseMessage != "" {
					dto.SafeMessage = dto.ResponseMessage
				}
				if dto.ErrorCode == "" && dto.ResponseErrorCode != "" {
					dto.ErrorCode = dto.ResponseErrorCode
				}

				for k, v := range ctx {
					dto.Context[k] = v
				}
			}
		}

		if dto.Service == "" {
			if s, ok := e.Labels["service_name"]; ok && s != "" {
				dto.Service = s
			}
		}

		// ---------------- فیلترهای پیشرفته سمت Go ----------------
		if routeFilter != "" && !strings.Contains(dto.Route, routeFilter) {
			continue
		}
		if methodFilter != "" && strings.ToUpper(dto.Method) != methodFilter {
			continue
		}
		if logTypeFilter != "" && dto.LogType != logTypeFilter {
			continue
		}
		if hasStatusFilter && dto.Status != statusFilter {
			continue
		}
		if hasMin && dto.DurationMs < minDur {
			continue
		}
		if hasMax && dto.DurationMs > maxDur {
			continue
		}

		// سرچ متنی شامل user و auth
		if text != "" {
			hayParts := []string{
				dto.Message,
				dto.Route,
				dto.Service,
				dto.Method,
				dto.ErrorCode,
				dto.SafeMessage,
				dto.ResponseMessage,
				dto.ResponseErrorCode,
			}

			if ctx != nil {
				// user
				if u, ok := ctx["user"].(map[string]any); ok {
					if v, ok := u["email"].(string); ok {
						hayParts = append(hayParts, v)
					}
					if v, ok := u["name"].(string); ok {
						hayParts = append(hayParts, v)
					}
				}
				// auth
				if a, ok := ctx["auth"].(map[string]any); ok {
					if v, ok := a["session_id"].(string); ok {
						hayParts = append(hayParts, v)
					}
					if v, ok := a["token_hash"].(string); ok {
						hayParts = append(hayParts, v)
					}
					if v, ok := a["token_id"].(float64); ok {
						hayParts = append(hayParts, fmt.Sprintf("%.0f", v))
					}
					if au, ok := a["user"].(map[string]any); ok {
						if v, ok := au["email"].(string); ok {
							hayParts = append(hayParts, v)
						}
						if v, ok := au["name"].(string); ok {
							hayParts = append(hayParts, v)
						}
					}
				}
			}

			hay := strings.ToLower(strings.Join(hayParts, " "))
			if !strings.Contains(hay, lowerText) {
				continue
			}
		}

		logs = append(logs, dto)
	}

	resp := map[string]any{"logs": logs}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(resp)
}

// ============================================================================
// /logs/detail – نمایش کامل یک لاگ بر اساس request_id
// ============================================================================

func (h *Handler) LogDetailPage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	reqID := q.Get("request_id")
	if reqID == "" {
		http.Error(w, "request_id is required", http.StatusBadRequest)
		return
	}

	rangeStr := q.Get("range")
	if rangeStr == "" {
		rangeStr = "24h"
	}
	dur, err := time.ParseDuration(rangeStr)
	if err != nil || dur <= 0 {
		dur = 24 * time.Hour
	}

	end := time.Now()
	start := end.Add(-dur)

	query := buildLogQL("", "", "", reqID)

	entries, err := h.lokiClient.QueryRange(r.Context(), query, start, end, 200)
	if err != nil {
		http.Error(w, "query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(entries) == 0 {
		http.Error(w, "log entry not found", http.StatusNotFound)
		return
	}

	var reqEntry *loki.LogEntry
	var excEntry *loki.LogEntry
	var dbQueries []DBQueryDTO
	var upstreams []UpstreamCallDTO

	for i := range entries {
		e := &entries[i]

		logType := ""
		if e.Parsed != nil && e.Parsed.Context != nil {
			if v, ok := e.Parsed.Context["log_type"].(string); ok {
				logType = v
			}
		}

		switch logType {
		case "http_request":
			if reqEntry == nil {
				reqEntry = e
			}
		case "exception":
			if excEntry == nil {
				excEntry = e
			}
		case "db_query":
			if e.Parsed != nil && e.Parsed.Context != nil {
				ctx := e.Parsed.Context
				dq := DBQueryDTO{}
				if v, ok := ctx["connection"].(string); ok {
					dq.Connection = v
				}
				if v, ok := ctx["sql"].(string); ok {
					dq.SQL = v
				}
				if v, ok := ctx["time_ms"].(float64); ok {
					dq.TimeMs = v
				}
				if rawBindings, ok := ctx["bindings"]; ok {
					if arr, ok := rawBindings.([]any); ok {
						dq.Bindings = arr
					}
				}
				dbQueries = append(dbQueries, dq)
			}
		case "http_upstream":
			if e.Parsed != nil && e.Parsed.Context != nil {
				ctx := e.Parsed.Context
				us := UpstreamCallDTO{}

				if raw, ok := ctx["upstream"]; ok {
					if m, ok := raw.(map[string]any); ok {
						if v, ok := m["url"].(string); ok {
							us.URL = v
						}
						if v, ok := m["method"].(string); ok {
							us.Method = v
						}
						if v, ok := m["duration_ms"].(float64); ok {
							us.DurationMs = int(v)
						}
						if v, ok := m["status"].(float64); ok {
							us.Status = int(v)
						}
						if v, ok := m["request_headers"].(map[string]any); ok {
							us.RequestHeaders = v
						}
						us.RequestBody = m["request_body"]
						if v, ok := m["response_headers"].(map[string]any); ok {
							us.ResponseHeaders = v
						}
						us.ResponseBody = m["response_body"]
					}
				}
				upstreams = append(upstreams, us)
			}
		}
	}

	if reqEntry == nil {
		reqEntry = &entries[0]
	}
	e := reqEntry
	tt := e.Timestamp.In(h.loc)

	dto := LogDTO{
		Timestamp: e.Timestamp.UnixNano(),
		Time:      tt.Format("2006-01-02 15:04:05.000"),
		Context:   map[string]any{},
		Raw:       e.Raw,
		Labels:    e.Labels,
	}

	if e.Parsed != nil {
		dto.Level = e.Parsed.LevelName
		dto.Message = e.Parsed.Message

		if ctx := e.Parsed.Context; ctx != nil {
			httpCtx := asMap(ctx["http"])
			for k, v := range ctx {
				dto.Context[k] = v
			}

			if v, ok := ctx["service"].(string); ok {
				dto.Service = v
			}
			if v, ok := firstString(ctx, httpCtx, "route"); ok {
				dto.Route = v
			} else if v, ok := firstString(ctx, httpCtx, "path"); ok {
				dto.Route = v
			}
			if v, ok := firstString(ctx, httpCtx, "method"); ok {
				dto.Method = v
			}
			if v, ok := ctx["request_id"].(string); ok {
				dto.RequestID = v
			}
			if v, ok := ctx["log_type"].(string); ok {
				dto.LogType = v
			}
			if v, ok := firstNumber(ctx, httpCtx, "status_code"); ok {
				dto.Status = int(v)
			}
			if v, ok := firstNumber(ctx, httpCtx, "duration_ms"); ok {
				dto.DurationMs = int(v)
			}
			if ec, ok := ctx["error_code"].(string); ok && ec != "" {
				dto.ErrorCode = ec
			}
			if sm, ok := ctx["safe_message"].(string); ok && sm != "" {
				dto.SafeMessage = sm
			}

			if v, ok := ctx["response_success"].(bool); ok {
				dto.ResponseSuccess = &v
			}
			if v, ok := ctx["response_message"].(string); ok {
				dto.ResponseMessage = v
			}
			if v, ok := ctx["response_error_code"].(string); ok {
				dto.ResponseErrorCode = v
			}
			if rb, ok := ctx["response_body"].(map[string]any); ok {
				dto.ResponseBody = rb
			}

			if dto.SafeMessage == "" && dto.ResponseMessage != "" {
				dto.SafeMessage = dto.ResponseMessage
			}
			if dto.ErrorCode == "" && dto.ResponseErrorCode != "" {
				dto.ErrorCode = dto.ResponseErrorCode
			}
		}
	}

	if dto.Service == "" {
		if s, ok := e.Labels["service_name"]; ok && s != "" {
			dto.Service = s
		}
	}

	if excEntry != nil && excEntry.Parsed != nil && excEntry.Parsed.Context != nil {
		ctx := excEntry.Parsed.Context
		if ec, ok := ctx["error_code"].(string); ok && ec != "" {
			dto.ErrorCode = ec
			dto.Context["error_code"] = ec
		}
		if sm, ok := ctx["safe_message"].(string); ok && sm != "" {
			dto.SafeMessage = sm
			dto.Context["safe_message"] = sm
		}
	}

	dto.DBQueries = dbQueries
	dto.DBQueriesJSON = prettyJSON(dbQueries, "[]")

	dto.Upstreams = upstreams
	dto.UpstreamsJSON = prettyJSON(upstreams, "[]")

	if dto.Raw != "" {
		var tmp any
		if err := json.Unmarshal([]byte(dto.Raw), &tmp); err == nil {
			dto.RawPretty = prettyJSON(tmp, "{}")
		} else {
			dto.RawPretty = dto.Raw
		}
	} else {
		dto.RawPretty = "{}"
	}

	ctxMap := dto.Context
	httpCtx := asMap(ctxMap["http"])
	reqMap := asMap(ctxMap["request"])

	reqObj := map[string]any{
		"method":     firstValue(ctxMap, httpCtx, "method"),
		"path":       firstValue(ctxMap, httpCtx, "path"),
		"route":      firstValue(ctxMap, httpCtx, "route"),
		"ip":         firstValue(ctxMap, httpCtx, "ip", "client_ip"),
		"user_agent": firstValue(ctxMap, httpCtx, "user_agent"),
		"headers":    reqMap["headers"],
		"payload":    reqMap["payload"],
		"query":      reqMap["query"],
	}
	dto.RequestJSON = prettyJSON(reqObj, "{}")

	respObj := map[string]any{
		"status_code":     firstValue(ctxMap, httpCtx, "status_code"),
		"response_status": ctxMap["response_status"],
		"success":         ctxMap["response_success"],
		"message":         ctxMap["response_message"],
		"error_code":      ctxMap["response_error_code"],
		"safe_message":    ctxMap["safe_message"],
		"headers":         ctxMap["response_headers"],
		"body":            ctxMap["response_body"],
	}
	dto.ResponseJSON = prettyJSON(respObj, "{}")

	var queryPart any
	if qv, ok := reqMap["query"]; ok {
		queryPart = qv
	} else if qv, ok := ctxMap["query"]; ok {
		queryPart = qv
	}
	dto.QueryJSON = prettyJSON(queryPart, "[]")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "log_detail", dto); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// ============================================================================
// LogQL builder (با escape حداقلی برای جلوگیری از شکستن query)
// ============================================================================

func buildLogQL(service, level, text, requestID string) string {
	selector := `{job="laravel"}`
	if service != "" {
		selector = `{job="laravel",service_name="` + escapeLogQL(service) + `"}`
	}

	filter := ""
	if requestID != "" {
		filter += ` |= "request_id\":\"` + escapeLogQL(requestID) + `\""`
	}
	if level != "" {
		filter += ` |= "\"level_name\":\"` + escapeLogQL(level) + `\""`
	}
	if text != "" {
		filter += ` |= "` + escapeLogQL(text) + `"`
	}

	return selector + filter
}

func escapeLogQL(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return s
}

// ============================================================================
// JSON helpers
// ============================================================================

func firstString(primary, nested map[string]any, keys ...string) (string, bool) {
	for _, key := range keys {
		if v, ok := primary[key].(string); ok && v != "" {
			return v, true
		}
		if v, ok := nested[key].(string); ok && v != "" {
			return v, true
		}
	}
	return "", false
}

func firstNumber(primary, nested map[string]any, keys ...string) (float64, bool) {
	for _, key := range keys {
		if v, ok := numberValue(primary[key]); ok {
			return v, true
		}
		if v, ok := numberValue(nested[key]); ok {
			return v, true
		}
	}
	return 0, false
}

func firstValue(primary, nested map[string]any, keys ...string) any {
	for _, key := range keys {
		if v, ok := primary[key]; ok && v != nil {
			return v
		}
		if v, ok := nested[key]; ok && v != nil {
			return v
		}
	}
	return nil
}

func numberValue(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

func prettyJSON(v any, empty string) string {
	if v == nil {
		if empty != "" {
			return empty
		}
		return ""
	}

	if s, ok := v.(string); ok {
		s = strings.TrimSpace(s)
		if s == "" {
			if empty != "" {
				return empty
			}
			return ""
		}
		var tmp any
		if err := json.Unmarshal([]byte(s), &tmp); err == nil {
			b, err := json.MarshalIndent(tmp, "", "  ")
			if err == nil {
				return string(b)
			}
		}
		return s
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		if empty != "" {
			return empty
		}
		return ""
	}
	return string(b)
}

func asMap(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

func (h *Handler) RequestsPage(w http.ResponseWriter, r *http.Request) {
	h.LogsPage(w, r)
}
func (h *Handler) DocsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lang := normalizedDocLang(r.URL.Query().Get("lang"))

	files, err := filepath.Glob(filepath.Join("docs", "progress", lang, "SECTION-*.md"))
	if err != nil {
		http.Error(w, "docs lookup error", http.StatusInternalServerError)
		return
	}

	if len(files) == 0 && lang == "fa" {
		files, err = filepath.Glob(filepath.Join("docs", "progress", "SECTION-*.md"))
		if err != nil {
			http.Error(w, "docs lookup error", http.StatusInternalServerError)
			return
		}
	}

	sort.Strings(files)

	if len(files) == 0 {
		http.Error(w, "no docs found for selected language", http.StatusNotFound)
		return
	}

	page := 1
	if raw := r.URL.Query().Get("page"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > len(files) {
			http.Error(w, "invalid docs page", http.StatusBadRequest)
			return
		}
		page = n
	}

	currentFile := files[page-1]

	b, err := os.ReadFile(currentFile)
	if err != nil {
		http.Error(w, "docs read error", http.StatusInternalServerError)
		return
	}

	title := docTitle(string(b))
	if title == "" {
		title = currentFile
	}

	oppositeLang := oppositeDocLang(lang)

	data := DocsPageData{
		Title:             title,
		BodyHTML:          renderDocMarkdown(string(b)),
		Page:              page,
		Total:             len(files),
		PrevPage:          page - 1,
		NextPage:          page + 1,
		HasPrev:           page > 1,
		HasNext:           page < len(files),
		CurrentFile:       currentFile,
		Lang:              lang,
		LangLabel:         docLangLabel(lang),
		TextDir:           docTextDir(lang),
		OppositeLang:      oppositeLang,
		OppositeLangLabel: docLangLabel(oppositeLang),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := h.tmpl.ExecuteTemplate(w, "docs", data); err != nil {
		http.Error(w, "render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func docTitle(md string) string {
	for _, line := range strings.Split(md, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}
func normalizedDocLang(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "en":
		return "en"
	default:
		return "fa"
	}
}

func oppositeDocLang(lang string) string {
	if lang == "en" {
		return "fa"
	}

	return "en"
}

func docLangLabel(lang string) string {
	if lang == "en" {
		return "English"
	}

	return "فارسی"
}

func docTextDir(lang string) string {
	if lang == "en" {
		return "ltr"
	}

	return "rtl"
}
func renderDocMarkdown(md string) template.HTML {
	var out bytes.Buffer
	lines := strings.Split(md, "\n")

	inCode := false
	inList := false
	var paragraph []string

	flushParagraph := func() {
		if len(paragraph) == 0 {
			return
		}

		text := strings.Join(paragraph, " ")
		out.WriteString("<p>")
		out.WriteString(renderInlineMarkdown(text))
		out.WriteString("</p>\n")
		paragraph = nil
	}

	flushList := func() {
		if !inList {
			return
		}

		out.WriteString("</ul>\n")
		inList = false
	}

	for _, raw := range lines {
		line := strings.TrimRight(raw, "\r")
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			flushParagraph()
			flushList()

			if inCode {
				out.WriteString("</code></pre>\n")
				inCode = false
				continue
			}

			out.WriteString("<pre><code>")
			inCode = true
			continue
		}

		if inCode {
			out.WriteString(stdhtml.EscapeString(line))
			out.WriteByte('\n')
			continue
		}

		if trimmed == "" {
			flushParagraph()
			flushList()
			continue
		}

		switch {
		case strings.HasPrefix(trimmed, "### "):
			flushParagraph()
			flushList()
			out.WriteString("<h3>")
			out.WriteString(renderInlineMarkdown(strings.TrimSpace(strings.TrimPrefix(trimmed, "### "))))
			out.WriteString("</h3>\n")

		case strings.HasPrefix(trimmed, "## "):
			flushParagraph()
			flushList()
			out.WriteString("<h2>")
			out.WriteString(renderInlineMarkdown(strings.TrimSpace(strings.TrimPrefix(trimmed, "## "))))
			out.WriteString("</h2>\n")

		case strings.HasPrefix(trimmed, "# "):
			flushParagraph()
			flushList()
			out.WriteString("<h1>")
			out.WriteString(renderInlineMarkdown(strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))))
			out.WriteString("</h1>\n")

		case strings.HasPrefix(trimmed, "- "):
			flushParagraph()
			if !inList {
				out.WriteString("<ul>\n")
				inList = true
			}
			out.WriteString("<li>")
			out.WriteString(renderInlineMarkdown(strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))))
			out.WriteString("</li>\n")

		default:
			paragraph = append(paragraph, trimmed)
		}
	}

	flushParagraph()
	flushList()

	if inCode {
		out.WriteString("</code></pre>\n")
	}

	return template.HTML(out.String())
}

func renderInlineMarkdown(s string) string {
	escaped := stdhtml.EscapeString(s)

	var out strings.Builder
	inCode := false
	var code strings.Builder

	for _, r := range escaped {
		if r == '`' {
			if inCode {
				out.WriteString("<code>")
				out.WriteString(code.String())
				out.WriteString("</code>")
				code.Reset()
				inCode = false
			} else {
				inCode = true
			}
			continue
		}

		if inCode {
			code.WriteRune(r)
			continue
		}

		out.WriteRune(r)
	}

	if inCode {
		out.WriteRune('`')
		out.WriteString(code.String())
	}

	return out.String()
}

func (h *Handler) RequestsAPI(w http.ResponseWriter, r *http.Request) {
	h.LogsAPI(w, r)
}
