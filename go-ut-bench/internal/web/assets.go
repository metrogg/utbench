// 资产管理 API：扫描 artifacts/runs/ 目录，返回每个 run 的完整磁盘视图，
// 用于「评测资产」页面的筛选与清理（避免脏数据混入排行榜 / 数据库）。
//
// 与 /api/runs 的差异：
//   - /api/runs 只关心可恢复成 RunSummary 的条目，缺少 run_summary.json 的目录会被跳过；
//     activeRuns 会覆盖磁盘条目（用于实时状态）。
//   - /api/assets/runs 关心磁盘上**所有** run 目录，哪怕 run_summary.json 缺失或损坏，
//     也会作为「dirty」条目返回，方便用户决定是否删除。
package web

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// AssetItem 描述一条评测资产（一个 artifacts/runs/<run_id>/ 目录）。
type AssetItem struct {
	RunID         string   `json:"run_id"`
	Label         string   `json:"label,omitempty"`
	Status        string   `json:"status"`        // completed / failed / unknown / dirty
	Phase         string   `json:"phase"`         // full / generate / evaluate / report
	StartedAt     string   `json:"started_at"`    // RFC3339
	SizeBytes     int64    `json:"size_bytes"`    // 整个 run 目录的累计大小
	HasSummary    bool     `json:"has_summary"`   // run_summary.json 存在且可解析
	HasManifest   bool     `json:"has_manifest"`  // generated/generated_manifest.json
	HasEvaluation bool     `json:"has_evaluation"`
	HasReportJSON bool     `json:"has_report_json"`
	HasReportHTML bool     `json:"has_report_html"`
	Models        []string `json:"models,omitempty"`
	Subjects      []string `json:"subjects,omitempty"`
	Languages     []string `json:"languages,omitempty"`
	MaxSamples    int      `json:"max_samples,omitempty"`
	Workers       int      `json:"workers,omitempty"`
	IsMerged      bool     `json:"is_merged_report,omitempty"`
	// SampleCount 是从 generated_manifest.json 或 evaluation_result.json 中读取到的样本数。
	SampleCount int `json:"sample_count,omitempty"`
	// Ingested 表示该 run_id 是否已入库到 SQLite（generation_runs 表中存在）。
	Ingested bool `json:"ingested"`
	// QualityFlag 用于前端筛选脏数据：
	//   "ok"        - 有 summary、有 evaluation、有 report
	//   "partial"   - 有 summary 但缺 evaluation 或 report
	//   "dirty"     - summary 缺失/损坏
	//   "running"   - 内存中仍在跑（不会出现在磁盘扫描里，仅占位）
	QualityFlag string `json:"quality_flag"`
	// Issues 列出导致 dirty/partial 的具体原因，前端 tooltip 用。
	Issues []string `json:"issues,omitempty"`
}

// handleAssets 处理 GET /api/assets/runs，返回 artifacts/runs/ 全量扫描结果。
func (s *Server) handleAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	root := filepath.Join(s.outputRoot, "runs")
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			writeJSON(w, http.StatusOK, []AssetItem{})
			return
		}
		errJSON(w, http.StatusInternalServerError, "read runs dir: "+err.Error())
		return
	}

	out := make([]AssetItem, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		runID := e.Name()
		runDir := filepath.Join(root, runID)
		out = append(out, scanAssetDir(runID, runDir))
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].StartedAt > out[j].StartedAt // ISO 字符串字典序==时间倒序
	})

	// 批量查 DB，标记哪些 run_id 已入库。DB 不可用时静默跳过（ingested 默认 false）。
	if s.db != nil {
		runIDs := make([]string, len(out))
		for i, a := range out {
			runIDs[i] = a.RunID
		}
		if ingested, err := s.db.IngestedRunIDs(r.Context(), runIDs); err == nil {
			for i := range out {
				out[i].Ingested = ingested[out[i].RunID]
			}
		}
	}

	writeJSON(w, http.StatusOK, out)
}

