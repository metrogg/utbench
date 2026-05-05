package web

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"sort"
	"strings"
)

// knownAPIKeys 定义系统已知的 API Key 及其用途描述。
// 用户可以在 Web UI 中查看状态并填入值，保存后写入 .env 文件。
var knownAPIKeys = []apiKeyDef{
	// Agent 平台认证
	{
		Key:         "CODEBUDDY_API_KEY",
		Label:       "CodeBuddy API Key",
		Category:    "agent",
		Description: "CodeBuddy 平台认证密钥，使用 CodeBuddy Agent 时必须。从 copilot.tencent.com/profile/ 获取。",
		DocURL:      "https://www.codebuddy.ai/docs/zh/cli/authentication",
		Required:    false,
	},
	{
		Key:         "CODEBUDDY_INTERNET_ENVIRONMENT",
		Label:       "CodeBuddy 网络环境",
		Category:    "agent",
		Description: "中国版设为 internal，海外版留空或 public，iOA 版设为 ioa。",
		Required:    false,
	},
	{
		Key:         "ANTHROPIC_API_KEY",
		Label:       "Anthropic API Key",
		Category:    "agent",
		Description: "Claude Code 直连 Anthropic API 时使用。",
		DocURL:      "https://docs.anthropic.com/en/docs/claude-code/quickstart",
		Required:    false,
	},
	{
		Key:         "ANTHROPIC_AUTH_TOKEN",
		Label:       "Anthropic Auth Token",
		Category:    "agent",
		Description: "Claude Code 通过 Gateway 或代理访问时可使用的鉴权令牌。",
		DocURL:      "https://code.claude.com/docs/en/llm-gateway",
		Required:    false,
	},
	{
		Key:         "ANTHROPIC_BASE_URL",
		Label:       "Anthropic Base URL",
		Category:    "agent",
		Description: "Claude Code 第三方 Gateway / Proxy 入口地址。",
		DocURL:      "https://code.claude.com/docs/en/third-party-integrations",
		Required:    false,
	},
	// 模型服务商
	{
		Key:         "DEEPSEEK_API_KEY",
		Label:       "DeepSeek API Key",
		Category:    "model",
		Description: "DeepSeek 模型 API 密钥。",
		Required:    false,
	},
	{
		Key:         "DASHSCOPE_API_KEY",
		Label:       "Qwen (Dashscope) API Key",
		Category:    "model",
		Description: "通义千问 Dashscope API 密钥。",
		Required:    false,
	},
	{
		Key:         "MINIMAX_API_KEY",
		Label:       "MiniMax API Key",
		Category:    "model",
		Description: "MiniMax 模型 API 密钥。",
		Required:    false,
	},
	{
		Key:         "MIMO_V25_API_KEY",
		Label:       "MiMo V2.5 API Key",
		Category:    "model",
		Description: "小米 MiMo V2.5 模型 API 密钥。",
		Required:    false,
	},
	{
		Key:         "MIMO_V25_PRO_API_KEY",
		Label:       "MiMo V2.5 Pro API Key",
		Category:    "model",
		Description: "小米 MiMo V2.5 Pro 模型 API 密钥。",
		Required:    false,
	},
	{
		Key:         "VOLCENGINE_API_KEY",
		Label:       "火山引擎 API Key",
		Category:    "model",
		Description: "火山引擎 (doubao) API 密钥。",
		Required:    false,
	},
	{
		Key:         "ARK_API_KEY",
		Label:       "Ark API Key",
		Category:    "model",
		Description: "火山方舟 API 密钥 (doubao-seed/glm/deepseek 系列)。",
		Required:    false,
	},
	{
		Key:         "BIGMODEL_API_KEY",
		Label:       "智谱 API Key",
		Category:    "model",
		Description: "智谱 BigModel API 密钥。",
		Required:    false,
	},
}

