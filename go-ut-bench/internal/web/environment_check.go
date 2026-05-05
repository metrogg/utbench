package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type environmentCheckResponse struct {
	CheckedAt   string                  `json:"checked_at"`
	OS          string                  `json:"os"`
	ProjectRoot string                  `json:"project_root"`
	Summary     environmentCheckSummary `json:"summary"`
	Groups      []environmentCheckGroup `json:"groups"`
}

type environmentCheckSummary struct {
	Total       int `json:"total"`
	OK          int `json:"ok"`
	Warning     int `json:"warning"`
	Missing     int `json:"missing"`
	Unknown     int `json:"unknown"`
	Installable int `json:"installable"`
}

type environmentCheckGroup struct {
	ID    string                 `json:"id"`
	Label string                 `json:"label"`
	Items []environmentCheckItem `json:"items"`
}

type environmentCheckItem struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Category       string   `json:"category"`
	Status         string   `json:"status"`
	Detected       bool     `json:"detected"`
	Required       bool     `json:"required"`
	Version        string   `json:"version,omitempty"`
	Path           string   `json:"path,omitempty"`
	Message        string   `json:"message,omitempty"`
	Advice         string   `json:"advice,omitempty"`
	Installable    bool     `json:"installable"`
	InstallLabel   string   `json:"install_label,omitempty"`
	CommandPreview []string `json:"command_preview,omitempty"`
}

type environmentInstallRequest struct {
	Tool      string `json:"tool"`
	Confirmed bool   `json:"confirmed"`
}

type environmentInstallResponse struct {
	Tool           string   `json:"tool"`
	OK             bool     `json:"ok"`
	CommandPreview []string `json:"command_preview"`
	Output         string   `json:"output,omitempty"`
	Error          string   `json:"error,omitempty"`
	CheckedAt      string   `json:"checked_at"`
}

type envCommand struct {
	Command string
	Prefix  []string
	Path    string
}

type installPlan struct {
	Tool           string
	Command        string
	Args           []string
	CommandPreview []string
}

const environmentInstallTimeout = 10 * time.Minute

func (s *Server) handleEnvironmentCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, s.buildEnvironmentCheckCached())
}

// buildEnvironmentCheckCached 带缓存的环境检测，10秒内复用。
func (s *Server) buildEnvironmentCheckCached() environmentCheckResponse {
	const cacheTTL = 10 * time.Second

	s.cacheMu.RLock()
	if s.envCheckCache != nil && time.Since(s.envCheckCache.loadedAt) < cacheTTL {
		data := s.envCheckCache.data
		s.cacheMu.RUnlock()
		return data
	}
	s.cacheMu.RUnlock()

	data := s.buildEnvironmentCheck()

	s.cacheMu.Lock()
	s.envCheckCache = &envCheckCacheEntry{data: data, loadedAt: time.Now()}
	s.cacheMu.Unlock()
	return data
}

func (s *Server) handleEnvironmentCheckOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	tool := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/environment/check/"), "/")
	if tool == "" {
		errJSON(w, http.StatusBadRequest, "tool id is required")
		return
	}
	for _, group := range s.buildEnvironmentCheck().Groups {
		for _, item := range group.Items {
			if item.ID == tool {
				writeJSON(w, http.StatusOK, item)
				return
			}
		}
	}
	errJSON(w, http.StatusNotFound, "unknown environment tool: "+tool)
}

