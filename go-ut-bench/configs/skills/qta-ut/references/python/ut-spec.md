# Python 单元测试规范

## 单元测试编写规范

- 对于每个被测函数仅修改或新增单个测试函数（或测试类中的单个测试方法），所有补充的测试逻辑都编写在同一个测试函数/方法中
- 仅编写单元测试代码，不修改被测试代码
- 尽量做最小变更
- 尽量接近已有代码编写风格
- 测试框架选择：
  - 优先检测项目中是否已使用 `pytest`、`unittest` 或其他测试框架，优先使用项目已有框架
  - 若项目未明确使用测试框架，推荐使用 `pytest`，其次是标准库的 `unittest`
  - 保持与已有测试代码风格一致
- 限制测试代码的副作用：
  - 测试逻辑应该尽量不依赖外部环境，也尽量少对外部环境产生影响
  - 如果测试中涉及读写文件：
    - 使用 `pytest` 时，应使用 `tmp_path` 或 `tmpdir` fixture 提供的临时目录
    - 使用 `unittest` 时，应使用 `tempfile.TemporaryDirectory()` 或 `tempfile.mkdtemp()` 创建临时目录，并在测试结束后清理
  - 为测试而修改环境变量：
    - 使用 `pytest` 时，应使用 `monkeypatch` fixture 的 `setenv` 方法
    - 使用 `unittest` 时，应使用 `unittest.mock.patch.dict('os.environ', {...})` 或在 `setUp` / `tearDown` 中手动保存和恢复
  - 被测逻辑依赖外部 HTTP 服务的，应在测试中通过 mock 库模拟网络请求和响应：
    - 可以使用 `unittest.mock` 模拟 `requests` 等 HTTP 库的调用
    - 或使用 `responses` 、`httpretty` 、`requests-mock` 等第三方库。若项目中已使用某个库优先使用项目中已有的库
- Mock 和断言：
  - Mock ：优先使用标准库 `unittest.mock` （参考 `./unittest-mock.md` ），也可以使用 `pytest-mock` 插件（如果项目已使用）
  - 断言：
    - 使用 `pytest` 时，优先使用 `assert` 语句，也可以使用 `pytest` 的高级断言功能
    - 使用 `unittest` 时，使用 `self.assertEqual`、`self.assertTrue` 等断言方法
    - 若项目中使用了其他断言库（如 `assertpy` 、 `nose.tools` 等），优先使用项目中已有的库

## 测试代码和函数命名规范

- 测试代码路径：
  - 推荐在被测试代码相同目录下创建测试文件，文件名添加 `_test` 后缀
  - 比如被测试代码在 `path/to/module.py`，测试代码文件路径推荐为 `path/to/module_test.py`
  - 也可以遵循项目已有的测试目录结构(如 `tests/` 目录)
- 测试函数命名：
  - 使用 `pytest` 时：
    - 测试函数名应以 `test_` 开头
    - 如果测试的是普通函数，命名为 `test_{func_name}`，如 `test_calculate`
    - 如果测试的是类方法，可以创建测试类 `Test{ClassName}`，在其中创建 `test_{method_name}` 方法
  - 使用 `unittest` 时：
    - 测试类应继承 `unittest.TestCase`，类名以 `Test` 开头，如 `TestCalculator`
    - 测试方法名应以 `test_` 开头，如 `test_add_method`

## 注释规范

```python
def test_xxx_yyy_zzz():
    """
    测试场景：<场景描述>
    前置条件：<前置条件>
    输入数据：<输入描述>
    预期结果：<预期输出>
    """
```

## 参数化测试（按需）

仅在多用例共享相同执行流程时使用：

```python
@pytest.mark.parametrize("input,expected", [
    ("", True),
    ("valid", False),
])
def test_validate_input(input, expected):
    err = validate(input)
    assert (err is not None) == expected
```

## 异步测试

```python
@pytest.mark.asyncio  # 必须添加
async def test_async_method():
    # Arrange
    mock_repo = Mock(spec=AsyncRepo)
    mock_repo.fetch = AsyncMock(return_value=data)

    # Act
    result = await service.fetch_data()

    # Assert
    assert result is not None
    mock_repo.fetch.assert_called_once()
```

## 单元测试运行方法

- 在项目根目录运行测试
- 使用 `pytest` 时：
  - 运行单个测试函数： `pytest {test_file_path}::{test_func_name} -v`
  - 运行测试文件中所有测试： `pytest {test_file_path} -v`
  - 如果需要输出覆盖率数据，添加参数： `pytest {test_file_path} --cov={module_path} --cov-report=xml:/tmp/coverage.xml`
- 使用 `unittest` 时：
  - 运行单个测试： `python -m unittest {module_path}.{test_class_name}.{test_method_name}`
  - 运行测试文件： `python -m unittest {module_path}.{test_class_name}`
  - 如果需要输出覆盖率数据，使用 `coverage` 工具： `coverage run -m unittest {module_path} && coverage xml -o /tmp/coverage.xml`
- 其中：
  - `{test_file_path}` 是测试文件路径
  - `{test_func_name}` 是测试函数名
  - `{module_path}` 是模块的导入路径(如 `path.to.module`)
  - `{test_class_name}` 是测试类名
  - `{test_method_name}` 是测试方法名
  - 覆盖率数据输出路径应在 `/tmp` 目录下
