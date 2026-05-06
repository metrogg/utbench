package evaluator

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// EnvironmentFingerprint 评测环境指纹
type EnvironmentFingerprint struct {
	OS              string            `json:"os"`
	Arch            string            `json:"arch"`
	GoVersion       string            `json:"go_version,omitempty"`
	PythonVersion   string            `json:"python_version,omitempty"`
	PytestVersion   string            `json:"pytest_version,omitempty"`
	MutmutVersion   string            `json:"mutmut_version,omitempty"`
	CoverageVersion string            `json:"coverage_version,omitempty"`
	JavaVersion     string            `json:"java_version,omitempty"`
	JavacVersion    string            `json:"javac_version,omitempty"`
	MavenVersion    string            `json:"maven_version,omitempty"`
	CPPVersion      string            `json:"cpp_version,omitempty"`
	GCCVersion      string            `json:"gcc_version,omitempty"`
	CMakeVersion    string            `json:"cmake_version,omitempty"`
	MullVersion     string            `json:"mull_version,omitempty"`
	GremlinsVersion string            `json:"gremlins_version,omitempty"`
	DockerDigest    string            `json:"docker_digest,omitempty"`
	DockerUsed      bool              `json:"docker_used"`
	Timestamp       string            `json:"timestamp"`
	Extra           map[string]string `json:"extra,omitempty"`
}

// CaptureEnvironmentFingerprint 采集当前评测环境指纹
func CaptureEnvironmentFingerprint(ctx context.Context, useDocker bool, dockerImage string) EnvironmentFingerprint {
	fp := EnvironmentFingerprint{
		OS:         runtime.GOOS + "/" + runtime.GOARCH,
		Arch:       runtime.GOARCH,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		DockerUsed: useDocker,
		Extra:      make(map[string]string),
	}

	// 如果使用Docker，获取镜像digest
	if useDocker && dockerImage != "" {
		fp.DockerDigest = getDockerImageDigest(ctx, dockerImage)
	} else {
		// 本地环境，采集工具版本
		fp.GoVersion = getToolVersion(ctx, "go", "version")
		fp.PythonVersion = getToolVersion(ctx, "python", "--version")
		fp.PytestVersion = getToolVersion(ctx, "pytest", "--version")
		fp.MutmutVersion = getToolVersion(ctx, "mutmut", "--version")
		fp.CoverageVersion = getToolVersion(ctx, "coverage", "--version")
		fp.JavaVersion = getToolVersion(ctx, "java", "-version")
		fp.JavacVersion = getToolVersion(ctx, "javac", "-version")
		fp.MavenVersion = getToolVersion(ctx, "mvn", "-version")
		fp.GCCVersion = getToolVersion(ctx, "gcc", "--version")
		fp.CMakeVersion = getToolVersion(ctx, "cmake", "--version")
		fp.MullVersion = getToolVersion(ctx, "mull", "--version")
		fp.GremlinsVersion = getToolVersion(ctx, "gremlins", "version")
	}

	// 采集环境变量（可选）
	if v := os.Getenv("UTBENCH_ENV_NOTE"); v != "" {
		fp.Extra["note"] = v
	}

	return fp
}

// FingerprintHash 计算指纹的稳定哈希值（用于数据库唯一标识）
func (fp EnvironmentFingerprint) FingerprintHash() string {
	// 只取影响评测结果的关键字段
	keyFields := []string{
		fp.OS,
		fp.GoVersion,
		fp.PythonVersion,
		fp.PytestVersion,
		fp.MutmutVersion,
		fp.CoverageVersion,
		fp.JavaVersion,
		fp.JavacVersion,
		fp.MavenVersion,
		fp.GCCVersion,
		fp.CMakeVersion,
		fp.MullVersion,
		fp.GremlinsVersion,
		fp.DockerDigest,
	}
	if fp.DockerUsed {
		keyFields = append(keyFields, "docker")
	}
	return stableHash(strings.Join(keyFields, "|"))
}

// ToJSON 序列化为JSON字符串
func (fp EnvironmentFingerprint) ToJSON() string {
	data, _ := json.Marshal(fp)
	return string(data)
}

// getToolVersion 获取工具版本
func getToolVersion(ctx context.Context, tool string, args ...string) string {
	cmd := exec.CommandContext(ctx, tool, args...)
	hideCommandWindow(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	// 提取第一行，清理格式
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		v := strings.TrimSpace(lines[0])
		// 简化版本字符串，只保留关键信息
		v = cleanVersionString(v)
		return v
	}
	return ""
}

// getDockerImageDigest 获取Docker镜像digest
func getDockerImageDigest(ctx context.Context, image string) string {
	cmd := exec.CommandContext(ctx, "docker", "inspect", "--format", "{{.Id}}", image)
	hideCommandWindow(cmd)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// cleanVersionString 清理版本字符串，只保留关键版本信息
func cleanVersionString(v string) string {
	// 提取版本号部分
	// 例如: "go version go1.21.5 windows/amd64" -> "go1.21.5"
	// 例如: "Python 3.10.11" -> "Python 3.10.11"
	v = strings.TrimSpace(v)
	// 限制长度
	if len(v) > 50 {
		v = v[:50]
	}
	return v
}

// stableHash 计算稳定哈希
func stableHash(s string) string {
	// 使用简单的SHA256截断
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])[:32]
}

// CompareFingerprints 比较两个指纹是否兼容（可用于跨Run对比）
func CompareFingerprints(a, b EnvironmentFingerprint) (compatible bool, differences []string) {
	differences = []string{}

	// Docker模式：比较镜像digest
	if a.DockerUsed && b.DockerUsed {
		if a.DockerDigest != b.DockerDigest {
			differences = append(differences, fmt.Sprintf("Docker镜像不同: %s vs %s", a.DockerDigest[:16], b.DockerDigest[:16]))
		}
		return len(differences) == 0, differences
	}

	// 本地模式：比较工具版本
	if a.GoVersion != b.GoVersion && a.GoVersion != "" && b.GoVersion != "" {
		differences = append(differences, fmt.Sprintf("Go版本不同: %s vs %s", a.GoVersion, b.GoVersion))
	}
	if a.PythonVersion != b.PythonVersion && a.PythonVersion != "" && b.PythonVersion != "" {
		differences = append(differences, fmt.Sprintf("Python版本不同: %s vs %s", a.PythonVersion, b.PythonVersion))
	}
	if a.PytestVersion != b.PytestVersion && a.PytestVersion != "" && b.PytestVersion != "" {
		differences = append(differences, fmt.Sprintf("pytest版本不同: %s vs %s", a.PytestVersion, b.PytestVersion))
	}
	if a.JavaVersion != b.JavaVersion && a.JavaVersion != "" && b.JavaVersion != "" {
		differences = append(differences, fmt.Sprintf("Java版本不同: %s vs %s", a.JavaVersion, b.JavaVersion))
	}

	// Docker vs 本地是重大差异
	if a.DockerUsed != b.DockerUsed {
		differences = append(differences, "执行模式不同: Docker vs 本地")
	}

	return len(differences) == 0, differences
}
