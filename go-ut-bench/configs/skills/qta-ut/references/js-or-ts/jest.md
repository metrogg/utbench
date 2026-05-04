# Jest 测试框架使用指南

Jest 是最流行的 JavaScript 测试框架，内置了测试运行器、断言库、Mock 功能和覆盖率报告。

## 基本结构

### 测试文件组织

```js
// calculator.test.js
import { add, subtract, multiply } from './calculator';

describe('Calculator functions', () => {
  describe('add', () => {
    it('应该正确相加两个正数', () => {
      expect(add(2, 3)).toBe(5);
    });

    it('应该正确处理负数', () => {
      expect(add(-1, -2)).toBe(-3);
    });
  });

  describe('subtract', () => {
    it('应该正确相减两个数', () => {
      expect(subtract(5, 3)).toBe(2);
    });
  });
});
```

### TypeScript 支持

```typescript
// calculator.test.ts
import { add, Calculator } from './calculator';

describe('Calculator', () => {
  let calculator: Calculator;

  beforeEach(() => {
    calculator = new Calculator();
  });

  it('应该正确执行加法运算', () => {
    const result: number = calculator.add(2, 3);
    expect(result).toBe(5);
  });
});
```

## 断言 (Expect)

### 基本断言

```js
// 相等性
expect(value).toBe(4); // 严格相等 (===)
expect(object).toEqual({ a: 1, b: 2 }); // 深度相等
expect(array).toStrictEqual([1, 2, 3]); // 严格深度相等

// 真假值
expect(value).toBeTruthy(); // 真值
expect(value).toBeFalsy(); // 假值
expect(value).toBeNull(); // null
expect(value).toBeUndefined(); // undefined
expect(value).toBeDefined(); // 已定义

// 数字比较
expect(value).toBeGreaterThan(3); // 大于
expect(value).toBeGreaterThanOrEqual(3); // 大于等于
expect(value).toBeLessThan(5); // 小于
expect(value).toBeLessThanOrEqual(5); // 小于等于
expect(value).toBeCloseTo(0.3); // 浮点数近似相等

// 字符串
expect(string).toMatch(/pattern/); // 正则匹配
expect(string).toContain('substring'); // 包含子串

// 数组/可迭代对象
expect(array).toContain(item); // 包含元素
expect(array).toHaveLength(3); // 长度检查

// 对象
expect(object).toHaveProperty('key'); // 有属性
expect(object).toHaveProperty('key', value); // 属性值检查
```

### 异常断言

```js
// 同步函数
expect(() => {
  throw new Error('error');
}).toThrow();

expect(() => {
  throw new Error('error message');
}).toThrow('error message');

expect(() => {
  throw new TypeError('type error');
}).toThrow(TypeError);

// 异步函数
await expect(asyncFunction()).rejects.toThrow();
await expect(asyncFunction()).rejects.toThrow('error message');
```

## Mock 功能

### Mock 函数

```js
// 创建 Mock 函数
const mockFn = jest.fn();

// 设置返回值
mockFn.mockReturnValue(42);
mockFn.mockReturnValueOnce(1).mockReturnValueOnce(2);

// 设置异步返回值
mockFn.mockResolvedValue('success');
mockFn.mockRejectedValue(new Error('failed'));

// 设置实现
mockFn.mockImplementation((a, b) => a + b);

// 验证调用
expect(mockFn).toHaveBeenCalled();
expect(mockFn).toHaveBeenCalledTimes(2);
expect(mockFn).toHaveBeenCalledWith(arg1, arg2);
expect(mockFn).toHaveBeenLastCalledWith(arg1, arg2);
```

### Mock 模块

```js
// 完全 Mock 整个模块
jest.mock('./utils');

// 部分 Mock 模块
jest.mock('./utils', () => ({
  ...jest.requireActual('./utils'),
  fetchData: jest.fn(),
}));

// 使用 Mock 模块
import { fetchData } from './utils';

describe('API tests', () => {
  it('应该成功获取数据', async () => {
    fetchData.mockResolvedValue({ data: 'test' });

    const result = await fetchData('api/users');

    expect(result.data).toBe('test');
    expect(fetchData).toHaveBeenCalledWith('api/users');
  });
});
```

### Mock 实现示例

```js
// Mock HTTP 请求
jest.mock('axios');
import axios from 'axios';

axios.get.mockResolvedValue({
  data: { id: 1, name: 'Alice' },
});

// Mock 文件系统
jest.mock('fs');
import fs from 'fs';

fs.readFileSync.mockReturnValue('file content');

// Mock 定时器
jest.useFakeTimers();
jest.advanceTimersByTime(1000);
jest.runAllTimers();
jest.useRealTimers();
```

## 测试生命周期

### Setup 和 Teardown

