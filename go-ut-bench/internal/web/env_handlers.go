package web

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
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
	force := strings.TrimSpace(r.URL.Query().Get("refresh")) == "1"
	st := detectEnvCached(s.dockerCfg, force)
	writeJSON(w, http.StatusOK, st)
}

// envCacheEntry caches DetectEnv results to avoid spawning subprocesses on every poll.
var (
	envCacheMu    sync.Mutex
	envCacheEntry *EnvStatus
	envCacheTime  time.Time
	envCacheTTL   = 15 * time.Second
)

func detectEnvCached(cfg DockerConfig, force bool) EnvStatus {
	envCacheMu.Lock()
	defer envCacheMu.Unlock()
	if !force && envCacheEntry != nil && time.Since(envCacheTime) < envCacheTTL {
		return *envCacheEntry
	}
	st := DetectEnv(cfg)
	envCacheEntry = &st
	envCacheTime = time.Now()
	return st
}

// handleBuildImage handles:
//
//	POST /api/env/build-image  → start a `docker build` job, return build_id
//	GET  /api/env/build-image  → return the latest build's status
func (s *Server) handleBuildImage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if latest := s.bld.Latest(); latest != nil {
			latest.mu.RLock()
			active := latest.Status == BuildPending || latest.Status == BuildRunning
			latest.mu.RUnlock()
			if active {
				writeJSON(w, http.StatusOK, buildJobSnapshot(latest))
				return
			}
		}
		if !isDockerReady(s.dockerCfg) {
			errJSON(w, http.StatusPreconditionFailed, "docker daemon is not available on the host")
			return
		}
		target := strings.TrimSpace(r.URL.Query().Get("target"))
		profile, err := defaultBuildProfile(target, s.dockerCfg)
		if err != nil {
			errJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		job := s.bld.Submit(profile)
		writeJSON(w, http.StatusCreated, map[string]string{
			"build_id":   job.BuildID,
			"image_name": job.ImageName,
			"target":     job.Target,
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
//
//	GET /api/env/build-image/{id}         → snapshot + log buffer
//	GET /api/env/build-image/{id}/events  → SSE stream of build log lines
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

	if r.Method == http.MethodDelete {
		if !job.Cancel() {
			errJSON(w, http.StatusConflict, "build is not running")
			return
		}
		writeJSON(w, http.StatusOK, buildJobSnapshot(job))
		return
	}
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
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
		"target":     job.Target,
		"dockerfile": job.Dockerfile,
		"image_name": job.ImageName,
		"build_args": job.BuildArgs,
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
	return detectEnvCached(cfg, false).DockerAvailable
}
