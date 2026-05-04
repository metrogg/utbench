# Java 单元测试规范

## 单元测试编写规范

### 基本原则

- 对于每个被测函数仅修改或新增单个测试函数，所有补充的测试逻辑都编写在同一个测试函数中
- 仅编写单元测试代码，不修改被测试代码
- 尽量做最小变更
- 尽量接近已有代码编写风格
- 限制测试代码的副作用：
  - 测试逻辑应该尽量不依赖外部环境，也尽量少对外部环境产生影响
  - 如果测试中涉及读写文件，应使用临时目录或在测试结束后清理
  - 被测逻辑依赖外部 HTTP 服务的，应在测试中通过 Mock 框架模拟外部服务响应
  - 对外部依赖、数据库访问、远程调用等使用 Mock 框架进行隔离

### 测试框架选择

#### 判断项目使用的测试框架

⚠️ 根据项目根目录pom.xml中声明的依赖进行判断

- 如果`pom.xml`中包含`org,testng`依赖，则该项目使用的是`TestNG`
  ```xml
  <dependency>
    <groupId>org.testng</groupId>
    <artifactId>testng</artifactId>
    <version>${testng.version}</version>
    <scope>test</scope>
  </dependency>
  ```
- 如果`pom.xml`中没有明确声明测试框架，则该项目用的就是`JUnit`

#### 各类框架使用说明

- **JUnit 5** （推荐）：现代化的测试框架,提供更强大的功能
  - 核心注解： `@Test`, `@BeforeEach`, `@AfterEach`, `@BeforeAll`, `@AfterAll`
  - 参数化测试： `@ParameterizedTest`, `@ValueSource`, `@MethodSource`
  - 断言：使用 `org.junit.jupiter.api.Assertions`
- **JUnit 4** ：如果项目已使用或有兼容性要求
  - 核心注解： `@Test`, `@Before`, `@After`, `@BeforeClass`, `@AfterClass`
  - 运行器： `@RunWith(MockitoJUnitRunner.class)`
  - 断言：使用 `org.junit.Assert`

- **TestNG** ：如果项目已使用
  - 核心注解： `@Test`, `@BeforeMethod`, `@AfterMethod`, `@BeforeClass`, `@AfterClass`
  - 数据驱动： `@DataProvider`
  - 断言：使用 `org.testng.Assert`

### Mock 框架使用

优先使用项目中已经使用的 Mock 框架。如果项目中未明确使用 Mock 框架，推荐使用 Mockito ：

- 核心注解： `@Mock`, `@InjectMocks`, `@Spy`
- 核心方法： `when()`, `thenReturn()`, `verify()`, `any()`, `eq()`
- 对外部调用、数据库访问、远程服务、随机数、时间等使用 Mock 进行隔离
- 不要在测试代码中重新实现被测试方法的核心逻辑
- 对于使用 Lombok `@AllArgsConstructor` 注解类的Mock，被测类使用`@AllArgsConstructor`注解修饰并且注入的对象使用`final`修饰，该对象被Mock后无法重复修改Mock行为，对这种情况，需要在`setUp`中使用`new`关键字创建被测类对象,例如被测类`BookService`，注入了`OrderService`对象并使用`final`修饰：

  ```java
  @Slf4j
  @Component
  @AllArgsConstructor
  public class BookService {

      private final OrderService orderService;

  }

  // 单测用例中需要在setUp中使用`new`构建被测对象`BookService`:
  public class BookServiceTest {
    @InjectMocks
    private BookService bookService;

    @Mock
    private OrderService orderService;

    @BeforeMethod
    public void setUp() {
        MockitoAnnotations.openMocks(this);
        // 因为被测类使用@AllArgsConstructor注解，因此可以直接使用`new`实例化对象
        bookService = new BookService(orderService);
    }
  }
  ```

### 测试类结构

