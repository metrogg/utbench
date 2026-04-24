package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// modelsMu 保护对 configs/models.yaml 的并发读写。
var modelsMu sync.Mutex

// modelEntry 是 API 层返回/接受的模型条目。
// 注意：api_key 字段仅写入 .env（或运行时环境变量）约定的变量名，
// 真实密钥值不会落盘到 YAML；YAML 只保留 api_key_env 字段名。
type modelEntry struct {
	Name        string            `json:"name"`
	Enabled     bool              `json:"enabled"`
	Provider    string            `json:"provider"`
	ModelID     string            `json:"model_id"`
	APIEndpoint string            `json:"api_endpoint"`
	APIKeyEnv   string            `json:"api_key_env"`
	APIKey      string            `json:"api_key,omitempty"` // 仅 POST/PUT 时接收，用于写入 .env
	APIKeySet   bool              `json:"api_key_set"`       // GET 时标示对应 env 变量是否已设置
	Parameters  map[string]any    `json:"parameters,omitempty"`
	Extra       map[string]any    `json:"extra,omitempty"` // 保留未识别字段
	_           map[string]string `json:"-"`
}

// handleModels 路由 GET/POST /api/models。
func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listModels(w, r)
	case http.MethodPost:
		s.createModel(w, r)
	default:
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleModelsSub 路由 PUT/DELETE /api/models/{name}。
func (s *Server) handleModelsSub(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/models/")
	action := ""
	if before, after, ok := strings.Cut(name, "/"); ok {
		name = before
		action = after
	}
	if name == "" || strings.Contains(name, "/") {
		errJSON(w, http.StatusBadRequest, "invalid model name")
		return
	}
	if action == "test" {
		s.handleTestModel(w, r, name)
		return
	}
	if action != "" {
		errJSON(w, http.StatusNotFound, "unknown model action")
		return
	}
	switch r.Method {
	case http.MethodPut:
		s.updateModel(w, r, name)
	case http.MethodDelete:
		s.deleteModel(w, r, name)
	default:
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

type modelTestResult struct {
	Name      string `json:"name"`
	OK        bool   `json:"ok"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	LatencyMS int64  `json:"latency_ms"`
	HTTPCode  int    `json:"http_code,omitempty"`
	CheckedAt string `json:"checked_at"`
}

func (s *Server) handleTestModel(w http.ResponseWriter, r *http.Request, name string) {
	if r.Method != http.MethodPost {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	modelsMu.Lock()
	root, err := readModelsRoot(s.configPath)
	modelsMu.Unlock()
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	entry, ok := findModelEntry(root, name)
	if !ok {
		errJSON(w, http.StatusNotFound, "model not found: "+name)
		return
	}
	writeJSON(w, http.StatusOK, testModelConnection(r.Context(), entry))
}

func (s *Server) handleTestAllModels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	modelsMu.Lock()
	root, err := readModelsRoot(s.configPath)
	modelsMu.Unlock()
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	entries := extractEntries(root)
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	results := make([]modelTestResult, 0, len(entries))
	for _, entry := range entries {
		results = append(results, testModelConnection(r.Context(), entry))
	}
	writeJSON(w, http.StatusOK, results)
}

func findModelEntry(root map[string]any, name string) (modelEntry, bool) {
	for _, entry := range extractEntries(root) {
		if entry.Name == name {
			return entry, true
		}
	}
	return modelEntry{}, false
}

func testModelConnection(parent context.Context, entry modelEntry) modelTestResult {
	result := modelTestResult{
		Name:      entry.Name,
		Status:    "failed",
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if strings.TrimSpace(entry.APIEndpoint) == "" {
		result.Message = "API endpoint is empty"
		return result
	}
	if strings.TrimSpace(entry.ModelID) == "" {
		result.Message = "model id is empty"
		return result
	}
	apiKey := strings.TrimSpace(os.Getenv(entry.APIKeyEnv))
	if apiKey == "" {
		result.Status = "missing_key"
		result.Message = "missing API key env var: " + entry.APIKeyEnv
		return result
	}

	ctx, cancel := context.WithTimeout(parent, 25*time.Second)
	defer cancel()

	body, err := json.Marshal(buildModelTestPayload(entry))
	if err != nil {
		result.Message = err.Error()
		return result
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resolveModelTestEndpoint(entry), bytes.NewReader(body))
	if err != nil {
		result.Message = err.Error()
		return result
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	started := time.Now()
	resp, err := (&http.Client{Timeout: 25 * time.Second}).Do(req)
	result.LatencyMS = time.Since(started).Milliseconds()
	if err != nil {
		result.Status = "network_error"
		result.Message = err.Error()
		return result
	}
	defer resp.Body.Close()
	result.HTTPCode = resp.StatusCode
	rawBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	rawText := string(rawBytes)
	if resp.StatusCode >= 400 {
		result.Status = "http_error"
		result.Message = fmt.Sprintf("http %d: %s", resp.StatusCode, trimModelTestText(rawText, 500))
		return result
	}
	if err := validateModelTestResponse(rawBytes, entry.Provider); err != nil {
		result.Status = "response_error"
		result.Message = err.Error()
		return result
	}
	result.OK = true
	result.Status = "ok"
	result.Message = "connection ok"
	return result
}

func resolveModelTestEndpoint(entry modelEntry) string {
	base := strings.TrimSuffix(strings.TrimSpace(entry.APIEndpoint), "/")
	if strings.EqualFold(entry.Provider, "dashscope") {
		if strings.Contains(base, "compatible-mode") {
			return base + "/chat/completions"
		}
		if strings.HasSuffix(base, "/api/v1") {
			return base + "/services/aigc/text-generation/generation"
		}
		return base + "/services/aigc/text-generation/generation"
	}
	return base + "/chat/completions"
}

func buildModelTestPayload(entry modelEntry) map[string]any {
	params := map[string]any{"max_tokens": 8, "temperature": 0}
	if strings.EqualFold(entry.Provider, "dashscope") && !strings.Contains(entry.APIEndpoint, "compatible-mode") {
		return map[string]any{
			"model":      entry.ModelID,
			"input":      map[string]any{"messages": []map[string]any{{"role": "user", "content": "ping"}}},
			"parameters": params,
		}
	}
	return map[string]any{
		"model": entry.ModelID,
		"messages": []map[string]any{
			{"role": "user", "content": "ping"},
		},
		"stream":      false,
		"max_tokens":  8,
		"temperature": 0,
	}
}

func validateModelTestResponse(raw []byte, provider string) error {
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return fmt.Errorf("invalid JSON response: %w", err)
	}
	if strings.EqualFold(provider, "dashscope") {
		if output, ok := payload["output"].(map[string]any); ok {
			if _, ok := output["text"].(string); ok {
				return nil
			}
			if choices, ok := output["choices"].([]any); ok && len(choices) > 0 {
				return nil
			}
		}
	}
	if choices, ok := payload["choices"].([]any); ok && len(choices) > 0 {
		return nil
	}
	return fmt.Errorf("response does not contain choices/output")
}

func trimModelTestText(v string, max int) string {
	v = strings.TrimSpace(v)
	if len(v) <= max {
		return v
	}
	return v[:max]
}

func (s *Server) listModels(w http.ResponseWriter, _ *http.Request) {
	modelsMu.Lock()
	defer modelsMu.Unlock()

	root, err := readModelsRoot(s.configPath)
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := extractEntries(root)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) createModel(w http.ResponseWriter, r *http.Request) {
	var in modelEntry
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		errJSON(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		errJSON(w, http.StatusBadRequest, "name is required")
		return
	}
	modelsMu.Lock()
	defer modelsMu.Unlock()
	root, err := readModelsRoot(s.configPath)
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	models := ensureModelsMap(root)
	if _, exists := models[in.Name]; exists {
		errJSON(w, http.StatusConflict, "model already exists: "+in.Name)
		return
	}
	upsertModelNode(models, in)
	if err := writeModelsRoot(s.configPath, root); err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := persistAPIKey(&in); err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"name": in.Name, "status": "created"})
}

func (s *Server) updateModel(w http.ResponseWriter, r *http.Request, name string) {
	var in modelEntry
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		errJSON(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	in.Name = name
	modelsMu.Lock()
	defer modelsMu.Unlock()
	root, err := readModelsRoot(s.configPath)
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	models := ensureModelsMap(root)
	if _, exists := models[name]; !exists {
		errJSON(w, http.StatusNotFound, "model not found: "+name)
		return
	}
	upsertModelNode(models, in)
	if err := writeModelsRoot(s.configPath, root); err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := persistAPIKey(&in); err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"name": name, "status": "updated"})
}

func (s *Server) deleteModel(w http.ResponseWriter, _ *http.Request, name string) {
	modelsMu.Lock()
	defer modelsMu.Unlock()
	root, err := readModelsRoot(s.configPath)
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	models := ensureModelsMap(root)
	if _, exists := models[name]; !exists {
		errJSON(w, http.StatusNotFound, "model not found: "+name)
		return
	}
	delete(models, name)
	if err := writeModelsRoot(s.configPath, root); err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"name": name, "status": "deleted"})
}

// ─── YAML read/write helpers ────────────────────────────────────────────────

// readModelsRoot 以 map[string]any 形式读取整个 YAML（保留 benchmark/languages 块）。
func readModelsRoot(path string) (map[string]any, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	root := map[string]any{}
	if len(raw) > 0 {
		if err := yaml.Unmarshal(raw, &root); err != nil {
			return nil, err
		}
	}
	return root, nil
}

func writeModelsRoot(path string, root map[string]any) error {
	out, err := yaml.Marshal(root)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}

func ensureModelsMap(root map[string]any) map[string]any {
	m, ok := root["models"].(map[string]any)
	if !ok {
		m = map[string]any{}
		root["models"] = m
	}
	return m
}

// extractEntries 把 YAML 的 models 块展平为前端用的 modelEntry 列表。
func extractEntries(root map[string]any) []modelEntry {
	models, _ := root["models"].(map[string]any)
	out := make([]modelEntry, 0, len(models))
	for name, v := range models {
		node, _ := v.(map[string]any)
		entry := modelEntry{Name: name}
		entry.Enabled, _ = node["enabled"].(bool)
		entry.Provider, _ = node["provider"].(string)
		cfg, _ := node["config"].(map[string]any)
		if cfg != nil {
			entry.ModelID, _ = cfg["model"].(string)
			entry.APIEndpoint, _ = cfg["api_endpoint"].(string)
			entry.APIKeyEnv, _ = cfg["api_key_env"].(string)
			if p, ok := cfg["parameters"].(map[string]any); ok {
				entry.Parameters = p
			}
		}
		if entry.APIKeyEnv != "" {
			entry.APIKeySet = os.Getenv(entry.APIKeyEnv) != ""
		}
		out = append(out, entry)
	}
	return out
}

// upsertModelNode 将前端提交的条目写回 YAML 的 models[name] 节点。
// 不覆盖未识别字段（parameters 以外的 config 子项、models[name] 的其它兄弟键）。
func upsertModelNode(models map[string]any, in modelEntry) {
	existing, _ := models[in.Name].(map[string]any)
	if existing == nil {
		existing = map[string]any{}
	}
	existing["enabled"] = in.Enabled
	existing["provider"] = in.Provider
	cfg, _ := existing["config"].(map[string]any)
	if cfg == nil {
		cfg = map[string]any{}
	}
	cfg["model"] = in.ModelID
	cfg["api_endpoint"] = in.APIEndpoint
	if in.APIKeyEnv == "" {
		in.APIKeyEnv = strings.ToUpper(strings.ReplaceAll(in.Name, "-", "_")) + "_API_KEY"
	}
	cfg["api_key_env"] = in.APIKeyEnv
	if in.Parameters != nil {
		cfg["parameters"] = in.Parameters
	} else if _, ok := cfg["parameters"]; !ok {
		cfg["parameters"] = map[string]any{"temperature": 0.7, "top_p": 0.9, "max_tokens": 4096}
	}
	existing["config"] = cfg
	models[in.Name] = existing
}

// persistAPIKey 若请求体带 api_key 字段，则写入进程环境变量，供后续运行使用。
// 注意：不写回磁盘 .env 以避免密钥永久落盘；如用户希望持久化，应自行修改 .env。
func persistAPIKey(in *modelEntry) error {
	if in.APIKey == "" || in.APIKeyEnv == "" {
		return nil
	}
	return os.Setenv(in.APIKeyEnv, in.APIKey)
}