```js
describe('Database tests', () => {
  let db;

  // 每个测试文件执行一次
  beforeAll(async () => {
    db = await connectDatabase();
  });

  afterAll(async () => {
    await db.close();
  });

  // 每个测试用例执行一次
  beforeEach(() => {
    db.clearCache();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('测试1', () => {
    // 测试代码
  });

  it('测试2', () => {
    // 测试代码
  });
});
```

## 异步测试

### async/await

```js
it('应该异步获取用户数据', async () => {
  const user = await fetchUser(123);
  expect(user.name).toBe('Alice');
});
```

### Promise

```js
it('应该异步获取用户数据', () => {
  return fetchUser(123).then((user) => {
    expect(user.name).toBe('Alice');
  });
});
```

### 回调函数

```js
it('应该通过回调返回数据', (done) => {
  fetchUser(123, (error, user) => {
    expect(error).toBeNull();
    expect(user.name).toBe('Alice');
    done();
  });
});
```

## 快照测试

```js
it('应该匹配快照', () => {
  const component = render(<MyComponent />);
  expect(component).toMatchSnapshot();
});

// 行内快照
it('应该返回正确的对象', () => {
  expect(getData()).toMatchInlineSnapshot(`
    {
      "id": 1,
      "name": "test"
    }
  `);
});
```

## 测试覆盖率

Jest 内置覆盖率支持，运行测试时添加 `--coverage` 选项：

```shell
npx jest --coverage
```

生成 LCOV 格式覆盖率报告：

```shell
npx jest --coverage --coverageReporters=lcov --coverageDirectory=/tmp/coverage
```

覆盖率文件将生成在 `/tmp/coverage/lcov.info`。

## 配置 (jest.config.js)

```js
module.exports = {
  // 测试环境
  testEnvironment: 'node', // 或 'jsdom' 用于浏览器环境

  // 测试文件匹配模式
  testMatch: ['**/__tests__/**/*.js', '**/?(*.)+(spec|test).js'],

  // 覆盖率收集
  collectCoverageFrom: [
    'src/**/*.{js,jsx,ts,tsx}',
    '!src/**/*.test.{js,jsx,ts,tsx}',
  ],

  // TypeScript 支持
  transform: {
    '^.+\\.tsx?$': 'ts-jest',
  },

  // 模块路径别名
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1',
  },
};
```

## 最佳实践

### 1. 测试结构清晰

使用 AAA (Arrange-Act-Assert) 模式：

```js
it('应该正确计算总价', () => {
  // Arrange - 准备测试数据
  const items = [
    { price: 10, quantity: 2 },
    { price: 5, quantity: 3 },
  ];

  // Act - 执行被测函数
  const total = calculateTotal(items);

  // Assert - 验证结果
  expect(total).toBe(35);
});
```

### 2. 测试隔离

每个测试应该独立，不依赖其他测试的执行结果：

```js
describe('Counter', () => {
  let counter;

  beforeEach(() => {
    counter = new Counter(); // 每个测试都创建新实例
  });

  it('测试1', () => {
    counter.increment();
    expect(counter.value).toBe(1);
  });

  it('测试2', () => {
    counter.increment();
    expect(counter.value).toBe(1); // 不受测试1影响
  });
});
```

### 3. Mock 清理

在 `afterEach` 中清理 Mock：

```js
afterEach(() => {
  jest.clearAllMocks(); // 清除 mock 调用记录
  jest.restoreAllMocks(); // 恢复原始实现
});
```

### 4. 避免过度 Mock

只 Mock 必要的外部依赖，不要 Mock 被测函数内部的逻辑：

```js
// ✅ 好的做法 - Mock 外部 API
jest.mock('axios');
axios.get.mockResolvedValue({ data: 'test' });

// ❌ 不好的做法 - Mock 被测函数本身
jest.mock('./calculator');
```

### 5. 描述性测试名称

测试名称应该清晰描述测试的场景和预期结果：

```js
// ✅ 好的测试名称
it('当输入为空数组时应该返回 0', () => {});
it('当价格为负数时应该抛出错误', () => {});

// ❌ 不好的测试名称
it('测试1', () => {});
it('works', () => {});
```

## 常见问题

### 1. ESM 模块支持

如果项目使用 ESM，在 `package.json` 中添加：

```json
{
  "type": "module"
}
```

并使用 `.mjs` 扩展名或配置 Jest：

```js
// jest.config.js
export default {
  transform: {},
  testMatch: ['**/*.test.mjs'],
};
```

### 2. TypeScript 支持

安装 `ts-jest` 和配置：

```shell
npm install --save-dev ts-jest @types/jest
```

```js
// jest.config.js
module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
};
```

### 3. Mock 不生效

确保 `jest.mock()` 在 import 语句之前调用（会被提升）：

```js
jest.mock('./utils'); // 正确位置

import { fetchData } from './utils';
```
