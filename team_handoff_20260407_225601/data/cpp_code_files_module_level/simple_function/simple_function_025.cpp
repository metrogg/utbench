#ifndef GOOGLETEST_INCLUDE_GTEST_GTEST_TYPED_TEST_H_
#define GOOGLETEST_INCLUDE_GTEST_GTEST_TYPED_TEST_H_
#if 0
template <typename T>
class FooTest : public testing::Test {
 public:
  ...
  typedef std::list<T> List;
  static T shared_;
  T value_;
};
typedef testing::Types<char, int, unsigned int> MyTypes;
TYPED_TEST_SUITE(FooTest, MyTypes);
TYPED_TEST(FooTest, DoesBlah) {
  TypeParam n = this->value_;
  n += TestFixture::shared_;
  typename TestFixture::List values;
  values.push_back(n);
  ...
}
TYPED_TEST(FooTest, HasPropertyA) { ... }
#endif  
#if 0
template <typename T>
class FooTest : public testing::Test {
  ...
};
TYPED_TEST_SUITE_P(FooTest);
TYPED_TEST_P(FooTest, DoesBlah) {
  TypeParam n = 0;
  ...
}
TYPED_TEST_P(FooTest, HasPropertyA) { ... }
REGISTER_TYPED_TEST_SUITE_P(FooTest,
                            DoesBlah, HasPropertyA);
typedef testing::Types<char, int, unsigned int> MyTypes;
INSTANTIATE_TYPED_TEST_SUITE_P(My, FooTest, MyTypes);
#endif  
#include "gtest/internal/gtest-internal.h"
#include "gtest/internal/gtest-port.h"
#include "gtest/internal/gtest-type-util.h"
#define GTEST_TYPE_PARAMS_(TestSuiteName) gtest_type_params_##TestSuiteName##_
#define GTEST_NAME_GENERATOR_(TestSuiteName) \
  gtest_type_params_##TestSuiteName##_NameGenerator
#define TYPED_TEST_SUITE(CaseName, Types, ...)                          \
  typedef ::testing::internal::GenerateTypeList<Types>::type            \
      GTEST_TYPE_PARAMS_(CaseName);                                     \
  typedef ::testing::internal::NameGeneratorSelector<__VA_ARGS__>::type \
  GTEST_NAME_GENERATOR_(CaseName)
#define TYPED_TEST(CaseName, TestName)                                       \
  static_assert(sizeof(GTEST_STRINGIFY_(TestName)) > 1,                      \
                "test-name must not be empty");                              \
  template <typename gtest_TypeParam_>                                       \
  class GTEST_TEST_CLASS_NAME_(CaseName, TestName)                           \
      : public CaseName<gtest_TypeParam_> {                                  \
   private:                                                                  \
    typedef CaseName<gtest_TypeParam_> TestFixture;                          \
    typedef gtest_TypeParam_ TypeParam;                                      \
    void TestBody() override;                                                \
  };                                                                         \
  GTEST_INTERNAL_ATTRIBUTE_MAYBE_UNUSED static bool                          \
      gtest_##CaseName##_##TestName##_registered_ =                          \
          ::testing::internal::TypeParameterizedTest<                        \
              CaseName,                                                      \
              ::testing::internal::TemplateSel<GTEST_TEST_CLASS_NAME_(       \
                  CaseName, TestName)>,                                      \
              GTEST_TYPE_PARAMS_(                                            \
                  CaseName)>::Register("",                                   \
                                       ::testing::internal::CodeLocation(    \
                                           __FILE__, __LINE__),              \
                                       GTEST_STRINGIFY_(CaseName),           \
                                       GTEST_STRINGIFY_(TestName), 0,        \
                                       ::testing::internal::GenerateNames<   \
                                           GTEST_NAME_GENERATOR_(CaseName),  \
                                           GTEST_TYPE_PARAMS_(CaseName)>()); \
  template <typename gtest_TypeParam_>                                       \
  void GTEST_TEST_CLASS_NAME_(CaseName,                                      \
                              TestName)<gtest_TypeParam_>::TestBody()
#ifndef GTEST_REMOVE_LEGACY_TEST_CASEAPI_
#define TYPED_TEST_CASE                                                \
  static_assert(::testing::internal::TypedTestCaseIsDeprecated(), ""); \
  TYPED_TEST_SUITE
#endif  
#define GTEST_SUITE_NAMESPACE_(TestSuiteName) gtest_suite_##TestSuiteName##_
#define GTEST_TYPED_TEST_SUITE_P_STATE_(TestSuiteName) \
  gtest_typed_test_suite_p_state_##TestSuiteName##_
#define GTEST_REGISTERED_TEST_NAMES_(TestSuiteName) \
  gtest_registered_test_names_##TestSuiteName##_
#define TYPED_TEST_SUITE_P(SuiteName)              \
  static ::testing::internal::TypedTestSuitePState \
  GTEST_TYPED_TEST_SUITE_P_STATE_(SuiteName)
#ifndef GTEST_REMOVE_LEGACY_TEST_CASEAPI_
#define TYPED_TEST_CASE_P                                                 \
  static_assert(::testing::internal::TypedTestCase_P_IsDeprecated(), ""); \
  TYPED_TEST_SUITE_P
