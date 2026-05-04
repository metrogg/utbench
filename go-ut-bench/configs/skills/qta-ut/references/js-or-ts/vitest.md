# Vitest 测试框架使用指南

Vitest 是一个由 Vite 驱动的极速单元测试框架，API 与 Jest 高度兼容，但性能更好，原生支持 ESM 和 TypeScript。

## 基本结构

### 测试文件组织

```js
// calculator.test.js
import { describe, it, expect } from 'vitest';
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

Vitest 原生支持 TypeScript，无需额外配置：

```typescript
// calculator.test.ts
import { describe, it, expect } from 'vitest';
import { add, Calculator } from './calculator';

describe('Calculator', () => {
  it('应该正确执行加法运算', () => {
    const calculator = new Calculator();
    const result: number = calculator.add(2, 3);
    expect(result).toBe(5);
  });
});
```

## 断言 (Expect)

Vitest 的断言 API 与 Jest 完全兼容：

### 基本断言

```js
import { expect } from 'vitest';

// 相等性
expect(value).toBe(4); // 严格相等
expect(object).toEqual({ a: 1, b: 2 }); // 深度相等
expect(array).toStrictEqual([1, 2, 3]); // 严格深度相等

// 真假值
expect(value).toBeTruthy();
expect(value).toBeFalsy();
expect(value).toBeNull();
expect(value).toBeUndefined();
expect(value).toBeDefined();

// 数字比较
expect(value).toBeGreaterThan(3);
expect(value).toBeGreaterThanOrEqual(3);
expect(value).toBeLessThan(5);
expect(value).toBeLessThanOrEqual(5);
expect(value).toBeCloseTo(0.3);

// 字符串
expect(string).toMatch(/pattern/);
expect(string).toContain('substring');

// 数组/可迭代对象
expect(array).toContain(item);
expect(array).toHaveLength(3);

// 对象
expect(object).toHaveProperty('key');
expect(object).toHaveProperty('key', value);
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

// 异步函数
await expect(asyncFunction()).rejects.toThrow();
await expect(asyncFunction()).rejects.toThrow('error message');
```

## Mock 功能

### Mock 函数

```js
import { vi } from 'vitest';

// 创建 Mock 函数
const mockFn = vi.fn();

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
import { vi } from 'vitest';

// 完全 Mock 整个模块
vi.mock('./utils');

