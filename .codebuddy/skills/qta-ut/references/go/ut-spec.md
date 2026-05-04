# Go 单元测试规范

## 单元测试编写规范

- 对于每个被测函数仅修改或新增单个测试函数，所有补充的测试逻辑都编写在同一个测试函数中
- 仅编写单元测试代码，不修改被测试代码
- 尽量做最小变更
- 尽量接近已有代码编写风格
- 限制测试代码的副作用：
  - 测试逻辑应该尽量不依赖外部环境，也尽量少对外部环境产生影响
  - 如果测试中涉及读写文件，读写的路径都必须在通过 `t.TempDir()` 创建的临时目录中。一个测试函数中的读写操作尽量复用一个临时目录。
  - 为测试而修改环境变量都必须通过 `t.Setenv` 函数实现
  - 被测逻辑依赖外部 HTTP 服务的，应在测试中通过 `net/http/httptest` 包提供的测试服务替代外部服务接收请求和返回响应
  - 可以使用 `github.com/bytedance/mockey` （参考 `./mockey.md` ）对依赖函数进行 Mock 。若项目中有使用其他类似 Mock 库优先使用项目中已经使用的库
- 断言：测试中可以使用 `github.com/stretchr/testify` 进行断言。若项目中有使用其他类似断言功能的库优先使用项目中已经使用的库
- 尽量复用已有结构体、 Mock 实现等可复用逻辑，减少重复代码

## 测试代码和函数命名规范

- 测试代码路径：应为被测试代码文件名添加 `_test` 后缀，比如被测试代码在 `path/to/file.go` 则测试代码文件路径为 `path/to/file_test.go`
- 测试函数名：
  - 如果被测函数不是成员函数，测试函数名应该是被测函数名加上 `Test` 前缀。比如 `Example` 函数的测试函数命名为 `TestExample`
  - 如果被测函数是成员函数，测试函数名应该是 `Test{reciver}_{funcName}` 格式，其中 `{reciver}` 是首字母转为大写的被测成员函数接收器名， `{funcName}` 是被测函数名。比如 `func (h *handler) example` 函数的测试函数名应为 `TestHandler_example`

## 静态检查运行方法

测试代码编写完成后必须执行静态检查，这是强制步骤，不可跳过

### 按照以下步骤进行静态检查

1. 严格按照以下原则执行静态检查命令
   - 情况A: 如果项目根目录存在 `.golangci.yml` 文件，则执行以下命令进行静态检查:
     `(command -v $(go env GOPATH)/bin/golangci-lint > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest) && $(go env GOPATH)/bin/golangci-lint run '{pkg}'`

   - 情况B: 如果项目根目录不存在 `.golangci.yml` 文件，则执行以下命令进行静态检查:
     `go vet '{pkg}'`
     - 其中 `{pkg}` 是被测试包路径（例如：`./path/to/pkg`）
     - 不要在没有 `.golangci.yml` 文件的项目中强制安装或使用 `golangci-lint`

2. 修复代码规范问题
   如果静态检查发现问题，按照以下步骤修复:
   - 分析错误信息: 分析静态检查结果中错误的原因
   - 问题修复: 根据错误类型进行相应修复
   - 重新检查: 修复后再次运行静态检查
   - 循环修复: 直到所有静态检查都通过

## 单元测试运行方法

- 在项目根目录运行测试，运行测试的命令为 `go test '{pkg}' -run '{testFuncNameRegexp}'` ，其中 `{pkg}` 是被测试包路径， `{testFuncNameRegexp}` 是测试函数名的正则
- 如果需要输出覆盖率数据，在 `go test` 命令中添加 `-coverprofile=/tmp/coverage.out` 参数，其中 `/tmp/coverage.out` 是覆盖率数据输出路径，可以根据情况替换为其它路径但都应该在 `/tmp` 目录下
