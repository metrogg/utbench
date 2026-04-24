package evaluator

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestMullInstallation 验证 Mull 是否正确安装
func TestMullInstallation(t *testing.T) {
	runner := findMullRunner()
	if runner == "" {
		t.Skip("Mull 未安装，跳过测试。安装命令: curl -1sLf 'https://dl.cloudsmith.io/public/mull-project/mull-stable/setup.deb.sh' | bash && apt-get install -y mull-19")
	}
	t.Logf("✅ 找到 Mull runner: %s", runner)

	frontend := findMullFrontend()
	if frontend == "" {
		t.Error("❌ mull-ir-frontend 未找到")
	} else {
		t.Logf("✅ 找到 Mull IR 前端: %s", frontend)
	}

	// 验证版本
	cmd := exec.Command(runner, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("⚠️  无法获取版本信息: %v", err)
	} else {
		t.Logf("📋 Mull 版本信息:\n%s", string(output))
	}
}

// TestCppMutationIntegration 完整的 C++ 变异测试集成测试
func TestCppMutationIntegration(t *testing.T) {
	// 检查 Mull 是否安装（支持 19/18 版本）
	runner := findMullRunner()
	if runner == "" {
		t.Skip("Mull 未安装 (mull-runner-19/18)，跳过变异测试")
	}
	t.Logf("使用 Mull: %s", runner)

	// 创建工作目录
	workdir, err := os.MkdirTemp("", "mull_test_*")
	if err != nil {
		t.Fatalf("创建工作目录失败: %v", err)
	}
	defer os.RemoveAll(workdir)

	// 创建简单的 C++ 源文件
	sourceContent := `int add(int a, int b) { return a + b; }
int multiply(int a, int b) { return a * b; }
`
	sourceFile := filepath.Join(workdir, "sample.cpp")
	if err := os.WriteFile(sourceFile, []byte(sourceContent), 0644); err != nil {
		t.Fatalf("写入源文件失败: %v", err)
	}

	// 创建测试文件
	testContent := `#include <gtest/gtest.h>
#include "sample.cpp"

TEST(AddTest, Basic) {
    EXPECT_EQ(add(2, 3), 5);
    EXPECT_EQ(add(-1, 1), 0);
}

TEST(MultiplyTest, Basic) {
    EXPECT_EQ(multiply(2, 3), 6);
    EXPECT_EQ(multiply(0, 5), 0);
}
`
	testFile := filepath.Join(workdir, "sample_test.cpp")
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("写入测试文件失败: %v", err)
	}

	t.Logf("📁 工作目录: %s", workdir)
	t.Logf("📄 源文件: %s", sourceFile)
	t.Logf("🧪 测试文件: %s", testFile)

	// 准备工作空间
	prepWorkdir, testName, sourceBase, _, prepErr := prepareCppWorkspace(testFile, sourceFile)
	if prepErr != "" {
		t.Fatalf("prepareCppWorkspace 失败: %s", prepErr)
	}
	defer os.RemoveAll(prepWorkdir)

	t.Logf("✅ 工作空间准备完成")
	t.Logf("   目录: %s", prepWorkdir)
	t.Logf("   测试: %s", testName)
	t.Logf("   源文件: %s", sourceBase)

	// 运行变异测试
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 使用 100% 测试通过率
	passRate := 1.0
	score, stats, mutErr := collectCppMutation(ctx, prepWorkdir, sourceBase, 60, &passRate, 4, 4)

	// 输出结果
	t.Logf("\n📊 变异测试结果:")
	t.Logf("   分数: %.2f%%", score*100)
	t.Logf("   总变异体: %d", stats.Total)
	t.Logf("   已杀死: %d", stats.Killed)
	t.Logf("   存活: %d", stats.Survived)
	t.Logf("   超时: %d", stats.Timeout)
	t.Logf("   无测试: %d", stats.NoTests)
	t.Logf("   未检查: %d", stats.NotChecked)
	t.Logf("   跳过: %d", stats.Skipped)
	t.Logf("   可疑: %d", stats.Suspicious)

	if mutErr != "" {
		t.Logf("⚠️  变异测试警告: %s", mutErr)
	}

	// 验证结果
	if stats.Total == 0 {
		t.Error("❌ 未产生任何变异体")
	} else {
		t.Logf("✅ 成功产生 %d 个变异体", stats.Total)
	}

	if stats.Killed+stats.Survived > 0 {
		t.Logf("✅ 有效变异体: %d (已杀死: %d, 存活: %d)",
			stats.Killed+stats.Survived, stats.Killed, stats.Survived)
	}
}

