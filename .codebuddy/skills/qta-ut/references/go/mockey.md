# Go Mockey Mock 框架参考指南

⚠️ **仅在测试涉及外部依赖（数据库、API、文件系统等）时使用 Mock**

## 一、安装与配置

```shell
go get github.com/bytedance/mockey@latest
```

```go
import "github.com/bytedance/mockey"
```

⚠️ **必须禁用内联优化**: `go test -gcflags=all=-l -v ./...`

## 二、使用场景

- ✅ 数据库操作 → Mock `database.Query()`
- ✅ HTTP/RPC 调用 → Mock `http.Get()`, `client.Call()`
- ✅ 时间函数 → Mock `time.Now()`
- ✅ 文件操作 → Mock `os.ReadFile()`
- ❌ 纯函数逻辑、数据转换 → 不需要 Mock

## 三、核心语法

### 3.1 Mock 函数

```go
// Mock 包级别函数
patch := mockey.Mock(database.QueryUser).To(func(ctx context.Context, id int) (*User, error) {
    return &User{ID: id, Name: "Test"}, nil
}).Build()
defer patch.UnPatch()
```

### 3.2 Mock 方法

```go
// Mock 结构体方法
patch := mockey.Mock((*sql.DB).Exec).To(func(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
    return &mockResult{rowsAffected: 1}, nil
}).Build()
defer patch.UnPatch()
```

### 3.3 Mock 变量

```go
// Mock 全局变量
patch := mockey.MockValue(&GlobalVar).To(123).Build()
defer patch.UnPatch()
```

### 3.4 条件 Mock

```go
// 根据参数返回不同结果
patch := mockey.Mock(database.Query).To(func(ctx context.Context, sql string) ([]Row, error) {
    if strings.Contains(sql, "user") {
        return []Row{{ID: 1}}, nil
    }
    return nil, errors.New("unknown")
}).Build()
defer patch.UnPatch()
```

## 四、Mockey API 使用指南

### 4.1 基础用法

```go
// ✅ 正确：Mock函数
patch := mockey.Mock(os.ReadFile).Return([]byte("data"), nil).Build()
defer patch.UnPatch()

// ✅ 正确：Mock方法（值接收器）
patch := mockey.Mock(A.Method).Return("result").Build()
defer patch.UnPatch()

// ✅ 正确：Mock方法（指针接收器）
patch := mockey.Mock((*B).Method).Return("result").Build()
defer patch.UnPatch()

// ✅ 正确：Mock变量
patch := mockey.MockValue(&GlobalVar).To(123).Build()
defer patch.UnPatch()
```

### 4.2 高级用法

```go
// ✅ 序列返回
patch := mockey.Mock(apiFunction).Return(
    mockey.Sequence("result1").Times(2).Then("result2").Build(),
).Build()
defer patch.UnPatch()

// ✅ 条件Mock
patch := mockey.Mock((*Service).Method).
    When(func(ctx context.Context, req *Request) bool {
        return req.ID == "specific-id"
    }).
    Return(&Response{Data: "mocked"}).
    Build()
defer patch.UnPatch()

// ✅ 装饰器Mock
patch := mockey.Mock(originalFunc).To(func(req *Request) *Response {
    // 先调用原函数
    resp := originalFunc(req)
    // 再修改返回值
    resp.Header = "modified"
    return resp
}).Origin(&originalFunc).Build()
defer patch.UnPatch()
```

### 4.3 清理机制

```go
// ✅ 必须立即defer清理
patch := mockey.Mock(Function).Return(result).Build()
defer patch.UnPatch()  // 立即执行，防止内存泄漏

// ✅ 使用t.Cleanup替代defer（推荐）
patch := mockey.Mock(Function).Return(result).Build()
t.Cleanup(patch.UnPatch)  // 更安全，测试失败时也会清理
```

### 4.4 调用计数

```go
patch := mockey.Mock(Function).Return(result).Build()
defer patch.UnPatch()

// 检查调用次数
assert.Equal(t, 1, patch.Times())      // 总调用次数
assert.Equal(t, 1, patch.MockTimes())   // 被Mock的调用次数
```

## 五、最佳实践

### 5.1 必须立即 defer

```go
patch := mockey.Mock(Func).To(mockFunc).Build()
defer patch.UnPatch()  // ✅ 防止 panic 导致未清理
```

### 5.2 Mock 原则

- 只 Mock 外部依赖，不 Mock 内部逻辑
- Mock 返回值应符合真实场景
- 使用 `t.Run()` 隔离多场景测试

```go
t.Run("success", func(t *testing.T) {
    patch := mockey.Mock(Func).To(mockSuccess).Build()
    defer patch.UnPatch()
    // 测试...
})

t.Run("error", func(t *testing.T) {
    patch := mockey.Mock(Func).To(mockError).Build()
    defer patch.UnPatch()
    // 测试...
})
```

## 六、常见问题

### Mock 不生效

确保使用 `-gcflags=all=-l` 参数禁用内联优化：

```shell
go test -gcflags=all=-l -v ./...
```

### 并发测试注意

避免在 `t.Parallel()` 测试中使用全局 Mock，或确保每个子测试独立创建和清理 Mock。

### 接口类型Mock错误

```go
// ❌ 错误：尝试直接mock接口
patch := mockey.Mock((*service.Repository).Method).To(...)

// ✅ 正确方式1：创建mock实现结构体
type mockRepository struct {
    result *Model
    err    error
}
func (m *mockRepository) Method(...) (*Model, error) {
    return m.result, m.err
}
// 在测试中注入mock实例

// ✅ 正确方式2：Mock接口的具体实现类型
patch := mockey.Mock((*ConcreteImpl).Method).To(...)
```

### Context 使用规范

```go
// ❌ 错误：使用nil context（会导致static check警告）
result, err := client.GetData(nil, "id")

// ✅ 正确：使用context.TODO()或context.Background()
result, err := client.GetData(context.TODO(), "id")

// ✅ 或者在测试中创建特定context
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
result, err := client.GetData(ctx, "id")
```

### 错误类型断言问题

```go
// 先读取项目错误定义，区分错误常量和错误类型

// ❌ 错误：将错误常量当作类型使用
var valErr *errors.ValidationError  // ValidationError是string常量

// ✅ 正确：使用实际的错误类型
var appErr *errors.AppError
assert.True(t, errors.As(err, &appErr))
assert.Equal(t, errors.ValidationError, appErr.Code)

// ✅ 或使用 errors.Is 判断
assert.True(t, errors.Is(err, errors.ErrValidation))
```

### 字段访问路径错误

```go
// ❌ 错误：未先读取结构体定义，猜测字段路径
cfg.QueueConfig.Timeout = 1

// ✅ 正确：先读取Config结构体定义，确认实际字段路径
// 假设实际结构为：
// type QueueManager struct {
//     queueConfig *QueueConfig  // 小写开头，私有字段
// }
// 需要通过其他方式设置，如构造函数参数或公开方法
```

### 重复定义检查

```shell
# 在添加新mock结构体前，先搜索是否已存在
grep "type mockXxx struct" *_test.go

# 避免重复定义导致编译错误
```
