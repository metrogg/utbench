# PHP 单元测试规范

## 单元测试编写规范

### 通用规范

- 仅编写单元测试代码，不修改被测试代码
- 尽量做最小变更
- 尽量接近已有代码编写风格
- 限制测试代码的副作用：
  - 测试逻辑应尽量不依赖外部环境，也尽量少对外部环境产生影响
  - 如果测试中涉及读写文件，应使用临时目录，测试结束后清理
  - 被测逻辑依赖外部 HTTP 服务的，应在测试中通过 Mock 或测试服务器替代
  - 被测逻辑依赖数据库的，应通过 Mock 、内存数据库等方式隔离
  - 被测逻辑依赖其他外部服务（ Redis / MQ 等）的，应通过 Mock 或测试替身替代

### PHPUnit 测试框架规范

当项目使用 PHPUnit 时：

#### 测试类结构

```php
<?php

namespace Tests\Unit;

use PHPUnit\Framework\TestCase;
use App\Services\UserService;

class UserServiceTest extends TestCase
{
    private UserService $userService;

    protected function setUp(): void
    {
        parent::setUp();
        // 测试前置准备
        $this->userService = new UserService();
    }

    protected function tearDown(): void
    {
        // 测试后置清理
        parent::tearDown();
    }

    public function testCreateUserSuccess(): void
    {
        // 准备测试数据
        $userData = [
            'name' => 'Test User',
            'email' => 'test@example.com'
        ];

        // 执行被测方法
        $result = $this->userService->createUser($userData);

        // 断言验证
        $this->assertNotNull($result);
        $this->assertEquals('Test User', $result->name);
        $this->assertEquals('test@example.com', $result->email);
    }
}
```

#### Mock 使用规范

**使用 PHPUnit Mock ：**

```php
public function testGetUserOrders(): void
{
    // 创建 Mock 对象
    $orderRepository = $this->createMock(OrderRepository::class);

    // 配置 Mock 行为
    $orderRepository->expects($this->once())
        ->method('findByUserId')
        ->with($this->equalTo(1))
        ->willReturn([
            ['id' => 1, 'total' => 100],
            ['id' => 2, 'total' => 200]
        ]);

    // 注入 Mock 并执行测试
    $userService = new UserService($orderRepository);
    $orders = $userService->getUserOrders(1);

    $this->assertCount(2, $orders);
}
```

**使用 Mockery （如果项目已使用）：**

```php
public function testGetUserOrders(): void
{
    $orderRepository = \Mockery::mock(OrderRepository::class);

    $orderRepository->shouldReceive('findByUserId')
        ->once()
        ->with(1)
        ->andReturn([
            ['id' => 1, 'total' => 100],
            ['id' => 2, 'total' => 200]
        ]);

    $userService = new UserService($orderRepository);
    $orders = $userService->getUserOrders(1);

    $this->assertCount(2, $orders);
}
```

#### 数据提供器

使用数据提供器进行参数化测试：

```php
/**
 * @dataProvider userValidationProvider
 */
public function testUserValidation(array $data, bool $expected): void
{
    $result = $this->userService->validateUser($data);
    $this->assertEquals($expected, $result);
}

public static function userValidationProvider(): array
{
    return [
        'valid user' => [
            ['name' => 'John', 'email' => 'john@example.com'],
            true
        ],
        'empty name' => [
            ['name' => '', 'email' => 'john@example.com'],
            false
        ],
        'invalid email' => [
            ['name' => 'John', 'email' => 'invalid-email'],
            false
        ]
    ];
}
```

#### 异常测试

```php
public function testCreateUserWithDuplicateEmail(): void
{
    $this->expectException(\InvalidArgumentException::class);
    $this->expectExceptionMessage('Email already exists');

    $this->userService->createUser([
        'email' => 'duplicate@example.com'
    ]);
}
```

### Pest 测试框架规范

当项目使用 Pest 时：

```php
<?php

use App\Services\UserService;

beforeEach(function () {
    $this->userService = new UserService();
});

test('create user successfully', function () {
    $result = $this->userService->createUser([
        'name' => 'Test User',
        'email' => 'test@example.com'
    ]);

    expect($result)->not->toBeNull()
        ->and($result->name)->toBe('Test User')
        ->and($result->email)->toBe('test@example.com');
});

test('validate user with invalid email', function () {
    $result = $this->userService->validateUser([
        'name' => 'John',
        'email' => 'invalid'
    ]);

    expect($result)->toBeFalse();
});

it('throws exception for duplicate email', function () {
    $this->userService->createUser([
        'email' => 'duplicate@example.com'
    ]);
})->throws(\InvalidArgumentException::class, 'Email already exists');
```

### Laravel 项目特定规范

当项目使用 Laravel 框架时：