#endif  
#define TYPED_TEST_P(SuiteName, TestName)                         \
  namespace GTEST_SUITE_NAMESPACE_(SuiteName) {                   \
  template <typename gtest_TypeParam_>                            \
  class TestName : public SuiteName<gtest_TypeParam_> {           \
   private:                                                       \
    typedef SuiteName<gtest_TypeParam_> TestFixture;              \
    typedef gtest_TypeParam_ TypeParam;                           \
    void TestBody() override;                                     \
  };                                                              \
  GTEST_INTERNAL_ATTRIBUTE_MAYBE_UNUSED static bool               \
      gtest_##TestName##_defined_ =                               \
          GTEST_TYPED_TEST_SUITE_P_STATE_(SuiteName).AddTestName( \
              __FILE__, __LINE__, GTEST_STRINGIFY_(SuiteName),    \
              GTEST_STRINGIFY_(TestName));                        \
  }                                                               \
  template <typename gtest_TypeParam_>                            \
  void GTEST_SUITE_NAMESPACE_(                                    \
      SuiteName)::TestName<gtest_TypeParam_>::TestBody()
#define REGISTER_TYPED_TEST_SUITE_P(SuiteName, ...)                         \
  namespace GTEST_SUITE_NAMESPACE_(SuiteName) {                             \
  typedef ::testing::internal::Templates<__VA_ARGS__> gtest_AllTests_;      \
  }                                                                         \
  GTEST_INTERNAL_ATTRIBUTE_MAYBE_UNUSED static const char* const            \
  GTEST_REGISTERED_TEST_NAMES_(SuiteName) =                                 \
      GTEST_TYPED_TEST_SUITE_P_STATE_(SuiteName).VerifyRegisteredTestNames( \
          GTEST_STRINGIFY_(SuiteName), __FILE__, __LINE__, #__VA_ARGS__)
#ifndef GTEST_REMOVE_LEGACY_TEST_CASEAPI_
#define REGISTER_TYPED_TEST_CASE_P                                           \
  static_assert(::testing::internal::RegisterTypedTestCase_P_IsDeprecated(), \
                "");                                                         \
  REGISTER_TYPED_TEST_SUITE_P
#endif  
#define INSTANTIATE_TYPED_TEST_SUITE_P(Prefix, SuiteName, Types, ...)        \
  static_assert(sizeof(GTEST_STRINGIFY_(Prefix)) > 1,                        \
                "test-suit-prefix must not be empty");                       \
  GTEST_INTERNAL_ATTRIBUTE_MAYBE_UNUSED static bool                          \
      gtest_##Prefix##_##SuiteName =                                         \
          ::testing::internal::TypeParameterizedTestSuite<                   \
              SuiteName, GTEST_SUITE_NAMESPACE_(SuiteName)::gtest_AllTests_, \
              ::testing::internal::GenerateTypeList<Types>::type>::          \
              Register(                                                      \
                  GTEST_STRINGIFY_(Prefix),                                  \
                  ::testing::internal::CodeLocation(__FILE__, __LINE__),     \
                  &GTEST_TYPED_TEST_SUITE_P_STATE_(SuiteName),               \
                  GTEST_STRINGIFY_(SuiteName),                               \
                  GTEST_REGISTERED_TEST_NAMES_(SuiteName),                   \
                  ::testing::internal::GenerateNames<                        \
                      ::testing::internal::NameGeneratorSelector<            \
                          __VA_ARGS__>::type,                                \
                      ::testing::internal::GenerateTypeList<Types>::type>())
#ifndef GTEST_REMOVE_LEGACY_TEST_CASEAPI_
#define INSTANTIATE_TYPED_TEST_CASE_P                                      \
  static_assert(                                                           \
      ::testing::internal::InstantiateTypedTestCase_P_IsDeprecated(), ""); \
  INSTANTIATE_TYPED_TEST_SUITE_P
#endif  
#endif  
#include "gtest/gtest-typed-test.h"
#include <set>
#include <string>
#include <vector>
#include "gtest/gtest.h"
namespace testing {
namespace internal {
static const char* SkipSpaces(const char* str) {
  while (IsSpace(*str)) str++;
  return str;
}
static std::vector<std::string> SplitIntoTestNames(const char* src) {
  std::vector<std::string> name_vec;
  src = SkipSpaces(src);
  for (; src != nullptr; src = SkipComma(src)) {
    name_vec.push_back(StripTrailingSpaces(GetPrefixUntilComma(src)));
  }
  return name_vec;
}
const char* TypedTestSuitePState::VerifyRegisteredTestNames(
    const char* test_suite_name, const char* file, int line,
    const char* registered_tests) {
  RegisterTypeParameterizedTestSuite(test_suite_name, CodeLocation(file, line));
  typedef RegisteredTestsMap::const_iterator RegisteredTestIter;
  registered_ = true;
  std::vector<std::string> name_vec = SplitIntoTestNames(registered_tests);
  Message errors;
  std::set<std::string> tests;
  for (std::vector<std::string>::const_iterator name_it = name_vec.begin();
       name_it != name_vec.end(); ++name_it) {
    const std::string& name = *name_it;
    if (tests.count(name) != 0) {
      errors << "Test " << name << " is listed more than once.\n";
      continue;
    }
    if (registered_tests_.count(name) != 0) {
      tests.insert(name);
    } else {
      errors << "No test named " << name
             << " can be found in this test suite.\n";
    }
  }
  for (RegisteredTestIter it = registered_tests_.begin();
       it != registered_tests_.end(); ++it) {
    if (tests.count(it->first) == 0) {
      errors << "You forgot to list test " << it->first << ".\n";
    }
  }
  const std::string& errors_str = errors.GetString();
  if (!errors_str.empty()) {
    fprintf(stderr, "%s %s", FormatFileLocation(file, line).c_str(),
            errors_str.c_str());
    fflush(stderr);
    posix::Abort();
  }
  return registered_tests;
}
}  
}  