func (s *Server) handleEnvironmentInstall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req environmentInstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errJSON(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}
	req.Tool = strings.TrimSpace(req.Tool)
	if req.Tool == "" {
		errJSON(w, http.StatusBadRequest, "tool is required")
		return
	}
	if !req.Confirmed {
		errJSON(w, http.StatusBadRequest, "installation requires explicit confirmation")
		return
	}

	plan, err := buildInstallPlan(req.Tool)
	if err != nil {
		errJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), environmentInstallTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, plan.Command, plan.Args...)
	hideCommandWindow(cmd)
	out, runErr := cmd.CombinedOutput()
	output := trimCommandOutput(string(out), 12000)
	resp := environmentInstallResponse{
		Tool:           req.Tool,
		OK:             runErr == nil,
		CommandPreview: plan.CommandPreview,
		Output:         output,
		CheckedAt:      time.Now().Format(time.RFC3339Nano),
	}
	if runErr != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			resp.Error = "installation timed out"
		} else {
			resp.Error = runErr.Error()
		}
		writeJSON(w, http.StatusInternalServerError, resp)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) buildEnvironmentCheck() environmentCheckResponse {
	py, pyOK, pyVersion := detectPythonCommand()

	// 检测 Docker 是否可用且镜像就绪
	dockerOK, dockerImageReady := false, false
	if _, version, _ := detectDocker(); version != "" {
		dockerOK = true
		if img := s.dockerCfg.EffectiveEvalImage(); img != "" {
			dockerImageReady, _ = detectImage(img)
		}
	}
	dockerAvailable := dockerOK && dockerImageReady

	groups := []environmentCheckGroup{
		{
			ID:    "basic",
			Label: "基础环境",
			Items: []environmentCheckItem{
				checkExecutable("go", "Go", "basic", []string{"go"}, []string{"version"}, true, "Go 1.21+ 是构建和 Go 样本评测所必需。", ""),
				checkPython(py, pyOK, pyVersion),
				checkExecutable("java", "Java", "basic", []string{"java"}, []string{"-version"}, false, "Java/JDK 17+ 用于 Java 样本评测。", "java"),
				checkExecutable("javac", "Javac", "basic", []string{"javac"}, []string{"-version"}, false, "Java 样本编译需要 JDK，而不仅是 JRE。", "jdk"),
				checkExecutable("maven", "Maven", "basic", []string{"mvn"}, []string{"-version"}, false, "Java 样本评测和 PITest 需要 Maven。", "maven"),
				checkDocker(s.dockerCfg.EffectiveEvalImage()),
			},
		},
		{
			ID:    "python",
			Label: "Python 评测",
			Items: []environmentCheckItem{
				checkPythonModule(py, pyOK, "pytest", "pytest", "python-pytest", true, "Python 测试执行需要 pytest。"),
				checkPythonModule(py, pyOK, "coverage", "coverage.py", "python-coverage", true, "Python 覆盖率统计需要 coverage.py。"),
				checkPythonModule(py, pyOK, "mutmut", "mutmut", "python-mutmut", false, "Python 变异测试需要 mutmut；Windows 本地运行不稳定时建议使用 Docker。"),
			},
		},
		{
			ID:    "go",
			Label: "Go 评测",
			Items: []environmentCheckItem{
				checkExecutable("go-mutesting", "go-mutesting", "go", []string{"go-mutesting"}, []string{"--help"}, false, "Go 变异测试需要 go-mutesting。", "go-mutesting"),
			},
		},
		{
			ID:    "java",
			Label: "Java 评测",
			Items: []environmentCheckItem{
				checkExecutable("pitest", "PITest Maven 插件", "java", []string{"mvn"}, []string{"-version"}, false, "PITest 由 Maven 在样本项目中解析执行；这里仅检查 Maven 是否可运行。", ""),
			},
		},
		{
			ID:    "cpp",
			Label: "C++ 评测",
			Items: []environmentCheckItem{
				checkExecutable("cmake", "CMake", "cpp", []string{"cmake"}, []string{"--version"}, false, "C++ 样本构建需要 CMake。", "cmake"),
				checkExecutable("clang", "Clang", "cpp", []string{"clang", "clang-19", "clang-18", "clang-15"}, []string{"--version"}, false, "C++ 编译和 Mull 变异测试需要 Clang/LLVM。", "clang"),
				checkExecutable("clangpp", "Clang++", "cpp", []string{"clang++", "clang++-19", "clang++-18", "clang++-15"}, []string{"--version"}, false, "C++ 测试编译需要 clang++。", ""),
				checkExecutable("gcov", "gcov", "cpp", []string{"gcov"}, []string{"--version"}, false, "C++ 覆盖率统计需要 gcov。", ""),
				checkExecutable("mull", "Mull", "cpp", []string{"mull-runner-19", "mull-runner-18", "mull-runner"}, []string{"--version"}, false, "C++ 变异测试需要 Mull；建议优先使用 Docker 镜像。", ""),
			},
		},
		{
			ID:    "project",
			Label: "项目目录",
			Items: []environmentCheckItem{
				checkPath("dataset-root", "datasets", "project", s.mgr.datasetRoot, true, "样本数据集根目录。"),
				checkPath("output-root", "artifacts", "project", s.mgr.outputRoot, false, "评测产物输出目录。"),
				checkPath("storage-root", "storage", "project", filepath.Dir(s.mgr.dbPath), false, "SQLite 数据库目录。"),
				checkPath("models-yaml", "models.yaml", "project", s.configPath, true, "模型配置文件。"),
				checkPath("env-file", ".env", "project", s.dockerCfg.EnvFile, false, "API Key 会同步到该 .env 文件，Docker 执行也会读取它。"),
			},
		},
	}

	// Docker 就绪时，将本地缺失的非必需工具升级为 "docker_ok"
	if dockerAvailable {
		for gi := range groups {
			for ii := range groups[gi].Items {
				item := &groups[gi].Items[ii]
				if item.Status == "missing" && !item.Required {
					item.Status = "docker_ok"
					item.Message = "本地未安装，但 Docker 镜像中可用。"
				}
			}
		}
	}

	return environmentCheckResponse{
		CheckedAt:   time.Now().Format(time.RFC3339Nano),
		OS:          runtime.GOOS,
		ProjectRoot: s.dockerCfg.ProjectRoot,
		Summary:     summarizeEnvironmentGroups(groups),
		Groups:      groups,
	}
}

