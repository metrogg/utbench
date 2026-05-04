# JavaScript/TypeScript 单元测试规范

## 关于读取代码的要求

- 编写单元测试不需要读取整份代码文件，只需要读取被测函数，因此应该通过 Grep 查找函数定义所在行，然后从函数定义所在行往下每次读取 200 行，直到读取到函数末尾
- JavaScript/TypeScript 函数定义形式多样（function 声明、箭头函数、类方法、对象方法等），需要正确识别

## 单元测试编写规范

### 基本原则

- 对于每个被测函数仅修改或新增单个测试函数或测试块，所有补充的测试逻辑都编写在同一个测试函数/块中
- 仅编写单元测试代码，不修改被测试代码
- 尽量做最小变更
- 尽量接近已有代码编写风格
- 限制测试代码的副作用：
  - 测试逻辑应该尽量不依赖外部环境，也尽量少对外部环境产生影响
  - 如果测试中涉及读写文件，应使用临时目录并在测试结束后清理
  - 为测试而修改环境变量应使用测试框架提供的机制（如 Jest 的 process.env 设置）
  - 被测逻辑依赖外部 HTTP 服务的，应在测试中通过 Mock 库模拟网络请求和响应

### 测试框架选择

优先使用项目中已经使用的测试框架。识别方法：检查 `package.json` 的 `dependencies` 或 `devDependencies` 字段。

- **Jest**（推荐，优先级最高）：最流行的 JavaScript 测试框架，内置 Mock、断言、覆盖率支持。具体规范参考 `./jest.md`
- **Vitest**：现代化的 Vite 原生测试框架，API 与 Jest 高度兼容。具体规范参考 `./vitest.md`
- **Mocha + Chai**：经典的测试框架组合，需要额外配置断言库和 Mock 库。具体规范参考 `./mocha.md`

如果项目未明确使用任何测试框架，推荐使用 Jest。

### 模块系统识别

JavaScript/TypeScript 项目可能使用 ESM（ES Modules）或 CommonJS 模块系统，需要正确识别并生成相应的导入语句。

**识别规则**：

1. 检查 `package.json` 中的 `"type"` 字段：
   - `"type": "module"` → 使用 ESM
   - 缺失或 `"type": "commonjs"` → 使用 CommonJS
2. 文件扩展名优先级：
   - `.mjs`, `.mts` → 强制使用 ESM
   - `.cjs`, `.cts` → 强制使用 CommonJS
   - `.js`, `.ts`, `.jsx`, `.tsx` → 根据 package.json 决定

**ESM 导入示例**：

```js
import { calculateTotal } from './calculator.js';
import { describe, it, expect } from '@jest/globals';
```

**CommonJS 导入示例**：

```js
const { calculateTotal } = require('./calculator');
const { describe, it, expect } = require('@jest/globals');
```

### TypeScript 类型处理

对于 TypeScript 文件（`.ts`, `.tsx`），生成的测试代码应当保留类型信息：

- 测试文件使用 `.ts` 或 `.tsx` 扩展名
- 在测试数据中使用明确的类型注解
- 导入并使用被测代码中定义的类型、接口
- 确保生成的代码能通过 TypeScript 类型检查

示例：

```typescript
import { User, createUser } from './user';

describe('createUser', () => {
  it('应该创建有效的用户对象', () => {
    const userData: Partial<User> = {
      name: 'Alice',
      email: 'alice@example.com',
    };

    const user: User = createUser(userData);

    expect(user.name).toBe('Alice');
    expect(user.email).toBe('alice@example.com');
  });
});
```

## 测试代码和函数命名规范

### 测试文件路径

优先遵循项目已有的测试文件组织方式。常见模式：

1. **同目录模式**（推荐）：
   - 被测文件：`src/utils/calculator.js`
   - 测试文件：`src/utils/calculator.test.js` 或 `src/utils/calculator.spec.js`

2. **测试目录模式**：
   - 被测文件：`src/utils/calculator.js`
   - 测试文件：`src/utils/__tests__/calculator.test.js`

3. **独立测试目录**：
   - 被测文件：`src/utils/calculator.js`
   - 测试文件：`tests/utils/calculator.test.js`

检查项目中是否存在 `.test.js`、`.spec.js` 或 `__tests__` 目录，优先使用项目已有的模式。

### 测试函数命名

- 使用 `describe` 块组织测试套件，描述被测函数或模块
- 使用 `it` 或 `test` 块编写具体测试用例，描述应清晰表达测试意图
- 测试函数名应以 `test_` 开头（仅在必要时），或使用描述性文本

示例：