func TestPrepareCppWorkspaceForcesSourceIncludeInsteadOfSynthesizingHeader(t *testing.T) {
	workdir := t.TempDir()
	sourcePath := filepath.Join(workdir, "boundary_000.cpp")
	source := `#include <string>
class BigInt {
public:
    BigInt() {}
    std::string toString() const { return "0"; }
};
`
	if err := os.WriteFile(sourcePath, []byte(source), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	testPath := filepath.Join(workdir, "boundary_000.test.cpp")
	testSource := `#include <gtest/gtest.h>
#include "BigInt.h"

TEST(BigIntTest, DefaultConstructor) {
    BigInt a;
    EXPECT_EQ(a.toString(), "0");
}
`
	if err := os.WriteFile(testPath, []byte(testSource), 0o644); err != nil {
		t.Fatalf("write test: %v", err)
	}

	prepWorkdir, testName, _, _, prepErr := prepareCppWorkspace(testPath, sourcePath)
	if prepErr != "" {
		t.Fatalf("prepareCppWorkspace failed: %s", prepErr)
	}
	defer os.RemoveAll(prepWorkdir)

	if _, err := os.Stat(filepath.Join(prepWorkdir, "BigInt.h")); !os.IsNotExist(err) {
		t.Fatalf("BigInt.h should not be synthesized, stat err=%v", err)
	}

	raw, err := os.ReadFile(filepath.Join(prepWorkdir, testName))
	if err != nil {
		t.Fatalf("read prepared test: %v", err)
	}
	prepared := string(raw)
	if !strings.Contains(prepared, `#include "boundary_000.cpp"`) {
		t.Fatalf("prepared test should include source cpp, got:\n%s", prepared)
	}
	if strings.Contains(prepared, `#include "BigInt.h"`) {
		t.Fatalf("prepared test should remove generated local header include, got:\n%s", prepared)
	}
}

// TestParseMullOutput 测试 Mull 输出解析
func TestParseMullOutput(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		wantErr  bool
		validate func(t *testing.T, stats mutationStats)
	}{
		{
			name: "全部杀死",
			output: `Running mutants (threads: 4)...
All mutations have been killed
10/10 . Finished`,
			wantErr: false,
			validate: func(t *testing.T, stats mutationStats) {
				if stats.Total != 10 {
					t.Errorf("期望 Total=10, 得到 %d", stats.Total)
				}
				if stats.Killed != 10 {
					t.Errorf("期望 Killed=10, 得到 %d", stats.Killed)
				}
				if stats.Survived != 0 {
					t.Errorf("期望 Survived=0, 得到 %d", stats.Survived)
				}
			},
		},
		{
			name: "标准输出格式",
			output: `Killed mutants (8/10)
Survived mutants (2/10)
Timeout mutants (0/10)
No tests mutants (0/10)`,
			wantErr: false,
			validate: func(t *testing.T, stats mutationStats) {
				if stats.Total != 10 {
					t.Errorf("期望 Total=10, 得到 %d", stats.Total)
				}
				if stats.Killed != 8 {
					t.Errorf("期望 Killed=8, 得到 %d", stats.Killed)
				}
				if stats.Survived != 2 {
					t.Errorf("期望 Survived=2, 得到 %d", stats.Survived)
				}
			},
		},
		{
			name: "warmup 失败",
			output: `Original test failed (warmup run)
Please check your test suite`,
			wantErr: true,
			validate: func(t *testing.T, stats mutationStats) {
				// 应该返回错误
			},
		},
		{
			name: "变异分数格式",
			output: `Mutation Score: 75.5%
Killed: 15
Survived: 5`,
			wantErr: false,
			validate: func(t *testing.T, stats mutationStats) {
				// 解析分数反推 killed
				if stats.Killed == 0 {
					t.Logf("警告: 未能从分数反推 killed")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats, err := parseMullOutput(tt.output)
			if tt.wantErr && err == "" {
				t.Error("期望错误但未得到")
			}
			if !tt.wantErr && err != "" {
				t.Errorf("不期望错误但得到: %s", err)
			}
			if tt.validate != nil {
				tt.validate(t, stats)
			}
		})
	}
}

// TestMullConfigTemplate 验证 Mull 配置模板
func TestMullConfigTemplate(t *testing.T) {
	// 验证配置模板不包含显式 mutators（使用默认全部）
	if contains := contains(cppMullConfigTemplate, "mutators:"); contains {
		t.Error("❌ 配置模板不应包含显式 mutators，应使用 Mull 默认全部")
	} else {
		t.Log("✅ 配置模板正确使用 Mull 默认变异器")
	}

	// 验证包含必要的排除路径
	requiredPaths := []string{".h$", ".hpp$", "googletest"}
	for _, path := range requiredPaths {
		if !contains(cppMullConfigTemplate, path) {
			t.Errorf("❌ 配置模板缺少排除路径: %s", path)
		}
	}

	t.Logf("📋 Mull 配置模板:\n%s", cppMullConfigTemplate)
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