func checkPython(py envCommand, ok bool, version string) environmentCheckItem {
	item := environmentCheckItem{
		ID:       "python",
		Name:     "Python",
		Category: "basic",
		Required: true,
		Advice:   "Python 3 用于 Python 样本评测。",
	}
	if !ok {
		item.Status = "missing"
		item.Message = "未找到可用的 python/python3/py。"
		return item
	}
	item.Status = "ok"
	item.Detected = true
	item.Version = firstLine(version)
	item.Path = py.Path
	return item
}

func checkDocker(imageName string) environmentCheckItem {
	ok, version, errText := detectDocker()
	item := environmentCheckItem{
		ID:       "docker",
		Name:     "Docker",
		Category: "basic",
		Required: false,
		Advice:   "Docker 用于隔离执行和补齐 Windows 下 mutation/C++ 工具链。",
	}
	if !ok {
		item.Status = "missing"
		item.Message = strings.TrimSpace(errText)
		return item
	}
	item.Status = "ok"
	item.Detected = true
	item.Version = firstLine(version)
	if imageName != "" {
		present, imageID := detectImage(imageName)
		if !present {
			item.Status = "warning"
			item.Message = "Docker 可用，但镜像 " + imageName + " 尚未构建。"
			item.Advice = "可使用顶部“构建镜像”按钮生成 " + imageName + "。"
		} else {
			item.Message = "镜像已就绪：" + shortID(imageID)
		}
	}
	return item
}