```js
describe('calculateTotal', () => {
  it('应该正确计算商品总价', () => {
    // 测试代码
  });

  it('应该在输入为空时返回 0', () => {
    // 测试代码
  });

  it('应该正确处理负数价格', () => {
    // 测试代码
  });
});
```

## 运行静态检查方法

**⚠️ 这是可选步骤，仅在项目配置了 ESLint 时执行**

### 检查是否需要运行 ESLint

检查项目根目录是否存在以下配置文件之一：

- `.eslintrc.js`
- `.eslintrc.json`
- `.eslintrc.yml`
- `.eslintrc.yaml`
- `eslint.config.js`
- `package.json` 中的 `eslintConfig` 字段

### 执行 ESLint 检查

如果存在 ESLint 配置，运行以下命令：

```shell
npx eslint {test_file_path}
```

如果存在可自动修复的问题，可以尝试：

```shell
npx eslint {test_file_path} --fix
```

### 处理检查结果

1. 如果没有错误，继续执行后续步骤
2. 如果存在错误，分析错误类型：
   - 缺少导入语句 → 自动添加缺失的导入
   - 代码风格问题 → 根据规则调整代码风格
   - 未使用的变量 → 移除或添加 `// eslint-disable-next-line` 注释
3. 修复后重新运行检查，最多重试 2 次
4. 如果无法自动修复，在输出中说明检查失败，但继续执行测试运行步骤

**注意**：不要在没有 ESLint 配置的项目中强制安装或运行 ESLint。

## 单元测试运行方法

根据检测到的测试框架，使用相应的命令运行测试。

### Jest

运行单个测试文件：

```shell
npx jest {test_file_path}
```

运行特定测试用例：

```shell
npx jest {test_file_path} -t "{test_name_pattern}"
```

生成覆盖率数据：

```shell
npx jest {test_file_path} --coverage --coverageReporters=lcov --coverageDirectory=/tmp/coverage
```

覆盖率文件位置：`/tmp/coverage/lcov.info`

### Vitest

运行单个测试文件：

```shell
npx vitest run {test_file_path}
```

运行特定测试用例：

```shell
npx vitest run {test_file_path} -t "{test_name_pattern}"
```

生成覆盖率数据：

```shell
npx vitest run {test_file_path} --coverage --coverage.reporter=lcov --coverage.reportsDirectory=/tmp/coverage
```

覆盖率文件位置：`/tmp/coverage/lcov.info`

### Mocha + Chai

运行单个测试文件：

```shell
npx mocha {test_file_path}
```

运行特定测试用例：

```shell
npx mocha {test_file_path} --grep "{test_name_pattern}"
```

生成覆盖率数据（需要 NYC）：

```shell
npx nyc --reporter=lcov --report-dir=/tmp/coverage mocha {test_file_path}
```

覆盖率文件位置：`/tmp/coverage/lcov.info`

### 提交覆盖率到 MCP

测试运行成功后，调用 MCP `qta_ut_agent_coverage` 服务提交覆盖率数据：

- 方法：`push_raw`
- 参数：
  - `path`: `/tmp/coverage/lcov.info`
  - `format`: `lcov`
  - `commitRef`: 当前 Git commit 或分支名（可选）

提交后查询被测函数的覆盖率变化，对比测试前后的行覆盖率百分比。

## 异步函数测试

JavaScript/TypeScript 中的异步函数需要特殊处理：

### async/await 函数

```js
describe('fetchUserData', () => {
  it('应该成功获取用户数据', async () => {
    const userId = 123;

    const userData = await fetchUserData(userId);

    expect(userData).toBeDefined();
    expect(userData.id).toBe(userId);
  });
});
```

### Promise 返回的函数

```js
describe('fetchUserData', () => {
  it('应该成功获取用户数据', () => {
    const userId = 123;

    return fetchUserData(userId).then((userData) => {
      expect(userData).toBeDefined();
      expect(userData.id).toBe(userId);
    });
  });
});
```

### 异步错误处理

```js
describe('fetchUserData', () => {
  it('应该在用户不存在时抛出错误', async () => {
    const userId = -1;

    await expect(fetchUserData(userId)).rejects.toThrow('User not found');
  });
});
```

## Mock 使用指南

不同测试框架的 Mock 使用方式不同，详见各框架的具体文档：

- Jest Mock：参考 `./jest.md`
- Vitest Mock：参考 `./vitest.md`（与 Jest 类似）
- Sinon.js（用于 Mocha）：参考 `./mocha.md`

基本原则：

1. Mock 外部依赖函数和模块
2. Mock HTTP 请求和响应
3. Mock 文件系统操作
4. Mock 时间相关函数（setTimeout, Date.now 等）
5. 测试结束后恢复 Mock（避免影响其他测试）