```java
/**
 * {@link ClassName} 的单元测试
 * 测试要点:
 * 1. 正常流程测试
 * 2. 边界条件测试
 * 3. 异常情况测试
 */
public class ClassNameTest {

    // 测试对象
    @InjectMocks
    private ClassName testInstance;

    // Mock 依赖
    @Mock
    private DependencyService dependencyService;

    // 测试前的初始化
    @BeforeEach  // JUnit 5
    // 或 @Before  // JUnit 4
    // 或 @BeforeMethod  // TestNG
    public void setUp() {
        MockitoAnnotations.openMocks(this);
        // 额外初始化
    }

    // 具体测试方法
    @Test
    public void testMethodNameScenario() {
        // 测试实现
    }
}
```

### 测试方法结构 (AAA 模式)

```java
@Test
public void testMethodNameScenario() {
    // Arrange: 准备测试数据和环境
    // 1. 构造输入参数
    // 2. 设置 Mock 行为

    // Act: 执行被测试的方法
    // 调用被测方法

    // Assert: 验证测试结果
    // 1. 验证返回值
    // 2. 验证状态变化
    // 3. 验证依赖调用
}
```

### 断言使用

选择合适的断言方法，提供清晰的错误信息：

```java
// JUnit 5
assertNotNull(result);
assertEquals(expected, actual);
assertTrue(condition);
assertThrows(ExceptionClass.class, () -> {
    // 会抛异常的代码
});

// JUnit 4
assertNotNull(result);
assertEquals(expected, actual);
assertTrue(condition);
// 异常测试
@Test(expected = ExceptionClass.class)
public void testMethodThrowsException() {
    // 会抛异常的代码
}

// TestNG
assertNotNull(result);
assertEquals(actual, expected);  // 注意参数顺序与 JUnit 相反
assertTrue(condition);
expectThrows(ExceptionClass.class, () -> {
    // 会抛异常的代码
});
```

### 测试数据构造

- 使用有意义的测试数据，避免无意义的 "test", "123" 等
- 对于使用 Lombok `@Builder` 或 `@SuperBuilder` 的类，使用 builder 模式：
  ```java
  User user = User.builder()
      .userId("user_001")
      .username("zhangsan")
      .email("zhangsan@example.com")
      .build();
  ```
- 对于使用普通类，使用构造函数或 setter 方法：
  ```java
  User user = new User();
  user.setUserId("user_001");
  user.setUsername("zhangsan");
  user.setEmail("zhangsan@example.com");
  ```
- 静态工厂方法构造对象：

  ```java
  @Data
  public class Result<T> {
      private boolean success;
      private T data;
      private String code;

      protected Result() {} // 构造函数为protected

      public static <T> Result<T> fromSuccess(T data) {
          Result<T> result = new Result<>();
          result.setSuccess(true);
          result.setData(data);
          return result;
      }
  }
  // 测试用例中使用静态工厂方法
  @Test
  public void testResultConstruction() {
      UserVo userData = new UserVo("user123", "张三");
      Result<UserVo> successResult = Result.fromSuccess(userData);

      assertTrue(successResult.isSuccess());
      assertEquals(successResult.getData().getName(), "张三");
  }
  ```

- 使用 Lombok `@Data` 注解的类, 直接食用new关键字构造：

  ```java
  @Data
  public class QueryDealShipmentRequest {
      @NotNull
      private List<String> dealIds;
      private String userId;
  }
  // 测试用例中使用new关键字构造对象
  @Test
  public void testQueryDealShipmentRequestConstruction() {
      QueryDealShipmentRequest request = new QueryDealShipmentRequest();

      request.setDealIds(Arrays.asList("deal001", "deal002"));
      request.setUserId("user123");

      assertNotNull(request);
      assertEquals(request.getDealIds().size(), 2);
  }
  ```

- 构造对象时需要处理必填字段：对于使用 Lombok `@NotNull/@NonNull` 注解修饰的字段，构造对象时必须赋值，不能为null
- 每个 setter 调用独立成行，不使用链式调用

### 注释规范

每个测试方法必须包含清晰的注释：

```java
/**
 * 测试创建用户 - 正常流程
 *
 * 测试场景: 使用有效的用户信息创建新用户
 * 前置条件: 用户名和邮箱在系统中不存在
 * 输入数据: 包含有效用户名、邮箱、年龄的用户信息
 * 预期结果: 成功创建用户并返回用户 ID
 */
@Test
public void testCreateUserSuccess() {
    // 测试实现
}
```

