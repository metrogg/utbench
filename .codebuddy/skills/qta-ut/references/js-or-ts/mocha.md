# Mocha + Chai + Sinon 测试框架使用指南

Mocha 是一个灵活的 JavaScript 测试框架，通常与 Chai（断言库）和 Sinon（Mock/Spy 库）配合使用。

## 基本结构

### 测试文件组织

```js
// calculator.test.js
const { expect } = require('chai');
const { add, subtract } = require('./calculator');

describe('Calculator functions', () => {
  describe('add', () => {
    it('应该正确相加两个正数', () => {
      expect(add(2, 3)).to.equal(5);
    });

    it('应该正确处理负数', () => {
      expect(add(-1, -2)).to.equal(-3);
    });
  });

  describe('subtract', () => {
    it('应该正确相减两个数', () => {
      expect(subtract(5, 3)).to.equal(2);
    });
  });
});
```

### ESM 支持

```js
// calculator.test.mjs
import { expect } from 'chai';
import { add, subtract } from './calculator.js';

describe('Calculator functions', () => {
  it('应该正确相加两个数', () => {
    expect(add(2, 3)).to.equal(5);
  });
});
```

### TypeScript 支持

```typescript
// calculator.test.ts
import { expect } from 'chai';
import { add, Calculator } from './calculator';

describe('Calculator', () => {
  let calculator: Calculator;

  beforeEach(() => {
    calculator = new Calculator();
  });

  it('应该正确执行加法运算', () => {
    const result: number = calculator.add(2, 3);
    expect(result).to.equal(5);
  });
});
```

## Chai 断言库

Chai 提供三种断言风格：expect、should、assert。推荐使用 expect 风格。

### BDD 风格 (expect)

```js
const { expect } = require('chai');

// 相等性
expect(value).to.equal(4); // 严格相等
expect(object).to.deep.equal({ a: 1, b: 2 }); // 深度相等
expect(array).to.have.ordered.members([1, 2, 3]); // 数组相等

// 真假值
expect(value).to.be.true;
expect(value).to.be.false;
expect(value).to.be.null;
expect(value).to.be.undefined;
expect(value).to.exist;

// 类型检查
expect(value).to.be.a('string');
expect(value).to.be.an('array');
expect(value).to.be.instanceof(MyClass);

// 数字比较
expect(value).to.be.above(3); // 大于
expect(value).to.be.at.least(3); // 大于等于
expect(value).to.be.below(5); // 小于
expect(value).to.be.at.most(5); // 小于等于
expect(value).to.be.closeTo(0.3, 0.01); // 近似相等

// 字符串
expect(string).to.match(/pattern/); // 正则匹配
expect(string).to.include('substring'); // 包含子串

// 数组
expect(array).to.include(item); // 包含元素
expect(array).to.have.lengthOf(3); // 长度检查
expect(array).to.be.empty; // 空数组

// 对象
expect(object).to.have.property('key'); // 有属性
expect(object).to.have.property('key', value); // 属性值检查
expect(object).to.have.all.keys('a', 'b'); // 包含所有键
```

### 异常断言

```js
// 同步函数
expect(() => {
  throw new Error('error');
}).to.throw();

expect(() => {
  throw new Error('error message');
}).to.throw('error message');

expect(() => {
  throw new TypeError('type error');
}).to.throw(TypeError);

// 异步函数
await expect(asyncFunction()).to.be.rejected;
await expect(asyncFunction()).to.be.rejectedWith('error message');
await expect(asyncFunction()).to.be.rejectedWith(Error);
```

### 否定断言

```js
expect(value).to.not.equal(4);
expect(value).to.not.be.null;
expect(array).to.not.include(item);
```

## Sinon Mock 和 Spy

### Spy 函数

Spy 用于监控函数调用，但不改变函数行为：

