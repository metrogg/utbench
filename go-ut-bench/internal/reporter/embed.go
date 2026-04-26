package reporter

import (
	_ "embed"
	"strings"
)

// 嵌入静态资源文件
//
//go:embed static/style.css
var styleCSS string

//go:embed static/charts.js
var chartsJSTemplate string

// GetStyleCSS 返回 CSS 样式内容
func GetStyleCSS() string {
	return styleCSS
}

// BuildChartsJS 构建包含数据的完整 JavaScript
func BuildChartsJS(topModelsJSON, rowsJSON string) string {
	js := chartsJSTemplate
	js = strings.Replace(js, "__TOP_MODELS_JSON_PLACEHOLDER__", topModelsJSON, 1)
	js = strings.Replace(js, "__ROWS_JSON_PLACEHOLDER__", rowsJSON, 1)
	return js
}