### 导入依赖处理

- 需要从被测类中复制相关的 import 语句到测试类
- 添加测试框架相关的 import ：

  ```java
  // JUnit 5
  import org.junit.jupiter.api.Test;
  import org.junit.jupiter.api.BeforeEach;
  import static org.junit.jupiter.api.Assertions.*;

  // JUnit 4
  import org.junit.Test;
  import org.junit.Before;
  import static org.junit.Assert.*;

  // TestNG
  import org.testng.annotations.Test;
  import org.testng.annotations.BeforeMethod;
  import static org.testng.Assert.*;

  // Mockito
  import org.mockito.Mock;
  import org.mockito.InjectMocks;
  import org.mockito.MockitoAnnotations;
  import static org.mockito.Mockito.*;
  ```

### 私有成员访问

如果必须访问或修改私有成员，使用反射：

```java
private <T> T getPrivateField(Object object, String fieldName)
        throws NoSuchFieldException, IllegalAccessException {
    Field field = object.getClass().getDeclaredField(fieldName);
    field.setAccessible(true);
    return (T) field.get(object);
}

private void setPrivateField(Object object, String fieldName, Object value)
        throws NoSuchFieldException, IllegalAccessException {
    Field field = object.getClass().getDeclaredField(fieldName);
    field.setAccessible(true);
    field.set(object, value);
}
```

## 测试代码和函数命名规范

- **测试类命名**：被测试类名 + `Test` 后缀
  - 示例： `UserService` → `UserServiceTest`
- **测试文件路径**：
  - Maven 项目： `src/test/java/{package_path}/{ClassName}Test.java`
  - Gradle 项目： `src/test/java/{package_path}/{ClassName}Test.java`
  - 目录结构与被测类保持一致
- **测试函数命名**： `test{方法名}{场景描述}` 清晰表达测试场景, 这里必须使用大驼峰命名
  - 示例： `testCreateUserSuccess`, `testFindByIdNotFound`

## 运行编译检查方法

### Maven 项目

在项目根目录运行编译检查：

```shell
mvn clean compile -DskipTests
```

如果是多模块项目且只需要编译特定模块：

```shell
mvn clean compile -DskipTests -pl {module-name} -am
```

其中 `{module-name}` 是模块名， `-am` 表示同时编译依赖的模块。

### Gradle 项目

在项目根目录运行编译检查：

```shell
./gradlew clean compileJava compileTestJava
```

如果是多模块项目且只需要编译特定模块：

```shell
./gradlew :{module-name}:clean :{module-name}:compileJava :{module-name}:compileTestJava
```

## 单元测试运行方法

### Maven 项目

在项目根目录运行测试：

```shell
mvn test -Dtest={TestClassName}#{testMethodName}
```

- `{TestClassName}`：测试类名（不含包名）
- `{testMethodName}`：测试方法名（可选,不指定则运行整个测试类）

多模块项目指定模块运行：

```shell
mvn test -pl {module-name} -am -Dtest={TestClassName}#{testMethodName}
```

生成覆盖率报告（使用 JaCoCo ）：

```shell
mvn test -Dtest={TestClassName} jacoco:report
```

覆盖率数据文件位于： `target/site/jacoco/jacoco.xml`

### Gradle 项目

在项目根目录运行测试：

```shell
./gradlew test --tests {package}.{TestClassName}.{testMethodName}
```

多模块项目指定模块运行：

```shell
./gradlew :{module-name}:test --tests {package}.{TestClassName}.{testMethodName}
```

生成覆盖率报告（使用 JaCoCo ）：

```shell
./gradlew test jacocoTestReport
```

覆盖率数据文件位于： `build/reports/jacoco/test/jacocoTestReport.xml`

### 提交覆盖率到 MCP

根据构建工具提交相应格式的覆盖率数据：

- **Maven + JaCoCo** ：使用 `jacoco` XML 格式，路径为 `target/site/jacoco/jacoco.xml`
- **Gradle + JaCoCo** ：使用 `jacoco` XML 格式，路径为 `build/reports/jacoco/test/jacocoTestReport.xml`