```js
const sinon = require('sinon');

// 创建独立 Spy
const spy = sinon.spy();

// Spy 对象方法
const obj = {
  method: (a, b) => a + b,
};
const spy = sinon.spy(obj, 'method');

obj.method(1, 2); // 正常调用

// 验证调用
expect(spy.called).to.be.true;
expect(spy.callCount).to.equal(1);
expect(spy.calledWith(1, 2)).to.be.true;
expect(spy.calledOnce).to.be.true;

// 恢复原始函数
spy.restore();
```

### Stub 函数

Stub 用于替换函数行为并控制返回值：

```js
const sinon = require('sinon');

// 创建独立 Stub
const stub = sinon.stub();

// Stub 对象方法
const obj = {
  method: (a, b) => a + b,
};
const stub = sinon.stub(obj, 'method');

// 设置返回值
stub.returns(42);
stub.onCall(0).returns(1);
stub.onCall(1).returns(2);

// 设置异步返回值
stub.resolves('success');
stub.rejects(new Error('failed'));

// 设置不同参数的返回值
stub.withArgs(1, 2).returns(10);
stub.withArgs(3, 4).returns(20);

// 调用
const result = obj.method(1, 2); // 返回 10

// 验证
expect(stub.calledWith(1, 2)).to.be.true;

// 恢复
stub.restore();
```

### Mock 对象

Mock 用于预设期望并验证：

```js
const sinon = require('sinon');

const obj = {
  method: () => {},
};

// 创建 Mock
const mock = sinon.mock(obj);

// 设置期望
mock.expects('method').once().withArgs(1, 2).returns(42);

// 调用
obj.method(1, 2);

// 验证所有期望
mock.verify(); // 如果期望未满足会抛出错误

// 恢复
mock.restore();
```

### Mock 模块

```js
const sinon = require('sinon');

// Mock 整个模块
const axios = require('axios');
sinon.stub(axios, 'get').resolves({ data: 'test' });

// 使用
const response = await axios.get('http://api.example.com');
expect(response.data).to.equal('test');

// 恢复
axios.get.restore();
```

### Fake 定时器

```js
const sinon = require('sinon');

describe('Timer tests', () => {
  let clock;

  beforeEach(() => {
    clock = sinon.useFakeTimers();
  });

  afterEach(() => {
    clock.restore();
  });

  it('应该在1秒后执行回调', () => {
    const callback = sinon.spy();

    setTimeout(callback, 1000);

    // 快进1秒
    clock.tick(1000);

    expect(callback.calledOnce).to.be.true;
  });
});
```

## 测试生命周期

```js
describe('Database tests', () => {
  let db;

  // 所有测试前执行一次
  before(async () => {
    db = await connectDatabase();
  });

  // 所有测试后执行一次
  after(async () => {
    await db.close();
  });

  // 每个测试前执行
  beforeEach(() => {
    db.clearCache();
  });

  // 每个测试后执行
  afterEach(() => {
    sinon.restore(); // 恢复所有 Sinon stub/spy
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
  expect(user.name).to.equal('Alice');
});
```

### Promise

```js
it('应该异步获取用户数据', () => {
  return fetchUser(123).then((user) => {
    expect(user.name).to.equal('Alice');
  });
});
```

### 回调函数

```js
it('应该通过回调返回数据', (done) => {
  fetchUser(123, (error, user) => {
    expect(error).to.be.null;
    expect(user.name).to.equal('Alice');
    done();
  });
});
```

**注意**：如果使用回调，必须调用 `done()`，否则测试会超时。

## 测试覆盖率

Mocha 本身不提供覆盖率功能，需要使用 NYC (Istanbul)：

### 安装

```shell
npm install --save-dev nyc
```

### 运行测试并生成覆盖率

```shell
npx nyc mocha
```

生成 LCOV 格式报告：

```shell
npx nyc --reporter=lcov --report-dir=/tmp/coverage mocha
```

覆盖率文件将生成在 `/tmp/coverage/lcov.info`。

### NYC 配置 (.nycrc.json)

```json
{
  "all": true,
  "include": ["src/**/*.js"],
  "exclude": ["**/*.test.js", "**/*.spec.js"],
  "reporter": ["text", "lcov", "html"],
  "report-dir": "./coverage"
}
```

