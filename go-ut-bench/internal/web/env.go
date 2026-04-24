package web

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// EnvStatus describes the host environment relevant to benchmark execution.
// Returned from GET /api/env so the frontend can decide whether to recommend
// Docker execution, show a "Build Image" button, etc.
type EnvStatus struct {
	OS              string          `json:"os"`
	ProjectRoot     string          `json:"project_root"`
	DockerAvailable bool            `json:"docker_available"`
	DockerVersion   string          `json:"docker_version,omitempty"`
	DockerError     string          `json:"docker_error,omitempty"`
	ImageName       string          `json:"image_name"`
	ImagePresent    bool            `json:"image_present"`
	ImageID         string          `json:"image_id,omitempty"`
	EnvFilePresent  bool            `json:"env_file_present"`
	EnvFilePath     string          `json:"env_file_path,omitempty"`
	NativeTools     map[string]bool `json:"native_tools"`
	// Recommendation is a short hint string the frontend may display, e.g.
	//   "Mutation testing on Windows requires Docker; image not found."
	Recommendation string `json:"recommendation,omitempty"`
}

// detectTimeout bounds every external command we shell out to.
const detectTimeout = 4 * time.Second

// DetectEnv probes the host for Docker + the image + native tool chain.
// It never returns an error; any per-check failure is recorded in the fields.
func DetectEnv(imageName, projectRoot, envFilePath string) EnvStatus {
	st := EnvStatus{
		OS:          runtime.GOOS,
		ProjectRoot: projectRoot,
		ImageName:   imageName,
		NativeTools: detectNativeTools(),
	}
	st.DockerAvailable, st.DockerVersion, st.DockerError = detectDocker()
	if st.DockerAvailable {
		st.ImagePresent, st.ImageID = detectImage(imageName)
	}
	if envFilePath != "" {
		if _, err := os.Stat(envFilePath); err == nil {
			st.EnvFilePresent = true
			st.EnvFilePath = envFilePath
		}
	}
	st.Recommendation = buildRecommendation(st)
	return st
}

// detectDocker runs `docker version --format {{.Server.Version}}`.
func detectDocker() (bool, string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), detectTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, "docker", "version", "--format", "{{.Server.Version}}").CombinedOutput()
	if err != nil {
		return false, "", strings.TrimSpace(string(out)) + " " + err.Error()
	}
	return true, strings.TrimSpace(string(out)), ""
}

// detectImage checks whether the named image exists locally.
// `docker image inspect <name> --format {{.Id}}` exits non-zero if missing.
func detectImage(name string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), detectTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, "docker", "image", "inspect", name, "--format", "{{.Id}}").Output()
	if err != nil {
		return false, ""
	}
	return true, strings.TrimSpace(string(out))
}

// detectNativeTools probes common evaluator toolchains via exec.LookPath.
// These are advisory only; failure to find a tool just tells the UI to
// suggest Docker execution.
func detectNativeTools() map[string]bool {
	tools := []string{
		"python3", "pytest", "coverage", "mutmut",
		"go", "gremlins",
		"mvn", "java", "javac",
		"clang", "clang++", "cmake", "mull",
	}
	out := make(map[string]bool, len(tools))
	for _, t := range tools {
		_, err := exec.LookPath(t)
		out[t] = err == nil
	}
	return out
}

// buildRecommendation crafts a short advisory message for the UI.
func buildRecommendation(st EnvStatus) string {
	switch {
	case !st.DockerAvailable:
		if st.OS == "windows" && !st.NativeTools["mutmut"] {
			return "Docker is not available. Mutation testing (mutmut) is not supported natively on Windows; install Docker Desktop or use WSL."
		}
		return ""
	case !st.ImagePresent:
		return "Docker is available but the image '" + st.ImageName + "' is not built yet. Click 'Build Image' to build it once."
	case st.OS == "windows" && !st.NativeTools["mutmut"]:
		return "Docker image is ready. Recommend enabling 'Execute in Docker' when running mutation tests on Windows."
	default:
		return ""
	}
}