func checkPythonModule(py envCommand, pyOK bool, module, displayName, installTool string, required bool, advice string) environmentCheckItem {
	item := environmentCheckItem{
		ID:             installTool,
		Name:           displayName,
		Category:       "python",
		Required:       required,
		Advice:         advice,
		Installable:    pyOK,
		InstallLabel:   "安装",
		CommandPreview: installPreview(installTool, py),
	}
	if !pyOK {
		item.Status = "missing"
		item.Installable = false
		item.Message = "Python 不可用，无法检查或安装 " + displayName + "。"
		return item
	}
	out, err := runCheckCommand(py.Command, append(py.Prefix, "-m", module, "--version")...)
	if err != nil {
		item.Status = "missing"
		item.Message = "未检测到 Python 模块 " + module + "。"
		return item
	}
	item.Status = "ok"
	item.Detected = true
	item.Version = firstLine(out)
	item.Path = py.Path
	return item
}

func checkExecutable(id, name, category string, candidates, versionArgs []string, required bool, advice, installTool string) environmentCheckItem {
	item := environmentCheckItem{
		ID:       id,
		Name:     name,
		Category: category,
		Required: required,
		Advice:   advice,
	}
	for _, candidate := range candidates {
		path, err := exec.LookPath(candidate)
		if err != nil {
			continue
		}
		item.Status = "ok"
		item.Detected = true
		item.Path = path
		if len(versionArgs) > 0 {
			if out, err := runCheckCommand(candidate, versionArgs...); err == nil {
				item.Version = firstLine(out)
			}
		}
		return item
	}
	item.Status = "missing"
	item.Message = "未找到命令：" + strings.Join(candidates, " / ")
	if installTool != "" {
		if plan, err := buildInstallPlan(installTool); err == nil {
			item.ID = installTool
			item.Installable = true
			item.InstallLabel = "安装"
			item.CommandPreview = plan.CommandPreview
		}
	}
	return item
}

func checkPath(id, name, category, path string, required bool, advice string) environmentCheckItem {
	item := environmentCheckItem{
		ID:       id,
		Name:     name,
		Category: category,
		Required: required,
		Path:     path,
		Advice:   advice,
	}
	if strings.TrimSpace(path) == "" {
		item.Status = "unknown"
		item.Message = "路径未配置。"
		return item
	}
	if _, err := os.Stat(path); err != nil {
		item.Status = "missing"
		item.Message = "路径不存在。"
		return item
	}
	item.Status = "ok"
	item.Detected = true
	return item
}

func summarizeEnvironmentGroups(groups []environmentCheckGroup) environmentCheckSummary {
	var sum environmentCheckSummary
	for _, group := range groups {
		for _, item := range group.Items {
			sum.Total++
			if item.Installable && item.Status != "ok" {
				sum.Installable++
			}
			switch item.Status {
			case "ok", "docker_ok":
				sum.OK++
			case "warning":
				sum.Warning++
			case "missing":
				sum.Missing++
			default:
				sum.Unknown++
			}
		}
	}
	return sum
}

func detectPythonCommand() (envCommand, bool, string) {
	candidates := []envCommand{
		{Command: "python"},
		{Command: "python3"},
		{Command: "py", Prefix: []string{"-3"}},
	}
	for _, c := range candidates {
		path, err := exec.LookPath(c.Command)
		if err != nil {
			continue
		}
		out, err := runCheckCommand(c.Command, append(c.Prefix, "--version")...)
		if err != nil {
			continue
		}
		c.Path = path
		return c, true, out
	}
	return envCommand{}, false, ""
}

func runCheckCommand(command string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), detectTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, command, args...)
	hideCommandWindow(cmd)
	out, err := cmd.CombinedOutput()
	return trimCommandOutput(string(out), 1200), err
}

