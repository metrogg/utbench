// dataset 包提供数据集管理功能
// 负责数据集样本的发现、索引构建、清单生成和数据验证
package dataset

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go-ut-bench/internal/contracts"
)

// indexSample 数据集索引中的样本结构
// 用于存储在索引文件中的样本基本信息
type indexSample struct {
	ID        string                 `json:"id"`        // 样本唯一标识符
	Language  string                 `json:"language"`  // 编程语言
	Category  contracts.DatasetClass `json:"category"`  // 数据集类别
	Scenario  string                 `json:"scenario"`  // 场景名称
	Path      string                 `json:"path"`      // 样本文件路径
	SourceMD5 string                 `json:"source_md5"` // 源代码MD5哈希
}

// datasetIndexFile 数据集索引文件结构
// 包含索引版本号、生成时间和所有样本列表
type datasetIndexFile struct {
	SchemaVersion string        `json:"schema_version"` // 索引版本号
	GeneratedAt   time.Time     `json:"generated_at_utc"` // 生成时间
	DatasetRoot   string        `json:"dataset_root"`   // 数据集根目录
	Samples       []indexSample `json:"samples"`       // 样本列表
}

// BuildSummary 构建操作的摘要信息
// 返回构建的输出路径和处理的样本总数
type BuildSummary struct {
	Path  string // 输出文件路径
	Total int    // 处理的样本总数
}

// ManifestBuildOptions 清单构建选项
// 定义如何从索引中筛选和构建样本清单
type ManifestBuildOptions struct {
	IndexPath        string // 索引文件路径
	Level            string // 清单级别，如"l1"
	OutputPath       string // 输出清单文件路径
	Languages        []string // 要包含的编程语言列表
	ClassFilter      string // 数据集类别过滤，如"self_contained"
	ScenarioFilter   string // 场景过滤，如"boundary"
	LimitPerScenario int    // 每个场景最大样本数（默认20）
}

// BuildIndex 构建数据集索引文件
// 参数:
//   - datasetRoot: 数据集根目录路径
//   - outputPath: 索引输出文件路径
//
// 返回值:
//   - BuildSummary: 构建摘要（输出路径和样本数）
//   - error: 构建失败时的错误
//
// 功能说明:
//   - 扫描数据集根目录下的所有样本文件
//   - 生成包含所有样本信息的索引文件
//   - 样本按语言、类别、场景、ID排序
func (s *Service) BuildIndex(datasetRoot, outputPath string) (BuildSummary, error) {
	spec := contracts.RunSpec{
		RunID:       contracts.NewRunID(),
		DatasetRoot: strings.TrimSpace(datasetRoot),
		Languages:   append([]string{}, contracts.SupportedLanguages...),
		MaxSamples:  0,
	}
	samples, err := s.DiscoverSamples(spec)
	if err != nil {
		return BuildSummary{}, err
	}

	// 将样本转换为索引格式
	rows := make([]indexSample, 0, len(samples))
	for _, item := range samples {
		rows = append(rows, indexSample{
			ID:        item.ID,
			Language:  item.Language,
			Category:  item.Category,
			Scenario:  item.Scenario,
			Path:      item.Path,
			SourceMD5: item.SourceMD5,
		})
	}

	// 多级排序：语言 -> 类别 -> 场景 -> ID
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Language == rows[j].Language {
			if rows[i].Category == rows[j].Category {
				if rows[i].Scenario == rows[j].Scenario {
					return rows[i].ID < rows[j].ID
				}
				return rows[i].Scenario < rows[j].Scenario
			}
			return rows[i].Category < rows[j].Category
		}
		return rows[i].Language < rows[j].Language
	})

	// 构建并写入索引文件
	payload := datasetIndexFile{
		SchemaVersion: contracts.SchemaVersion,
		GeneratedAt:   time.Now().UTC(),
		DatasetRoot:   spec.DatasetRoot,
		Samples:       rows,
	}
	if err := contracts.WriteJSON(outputPath, payload); err != nil {
		return BuildSummary{}, err
	}
	return BuildSummary{Path: outputPath, Total: len(rows)}, nil
}

