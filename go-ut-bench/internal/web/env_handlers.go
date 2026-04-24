package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// handleEnv returns the host environment status (docker + image + tools).
// The frontend polls this to decide whether to show "Build Image" prompts
// or to auto-enable the Docker execution toggle.
func (s *Server) handleEnv(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	st := DetectEnv(s.dockerCfg.ImageName, s.dockerCfg.ProjectRoot, s.dockerCfg.EnvFile)
	writeJSON(w, http.StatusOK, st)
}

// handleBuildImage handles:
//   POST /api/env/build-image  → start a `docker build` job, return build_id
//   GET  /api/env/build-image  → return the latest build's status
func (s *Server) handleBuildImage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if !isDockerReady(s.dockerCfg) {
			errJSON(w, http.StatusPreconditionFailed, "docker daemon is not available on the host")
			return
		}
		job := s.bld.Submit(s.dockerCfg.ImageName)
		writeJSON(w, http.StatusCreated, map[string]string{
			"build_id":   job.BuildID,
			"image_name": job.ImageName,
			"status":     string(job.Status),
		})
	case http.MethodGet:
		job := s.bld.Latest()
		if job == nil {
			writeJSON(w, http.StatusOK, map[string]any{"build_id": "", "status": "idle"})
			return
		}
		writeJSON(w, http.StatusOK, buildJobSnapshot(job))
	default:
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleBuildImageSub routes:
//   GET /api/env/build-image/{id}         → snapshot + log buffer
//   GET /api/env/build-image/{id}/events  → SSE stream of build log lines
func (s *Server) handleBuildImageSub(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/env/build-image/")
	parts := strings.SplitN(path, "/", 2)
	id := parts[0]
	sub := ""
	if len(parts) > 1 {
		sub = parts[1]
	}

	job, ok := s.bld.Get(id)
	if !ok {
		errJSON(w, http.StatusNotFound, "build not found: "+id)
		return
	}

	switch sub {
	case "events":
		s.streamBuildEvents(w, r, job)
	default:
		snap := buildJobSnapshot(job)
		snap["logs"] = job.GetLogs()
		writeJSON(w, http.StatusOK, snap)
	}
}

// streamBuildEvents streams the build log via Server-Sent Events, mirroring
// the shape used by /api/runs/{id}/events.
func (s *Server) streamBuildEvents(w http.ResponseWriter, r *http.Request, job *BuildJob) {
	flusher, canFlush := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ch := job.Subscribe()
	defer job.Unsubscribe(ch)

	sendEvent := func(typ, payload string) {
		fmt.Fprintf(w, "data: {\"type\":%q,\"payload\":%q}\n\n", typ, payload)
		if canFlush {
			flusher.Flush()
		}
	}

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-job.Done:
			for {
				select {
				case line, ok := <-ch:
					if !ok {
						goto streamDone
					}
					sendEvent("log", line)
				default:
					goto streamDone
				}
			}
		streamDone:
			job.mu.RLock()
			status := string(job.Status)
			job.mu.RUnlock()
			sendEvent("done", status)
			return
		case line, ok := <-ch:
			if !ok {
				return
			}
			sendEvent("log", line)
		}
	}
}

// buildJobSnapshot renders a BuildJob into a JSON-serializable map without
// exposing internal fields like the subscriber slice.
func buildJobSnapshot(job *BuildJob) map[string]any {
	job.mu.RLock()
	defer job.mu.RUnlock()
	m := map[string]any{
		"build_id":   job.BuildID,
		"image_name": job.ImageName,
		"status":     string(job.Status),
		"started_at": job.StartedAt.Format(time.RFC3339Nano),
		"error":      job.Error,
	}
	if job.EndedAt != nil {
		m["ended_at"] = job.EndedAt.Format(time.RFC3339Nano)
	}
	return m
}

// isDockerReady returns true if a docker daemon is reachable. It is a
// light-weight check; we reuse DetectEnv so failure modes stay consistent
// with the /api/env endpoint.
func isDockerReady(cfg DockerConfig) bool {
	return DetectEnv(cfg.ImageName, cfg.ProjectRoot, cfg.EnvFile).DockerAvailable
}