func buildInstallPlan(tool string) (installPlan, error) {
	tool = strings.TrimSpace(tool)
	py, pyOK, _ := detectPythonCommand()
	switch tool {
	case "python-pytest":
		if !pyOK {
			return installPlan{}, errors.New("python is required before installing pytest")
		}
		return pythonInstallPlan(tool, py, "pytest"), nil
	case "python-coverage":
		if !pyOK {
			return installPlan{}, errors.New("python is required before installing coverage")
		}
		return pythonInstallPlan(tool, py, "coverage"), nil
	case "python-mutmut":
		if !pyOK {
			return installPlan{}, errors.New("python is required before installing mutmut")
		}
		return pythonInstallPlan(tool, py, "mutmut"), nil
	case "go-mutesting":
		if _, err := exec.LookPath("go"); err != nil {
			return installPlan{}, errors.New("go is required before installing go-mutesting")
		}
		args := []string{"install", "github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest"}
		return installPlan{
			Tool:           tool,
			Command:        "go",
			Args:           args,
			CommandPreview: append([]string{"go"}, args...),
		}, nil
	case "java", "jdk":
		return osPkgInstallPlan(tool, map[string][]string{
			"darwin":  {"brew", "install", "--cask", "temurin"},
			"linux":   {"sudo", "apt-get", "install", "-y", "default-jdk"},
			"windows": {"winget", "install", "--id", "Microsoft.OpenJDK.21", "--accept-source-agreements", "--accept-package-agreements"},
		})
	case "maven":
		return osPkgInstallPlan(tool, map[string][]string{
			"darwin":  {"brew", "install", "maven"},
			"linux":   {"sudo", "apt-get", "install", "-y", "maven"},
			"windows": {"winget", "install", "--id", "Apache.Maven", "--accept-source-agreements", "--accept-package-agreements"},
		})
	case "cmake":
		return osPkgInstallPlan(tool, map[string][]string{
			"darwin":  {"brew", "install", "cmake"},
			"linux":   {"sudo", "apt-get", "install", "-y", "cmake"},
			"windows": {"winget", "install", "--id", "Kitware.CMake", "--accept-source-agreements", "--accept-package-agreements"},
		})
	case "clang":
		return osPkgInstallPlan(tool, map[string][]string{
			"darwin": {"brew", "install", "llvm"},
			"linux":  {"sudo", "apt-get", "install", "-y", "clang"},
		})
	default:
		return installPlan{}, errors.New("tool is not installable by UTBench: " + tool)
	}
}

// osPkgInstallPlan returns a platform-specific package-manager install plan.
// cmds maps runtime.GOOS to the full command+args slice (first element is the binary).
func osPkgInstallPlan(tool string, cmds map[string][]string) (installPlan, error) {
	parts, ok := cmds[runtime.GOOS]
	if !ok {
		return installPlan{}, fmt.Errorf("no install plan for %s on %s", tool, runtime.GOOS)
	}
	if len(parts) < 2 {
		return installPlan{}, fmt.Errorf("invalid install plan for %s", tool)
	}
	return installPlan{
		Tool:           tool,
		Command:        parts[0],
		Args:           parts[1:],
		CommandPreview: parts,
	}, nil
}

func pythonInstallPlan(tool string, py envCommand, pkg string) installPlan {
	args := append(append([]string{}, py.Prefix...), "-m", "pip", "install", "-U", pkg)
	return installPlan{
		Tool:           tool,
		Command:        py.Command,
		Args:           args,
		CommandPreview: append([]string{py.Command}, args...),
	}
}

func installPreview(tool string, py envCommand) []string {
	plan, err := buildInstallPlan(tool)
	if err == nil {
		return plan.CommandPreview
	}
	if strings.HasPrefix(tool, "python-") && py.Command != "" {
		pkg := strings.TrimPrefix(tool, "python-")
		return append([]string{py.Command}, append(append([]string{}, py.Prefix...), "-m", "pip", "install", "-U", pkg)...)
	}
	return nil
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return s
}

func shortID(id string) string {
	id = strings.TrimSpace(id)
	if len(id) <= 20 {
		return id
	}
	return id[:20]
}

func trimCommandOutput(s string, max int) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\r\n", "\n"))
	if len(s) <= max {
		return s
	}
	return s[:max] + "\n... output truncated ..."
}