// BuildManifest 从索引文件构建样本清单
// 参数:
//   - opts: ManifestBuildOptions，构建选项
//
// 返回值:
//   - BuildSummary: 构建摘要
//   - error: 构建失败时的错误
//
// 功能说明:
//   - 根据选项从索引中筛选样本
//   - 支持按语言、类别、场景过滤
//   - 每个场景限制样本数量
func (s *Service) BuildManifest(opts ManifestBuildOptions) (BuildSummary, error) {
	idx, err := readDatasetIndex(opts.IndexPath)
	if err != nil {
		return BuildSummary{}, err
	}

	// 处理语言过滤
	langs := map[string]struct{}{}
	for _, lang := range opts.Languages {
		lang = strings.ToLower(strings.TrimSpace(lang))
		if lang != "" {
			langs[lang] = struct{}{}
		}
	}
	classFilter := strings.TrimSpace(opts.ClassFilter)
	scenarioFilter := normalizeScenario(opts.ScenarioFilter)
	limit := opts.LimitPerScenario
	if limit <= 0 {
		limit = 20 // 默认每个场景最多20个样本
	}

	// 定义分组键
	type groupKey struct {
		Lang     string
		Category contracts.DatasetClass
		Scenario string
	}
	grouped := map[groupKey][]indexSample{}

	// 遍历索引，应用过滤条件
	for _, item := range idx.Samples {
		// 语言过滤
		if len(langs) > 0 {
			if _, ok := langs[item.Language]; !ok {
				continue
			}
		}
		// 类别过滤
		if classFilter != "" {
			classFilters := strings.Split(classFilter, ",")
			for i := range classFilters {
				classFilters[i] = strings.TrimSpace(classFilters[i])
			}
			if !matchDatasetClassFilter(classFilters, item.Category) {
				continue
			}
		}
		// 场景过滤
		if scenarioFilter != "" && item.Scenario != scenarioFilter {
			continue
		}
		k := groupKey{Lang: item.Language, Category: item.Category, Scenario: item.Scenario}
		grouped[k] = append(grouped[k], item)
	}

	// 排序分组键
	keys := make([]groupKey, 0, len(grouped))
	for k := range grouped {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Lang == keys[j].Lang {
			if keys[i].Category == keys[j].Category {
				return keys[i].Scenario < keys[j].Scenario
			}
			return keys[i].Category < keys[j].Category
		}
		return keys[i].Lang < keys[j].Lang
	})

	// 构建样本列表，应用数量限制
	samples := make([]map[string]any, 0)
	for _, key := range keys {
		rows := grouped[key]
		sort.Slice(rows, func(i, j int) bool { return rows[i].ID < rows[j].ID })
		if len(rows) > limit {
			rows = rows[:limit]
		}
		for _, row := range rows {
			samples = append(samples, map[string]any{
				"id":       row.ID,
				"language": row.Language,
				"category": row.Category,
				"scenario": row.Scenario,
				"path":     row.Path,
			})
		}
	}

	if len(samples) == 0 {
		return BuildSummary{}, fmt.Errorf("no samples selected from index")
	}

	// 设置默认级别
	level := strings.TrimSpace(opts.Level)
	if level == "" {
		level = "l1"
	}
	payload := map[string]any{
		"level":   level,
		"samples": samples,
	}
	if err := contracts.WriteJSON(opts.OutputPath, payload); err != nil {
		return BuildSummary{}, err
	}
	return BuildSummary{Path: opts.OutputPath, Total: len(samples)}, nil
}

// readDatasetIndex 读取数据集索引文件
// 参数:
//   - path: 索引文件路径
//
// 返回值:
//   - datasetIndexFile: 解析后的索引结构
//   - error: 读取或解析失败时的错误
func readDatasetIndex(path string) (datasetIndexFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return datasetIndexFile{}, err
	}
	var out datasetIndexFile
	if err := json.Unmarshal(raw, &out); err != nil {
		return datasetIndexFile{}, err
	}
	// 规范化路径
	for i := range out.Samples {
		if !filepath.IsAbs(out.Samples[i].Path) {
			out.Samples[i].Path = filepath.Clean(out.Samples[i].Path)
		}
	}
	return out, nil
}