type apiKeyDef struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Category    string `json:"category"`
	Description string `json:"description"`
	DocURL      string `json:"doc_url,omitempty"`
	Required    bool   `json:"required"`
}

type apiKeyStatus struct {
	apiKeyDef
	ValueSet bool   `json:"value_set"`
	Preview  string `json:"preview,omitempty"` // 脱敏预览，如 "ck_f1...sU"
}

// handleAPIKeys 管理已知 API Key 的状态和值。
//
//	GET  /api/settings/api-keys  → 返回所有已知 key 的状态
//	POST /api/settings/api-keys  → 保存一个或多个 key 到 .env + 进程环境变量
func (s *Server) handleAPIKeys(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleAPIKeysGet(w, r)
	case http.MethodPost:
		s.handleAPIKeysPost(w, r)
	default:
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) handleAPIKeysGet(w http.ResponseWriter, _ *http.Request) {
	statuses := make([]apiKeyStatus, 0, len(knownAPIKeys))
	for _, def := range knownAPIKeys {
		val := s.lookupAPIKey(def.Key)
		st := apiKeyStatus{
			apiKeyDef: def,
			ValueSet:  hasUsableAPIKey(val),
			Preview:   maskKey(val),
		}
		statuses = append(statuses, st)
	}
	writeJSON(w, http.StatusOK, map[string]any{"keys": statuses})
}

func (s *Server) lookupAPIKey(key string) string {
	val := os.Getenv(key)
	if hasUsableAPIKey(val) {
		return val
	}
	if s == nil || s.dockerCfg.EnvFile == "" {
		return val
	}
	if v, ok := readEnvFileMap(s.dockerCfg.EnvFile)[key]; ok && hasUsableAPIKey(v) {
		return v
	}
	return val
}

type apiKeysPostRequest struct {
	Keys map[string]string `json:"keys"` // key_name → value
}

func (s *Server) handleAPIKeysPost(w http.ResponseWriter, r *http.Request) {
	var req apiKeysPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.Keys) == 0 {
		errJSON(w, http.StatusBadRequest, "no keys provided")
		return
	}

	// 校验 key 名称是否在已知列表中
	knownSet := make(map[string]bool, len(knownAPIKeys))
	for _, def := range knownAPIKeys {
		knownSet[def.Key] = true
	}

	var saved []string
	for key, val := range req.Keys {
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		if key == "" || val == "" {
			continue
		}
		if !knownSet[key] {
			continue
		}
		// 写入进程环境变量
		_ = os.Setenv(key, val)
		// 同步到 .env 文件
		if s.dockerCfg.EnvFile != "" {
			_ = upsertEnvFileValue(s.dockerCfg.EnvFile, key, val)
		}
		saved = append(saved, key)
	}
	sort.Strings(saved)
	writeJSON(w, http.StatusOK, map[string]any{"saved": saved})
}

// maskKey 对 API Key 做脱敏处理。
// 保留前 4 位和后 2 位，中间用 ... 代替。
func maskKey(val string) string {
	val = strings.TrimSpace(val)
	if val == "" {
		return ""
	}
	runes := []rune(val)
	if len(runes) <= 8 {
		return string(runes[:min(3, len(runes))]) + "..."
	}
	return string(runes[:4]) + "..." + string(runes[len(runes)-2:])
}

// readEnvFileMap 读取 .env 文件为 key→value 映射，忽略注释和空行。
func readEnvFileMap(path string) map[string]string {
	m := make(map[string]string)
	if path == "" {
		return m
	}
	f, err := os.Open(path)
	if err != nil {
		return m
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if i := strings.IndexByte(line, '='); i > 0 {
			key := strings.TrimSpace(line[:i])
			val := strings.TrimSpace(line[i+1:])
			// 去掉可能的引号
			if len(val) >= 2 && (val[0] == '"' && val[len(val)-1] == '"' || val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
			m[key] = val
		}
	}
	return m
}