## Mocha 配置

### 配置文件 (.mocharc.json)

```json
{
  "require": ["@babel/register"],
  "spec": ["test/**/*.test.js"],
  "timeout": 5000,
  "recursive": true,
  "exit": true
}
```

### 命令行选项

```shell
# 指定测试文件
mocha test/**/*.test.js

# 设置超时时间
mocha --timeout 10000

# 监听模式
mocha --watch

# 使用特定 reporter
mocha --reporter spec

# 运行特定测试
mocha --grep "pattern"
```

## 最佳实践

### 1. 使用 chai-as-promised 处理 Promise

```shell
npm install --save-dev chai-as-promised
```

```js
const chai = require('chai');
const chaiAsPromised = require('chai-as-promised');
chai.use(chaiAsPromised);

const { expect } = chai;

it('应该成功获取数据', () => {
  return expect(fetchData()).to.eventually.deep.equal({ data: 'test' });
});

it('应该在错误时拒绝', () => {
  return expect(fetchData()).to.be.rejectedWith('error message');
});
```

### 2. 清理 Sinon

在 `afterEach` 中清理所有 stub 和 spy：

```js
afterEach(() => {
  sinon.restore();
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
    expect(counter.value).to.equal(1);
  });

  it('测试2', () => {
    counter.increment();
    expect(counter.value).to.equal(1); // 独立运行
  });
});
```

### 4. 组织测试结构

```js
// ✅ 好的结构
describe('UserService', () => {
  describe('createUser', () => {
    it('应该创建新用户', () => {});
    it('应该验证邮箱格式', () => {});
  });

  describe('deleteUser', () => {
    it('应该删除存在的用户', () => {});
    it('应该在用户不存在时抛出错误', () => {});
  });
});

// ❌ 不好的结构
describe('tests', () => {
  it('test1', () => {});
  it('test2', () => {});
});
```

### 5. Mock HTTP 请求

```js
const sinon = require('sinon');
const axios = require('axios');
const { expect } = require('chai');

describe('API tests', () => {
  let axiosStub;

  beforeEach(() => {
    axiosStub = sinon.stub(axios, 'get');
  });

  afterEach(() => {
    axiosStub.restore();
  });

  it('应该成功获取用户数据', async () => {
    axiosStub.resolves({ data: { id: 1, name: 'Alice' } });

    const response = await axios.get('/api/users/1');

    expect(response.data.name).to.equal('Alice');
    expect(axiosStub.calledOnce).to.be.true;
  });
});
```

## 常见问题

### 1. TypeScript 支持

安装 ts-node：

```shell
npm install --save-dev ts-node @types/mocha @types/chai @types/sinon
```

配置 Mocha：

```json
// .mocharc.json
{
  "require": ["ts-node/register"],
  "spec": ["test/**/*.test.ts"],
  "extensions": ["ts"]
}
```

### 2. ESM 支持

使用 `.mjs` 扩展名或在 `package.json` 中添加 `"type": "module"`：

```shell
mocha --loader=esmock test/**/*.test.mjs
```

### 3. Timeout 错误

增加超时时间：

```js
// 单个测试
it('耗时操作', function () {
  this.timeout(10000); // 10秒
  // 测试代码
});

// 整个套件
describe('API tests', function () {
  this.timeout(10000);
  // 测试
});
```

### 4. 并行运行测试

```shell
mocha --parallel
```

注意：并行运行时共享状态可能导致问题，确保测试独立。

## Mocha vs Jest/Vitest

| 特性    | Mocha      | Jest/Vitest       |
| ------- | ---------- | ----------------- |
| 断言库  | 需要 Chai  | 内置              |
| Mock 库 | 需要 Sinon | 内置              |
| 覆盖率  | 需要 NYC   | 内置              |
| 配置    | 灵活但复杂 | 开箱即用          |
| 性能    | 中等       | 快（Vitest 最快） |
| 生态    | 成熟       | 现代              |