// 部分 Mock 模块
vi.mock('./utils', () => ({
  ...vi.importActual('./utils'),
  fetchData: vi.fn(),
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

### Spy 函数

```js
import { vi } from 'vitest';

const obj = {
  method: (a, b) => a + b,
};

// 创建 Spy
const spy = vi.spyOn(obj, 'method');

// 调用原始函数
obj.method(1, 2); // 返回 3

// 验证调用
expect(spy).toHaveBeenCalledWith(1, 2);

// 恢复原始实现
spy.mockRestore();
```

### Mock 实现示例

```js
import { vi } from 'vitest';

// Mock HTTP 请求
vi.mock('axios');
import axios from 'axios';

axios.get.mockResolvedValue({
  data: { id: 1, name: 'Alice' },
});

// Mock 定时器
vi.useFakeTimers();
vi.advanceTimersByTime(1000);
vi.runAllTimers();
vi.useRealTimers();

// Mock Date
vi.setSystemTime(new Date('2024-01-01'));
```

## 测试生命周期

```js
import {
  describe,
  it,
  beforeAll,
  afterAll,
  beforeEach,
  afterEach,
} from 'vitest';

describe('Database tests', () => {
  let db;

  // 所有测试前执行一次
  beforeAll(async () => {
    db = await connectDatabase();
  });

  // 所有测试后执行一次
  afterAll(async () => {
    await db.close();
  });

  // 每个测试前执行
  beforeEach(() => {
    db.clearCache();
  });

  // 每个测试后执行
  afterEach(() => {
    vi.clearAllMocks();
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

## 快照测试

```js
import { expect } from 'vitest';

it('应该匹配快照', () => {
  const data = getData();
  expect(data).toMatchSnapshot();
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

Vitest 使用 c8 或 istanbul 收集覆盖率：

```shell
npx vitest run --coverage
```

生成 LCOV 格式覆盖率报告：

```shell
npx vitest run --coverage --coverage.reporter=lcov --coverage.reportsDirectory=/tmp/coverage
```

覆盖率文件将生成在 `/tmp/coverage/lcov.info`。

## 配置 (vitest.config.js)

```js
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    // 测试环境
    environment: 'node', // 或 'jsdom', 'happy-dom'

    // 全局 API
    globals: true, // 启用后无需 import { describe, it, expect }

    // 测试文件匹配
    include: ['**/*.{test,spec}.{js,mjs,cjs,ts,mts,cts,jsx,tsx}'],

    // 覆盖率配置
    coverage: {
      provider: 'c8', // 或 'istanbul'
      reporter: ['text', 'lcov', 'html'],
      reportsDirectory: './coverage',
    },

    // 超时设置
    testTimeout: 10000,
  },
});
```

## Vitest 特有功能

### 1. 并发测试

```js
import { describe, it } from 'vitest';

// 测试套件并发
describe.concurrent('API tests', () => {
  it('测试1', async () => {
    // 这些测试将并发运行
  });

  it('测试2', async () => {
    //
  });
});

// 单个测试并发
it.concurrent('并发测试1', async () => {});
it.concurrent('并发测试2', async () => {});
```

### 2. 测试过滤

```js
// 只运行这个测试
it.only('焦点测试', () => {});

// 跳过这个测试
it.skip('跳过的测试', () => {});

// 待办测试
it.todo('待实现的测试');
```

### 3. 基准测试

```js
import { bench, describe } from 'vitest';

describe('Performance', () => {
  bench('sort', () => {
    const array = Array.from({ length: 1000 }, (_, i) => i);
    array.sort((a, b) => b - a);
  });
});
```

## 最佳实践

### 1. 使用全局 API（可选）

在 `vitest.config.js` 中启用 `globals: true`，则无需每次导入：

```js
// 无需 import { describe, it, expect } from 'vitest';

describe('tests', () => {
  it('works', () => {
    expect(true).toBe(true);
  });
});
```

### 2. Mock 清理

```js
import { vi, afterEach } from 'vitest';

afterEach(() => {
  vi.clearAllMocks(); // 清除 mock 调用记录
  vi.restoreAllMocks(); // 恢复原始实现
});
```

### 3. 测试隔离

```js
describe('Counter', () => {
  let counter;

  beforeEach(() => {
    counter = new Counter();
  });

  it('测试1', () => {
    counter.increment();
    expect(counter.value).toBe(1);
  });

  it('测试2', () => {
    counter.increment();
    expect(counter.value).toBe(1); // 独立运行
  });
});
```

### 4. 使用 Spy 而非完全 Mock

对于只需要监控调用的情况，使用 `vi.spyOn` 而不是完全 Mock：

```js
import { vi } from 'vitest';

// ✅ 好的做法
const spy = vi.spyOn(console, 'log');
myFunction();
expect(spy).toHaveBeenCalled();
spy.mockRestore();

// ❌ 不必要的完全 Mock
vi.mock('console', () => ({
  log: vi.fn(),
}));
```

## 与 Jest 的差异

虽然 API 高度兼容，但有一些细微差异：

### 1. Mock 导入

```js
// Jest
jest.mock('./module');

// Vitest
vi.mock('./module');
```

### 2. Mock 定时器

```js
// Jest
jest.useFakeTimers();
jest.advanceTimersByTime(1000);

// Vitest
vi.useFakeTimers();
vi.advanceTimersByTime(1000);
```

### 3. 配置文件

- Jest: `jest.config.js`
- Vitest: `vitest.config.js` 或集成到 `vite.config.js`

## 常见问题

### 1. 与 Vite 集成

Vitest 可以复用 Vite 配置：

```js
// vite.config.js
import { defineConfig } from 'vite';

export default defineConfig({
  // Vite 配置
  plugins: [],

  // Vitest 配置
  test: {
    environment: 'jsdom',
  },
});
```

### 2. ESM 支持

Vitest 原生支持 ESM，无需额外配置。

### 3. 覆盖率提供者

推荐使用 `c8`（更快）而非 `istanbul`：

```shell
npm install --save-dev @vitest/coverage-c8
```

```js
// vitest.config.js
export default defineConfig({
  test: {
    coverage: {
      provider: 'c8',
    },
  },
});
```