// scanAssetDir 扫描单个 run 目录，提取关键文件存在性、大小、以及来自 run_summary.json 的元信息。
func scanAssetDir(runID, runDir string) AssetItem {
	item := AssetItem{
		RunID:       runID,
		Status:      "unknown",
		QualityFlag: "dirty",
	}

	// 累计目录大小（递归）。失败不阻塞。
	item.SizeBytes = dirSize(runDir)

	// 关键文件存在性检查
	summaryPath := filepath.Join(runDir, "run_summary.json")
	manifestPath := filepath.Join(runDir, "generated", "generated_manifest.json")
	evaluationPath := filepath.Join(runDir, "evaluation", "evaluation_result.json")
	reportJSONPath := filepath.Join(runDir, "report", "report_summary.json")
	reportHTMLPath := filepath.Join(runDir, "report", "report.html")

	item.HasSummary = fileNonEmpty(summaryPath)
	item.HasManifest = fileNonEmpty(manifestPath)
	item.HasEvaluation = fileNonEmpty(evaluationPath)
	item.HasReportJSON = fileNonEmpty(reportJSONPath)
	item.HasReportHTML = fileNonEmpty(reportHTMLPath)

	// 解析 run_summary.json 提取 spec 元信息
	if item.HasSummary {
		if data, err := os.ReadFile(summaryPath); err == nil {
			var raw struct {
				RunID          string `json:"run_id"`
				Label          string `json:"label"`
				CreatedAtUTC   string `json:"created_at_utc"`
				Phase          string `json:"phase"`
				IsMergedReport bool   `json:"is_merged_report"`
				Spec           struct {
					Models     []string `json:"models"`
					Subjects   []string `json:"subjects"`
					Languages  []string `json:"languages"`
					MaxSamples int      `json:"max_samples"`
					Workers    int      `json:"workers"`
				} `json:"spec"`
			}
			if err := json.Unmarshal(data, &raw); err == nil {
				item.Label = strings.TrimSpace(raw.Label)
				item.Phase = raw.Phase
				item.IsMerged = raw.IsMergedReport
				item.Models = raw.Spec.Models
				item.Subjects = raw.Spec.Subjects
				item.Languages = raw.Spec.Languages
				item.MaxSamples = raw.Spec.MaxSamples
				item.Workers = raw.Spec.Workers
				if t, err := time.Parse(time.RFC3339Nano, raw.CreatedAtUTC); err == nil {
					item.StartedAt = t.UTC().Format(time.RFC3339)
				}
			} else {
				item.Issues = append(item.Issues, "run_summary.json 解析失败: "+err.Error())
				item.HasSummary = false
			}
		}
	}

	// 从 generated_manifest.json 中读取实际生成的样本数
	if item.HasManifest {
		if data, err := os.ReadFile(manifestPath); err == nil {
			var manifest struct {
				Cases []json.RawMessage `json:"cases"`
			}
			if json.Unmarshal(data, &manifest) == nil {
				item.SampleCount = len(manifest.Cases)
			}
		}
	}
	// 若 manifest 缺失但有 evaluation，则从 evaluation 中补充样本数
	if item.SampleCount == 0 && item.HasEvaluation {
		if data, err := os.ReadFile(evaluationPath); err == nil {
			var eval struct {
				Results []json.RawMessage `json:"results"`
			}
			if json.Unmarshal(data, &eval) == nil {
				item.SampleCount = len(eval.Results)
			}
		}
	}

	// 用目录的 ModTime 兜底 StartedAt
	if item.StartedAt == "" {
		if fi, err := os.Stat(runDir); err == nil {
			item.StartedAt = fi.ModTime().UTC().Format(time.RFC3339)
		}
	}

	// 状态与质量分类
	switch {
	case !item.HasSummary:
		item.Status = "dirty"
		item.QualityFlag = "dirty"
		item.Issues = append(item.Issues, "缺少 run_summary.json")
	case item.HasReportHTML && item.HasEvaluation:
		item.Status = "completed"
		item.QualityFlag = "ok"
	case item.HasManifest && !item.HasEvaluation:
		item.Status = "generated_only"
		item.QualityFlag = "partial"
		item.Issues = append(item.Issues, "仅完成 generate，未跑 evaluate")
	case item.HasEvaluation && !item.HasReportHTML:
		item.Status = "evaluated_only"
		item.QualityFlag = "partial"
		item.Issues = append(item.Issues, "已 evaluate，未生成 HTML 报告")
	default:
		item.Status = "incomplete"
		item.QualityFlag = "partial"
		if !item.HasManifest {
			item.Issues = append(item.Issues, "缺少 generated_manifest.json")
		}
		if !item.HasEvaluation {
			item.Issues = append(item.Issues, "缺少 evaluation_result.json")
		}
		if !item.HasReportHTML {
			item.Issues = append(item.Issues, "缺少 report.html")
		}
	}

	return item
}

// fileNonEmpty 检查文件存在且非空。
func fileNonEmpty(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fi.IsDir() && fi.Size() > 0
}

// dirSize 返回目录递归总字节数。无法访问的子项被跳过。
func dirSize(root string) int64 {
	var total int64
	_ = filepath.WalkDir(root, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if info, e := d.Info(); e == nil {
			total += info.Size()
		}
		return nil
	})
	return total
}
