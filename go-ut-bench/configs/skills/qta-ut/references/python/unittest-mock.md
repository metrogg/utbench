# Python unittest.mock 框架参考指南

⚠️ **仅在测试涉及外部依赖（数据库、API、文件系统等）时使用 Mock**

## 一、导入与配置

```python
from unittest.mock import Mock, AsyncMock, patch, MagicMock
```

## 二、使用场景

- ✅ 数据库操作 → Mock `db.query()`
- ✅ HTTP/RPC 调用 → Mock `requests.get()`, `httpx.Client()`
- ✅ 时间函数 → Mock `time.time()`, `datetime.now()`
- ✅ 文件操作 → Mock `open()`, `os.path.exists()`
- ❌ 纯函数逻辑、数据转换 → 不需要 Mock

## 三、核心语法

### 3.1 Mock 函数/方法

```python
# Mock 函数
@patch('module.function_name')
def test_something(mock_func):
    mock_func.return_value = "mocked_result"
    result = function_name()
    assert result == "mocked_result"
```

### 3.2 Mock 对象方法

```python
# 使用 spec 确保类型安全
mock_repo = Mock(spec=UserRepository)
mock_repo.find_by_id.return_value = User(id=1, name="test")

result = mock_repo.find_by_id(1)
assert result.id == 1
```

### 3.3 Mock 异步方法

```python
# 必须使用 AsyncMock
mock_repo = Mock(spec=AsyncRepository)
mock_repo.fetch_data = AsyncMock(return_value=[1, 2, 3])

result = await mock_repo.fetch_data()
assert result == [1, 2, 3]
```

### 3.4 条件 Mock

```python
# 根据参数返回不同结果
def mock_query(sql: str):
    if "users" in sql:
        return [{"id": 1, "name": "user1"}]
    return []

mock_db.query.side_effect = mock_query
```

## 四、最佳实践

### 4.1 必须使用 spec 参数

```python
# ✅ 推荐：防止拼写错误
mock_service = Mock(spec=UserService)

# ❌ 避免：任何方法都不会报错
mock_service = Mock()
```