```php
<?php

namespace Tests\Feature;

use Tests\TestCase;
use Illuminate\Foundation\Testing\RefreshDatabase;
use App\Models\User;

class UserServiceTest extends TestCase
{
    use RefreshDatabase;

    public function testCreateUser(): void
    {
        $response = $this->postJson('/api/users', [
            'name' => 'Test User',
            'email' => 'test@example.com'
        ]);

        $response->assertStatus(201)
            ->assertJson([
                'name' => 'Test User',
                'email' => 'test@example.com'
            ]);

        $this->assertDatabaseHas('users', [
            'email' => 'test@example.com'
        ]);
    }
}
```

### 测试最佳实践

1. **测试命名**：测试方法名应清晰描述测试场景，使用 `test` 前缀或 `@test` 注解
2. **AAA 模式**： Arrange(准备) -> Act(执行) -> Assert(断言)
3. **单一职责**：每个测试方法只测试一个场景
4. **避免硬编码**： 使用常量或工厂方法生成测试数据

## 测试代码和函数命名规范

### 测试文件路径

**标准 PHP 项目**：

- 源代码在 `src/` 目录，测试代码在 `tests/` 目录
- 保持目录结构一致，如 `src/Services/UserService.php` -> `tests/Services/UserServiceTest.php`

**Laravel 项目**：

- 单元测试放在 `tests/Unit/` 目录
- 如 `app/Services/UserService.php` -> `tests/Unit/Services/UserServiceTest.php`

**PSR-4 项目**：

- 根据 composer.json 中的 autoload-dev 配置确定测试目录

### 测试类命名

- 测试类名应为被测类名加 `Test` 后缀
- 如 `UserService` -> `UserServiceTest`
- 命名空间应对应测试目录结构

### 测试方法命名

**PHPUnit 风格**：

- 使用 `test` 前缀： `testCreateUser`
- 或使用 `@test` 注解，方法名无前缀限制
- 建议格式： `test{方法名}{场景}`，如 `testCreateUserSuccess`, `testCreateUserWithInvalidEmail`

**Pest 风格**：

- 使用 `test()` 或 `it()` 函数
- 描述性字符串： `test('create user successfully')`

## 运行静态检查方法

在项目根目录运行静态检查：

**使用 PHP_CodeSniffer ：**

```shell
vendor/bin/phpcs --standard=PSR12 tests/
```

**使用 PHPStan ：**

```shell
vendor/bin/phpstan analyse tests/
```

**使用 Psalm ：**

```shell
vendor/bin/psalm tests/
```

优先使用项目中已配置的静态检查工具。检查配置文件如 `phpcs.xml`, `phpstan.neon`, `psalm.xml` 等。

## 单元测试运行方法

### PHPUnit

**运行所有测试：**

```shell
vendor/bin/phpunit
```

**运行指定测试类：**

```shell
vendor/bin/phpunit tests/Unit/Services/UserServiceTest.php
```

**运行指定测试方法：**

```shell
vendor/bin/phpunit --filter testCreateUser
```

**生成覆盖率报告（ `clover` XML 格式）：**

```shell
vendor/bin/phpunit --coverage-clover /tmp/coverage.xml
```

### Pest

**运行所有测试：**

```shell
vendor/bin/pest
```

**运行指定测试文件：**

```shell
vendor/bin/pest tests/Unit/UserServiceTest.php
```

**运行指定测试：**

```shell
vendor/bin/pest --filter "create user successfully"
```

**生成覆盖率报告（ `clover` XML 格式）：**

```shell
vendor/bin/pest --coverage --coverage-clover /tmp/coverage.xml
```

### Laravel Artisan

**运行所有测试：**

```shell
php artisan test
```

**运行指定测试：**

```shell
php artisan test --filter UserServiceTest
```

**生成覆盖率：**

```shell
php artisan test --coverage
```

## 常见测试场景处理

### 测试私有方法

通过反射访问私有方法（谨慎使用，优先测试公共接口）：

```php
public function testPrivateMethod(): void
{
    $reflection = new \ReflectionClass(UserService::class);
    $method = $reflection->getMethod('privateMethod');
    $method->setAccessible(true);

    $result = $method->invoke($this->userService, $arg1, $arg2);

    $this->assertEquals($expected, $result);
}
```

### 测试依赖注入

```php
public function testServiceWithDependencies(): void
{
    $repository = $this->createMock(UserRepository::class);
    $logger = $this->createMock(LoggerInterface::class);

    $service = new UserService($repository, $logger);

    // 测试逻辑...
}
```

### 测试时间相关逻辑

使用时间 Mock 库（如 Carbon ）：

```php
public function testExpiredUser(): void
{
    Carbon::setTestNow('2024-01-01 00:00:00');

    $result = $this->userService->isExpired($user);

    $this->assertTrue($result);

    Carbon::setTestNow(); // 重置
}
```

### 测试文件操作

```php
public function testFileUpload(): void
{
    $tempDir = sys_get_temp_dir() . '/test_' . uniqid();
    mkdir($tempDir);

    try {
        // 测试逻辑
        $result = $this->fileService->upload($file, $tempDir);
        $this->assertFileExists($tempDir . '/uploaded.txt');
    } finally {
        // 清理
        if (is_dir($tempDir)) {
            array_map('unlink', glob("$tempDir/*"));
            rmdir($tempDir);
        }
    }
}
```
