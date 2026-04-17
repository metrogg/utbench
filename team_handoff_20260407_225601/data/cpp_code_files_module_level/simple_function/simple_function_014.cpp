#ifndef TENSORFLOW_CORE_PLATFORM_PLATFORM_STRINGS_H_
#define TENSORFLOW_CORE_PLATFORM_PLATFORM_STRINGS_H_
#include <stdio.h>
#include <string>
#include <vector>
#define TF_PLAT_STR_VERSION_ "1.0"
#define TF_PLAT_STR_MAGIC_PREFIX_ "\0S\\s\":^p*L}"
#define TF_PLAT_STR_STR_1_(x) #x
#define TF_PLAT_STR_AS_STR_(x) TF_PLAT_STR_STR_1_(x)
#define TF_PLAT_STR_TERMINATOR_
#define TF_PLAT_STR_(x) TF_PLAT_STR_MAGIC_PREFIX_ #x "=" TF_PLAT_STR_AS_STR_(x)
#include "tensorflow/core/platform/platform_strings_computed.h"
#define TF_PLAT_STR_LIST___x86_64__()                                      \
        TF_PLAT_STR__M_IX86_FP                                             \
        TF_PLAT_STR__NO_PREFETCHW                                          \
        TF_PLAT_STR___3dNOW_A__                                            \
        TF_PLAT_STR___3dNOW__                                              \
        TF_PLAT_STR___ABM__                                                \
        TF_PLAT_STR___ADX__                                                \
        TF_PLAT_STR___AES__                                                \
        TF_PLAT_STR___AVX2__                                               \
        TF_PLAT_STR___AVX512BW__                                           \
        TF_PLAT_STR___AVX512CD__                                           \
        TF_PLAT_STR___AVX512DQ__                                           \
        TF_PLAT_STR___AVX512ER__                                           \
        TF_PLAT_STR___AVX512F__                                            \
        TF_PLAT_STR___AVX512IFMA__                                         \
        TF_PLAT_STR___AVX512PF__                                           \
        TF_PLAT_STR___AVX512VBMI__                                         \
        TF_PLAT_STR___AVX512VL__                                           \
        TF_PLAT_STR___AVX__                                                \
        TF_PLAT_STR___BMI2__                                               \
        TF_PLAT_STR___BMI__                                                \
        TF_PLAT_STR___CLFLUSHOPT__                                         \
        TF_PLAT_STR___CLZERO__                                             \
        TF_PLAT_STR___F16C__                                               \
        TF_PLAT_STR___FMA4__                                               \
        TF_PLAT_STR___FMA__                                                \
        TF_PLAT_STR___FP_FAST_FMA                                          \
        TF_PLAT_STR___FP_FAST_FMAF                                         \
        TF_PLAT_STR___FSGSBASE__                                           \
        TF_PLAT_STR___FXSR__                                               \
        TF_PLAT_STR___LWP__                                                \
        TF_PLAT_STR___LZCNT__                                              \
        TF_PLAT_STR___MMX__                                                \
        TF_PLAT_STR___MWAITX__                                             \
        TF_PLAT_STR___PCLMUL__                                             \
        TF_PLAT_STR___PKU__                                                \
        TF_PLAT_STR___POPCNT__                                             \
        TF_PLAT_STR___PRFCHW__                                             \
        TF_PLAT_STR___RDRND__                                              \
        TF_PLAT_STR___RDSEED__                                             \
        TF_PLAT_STR___RTM__                                                \
        TF_PLAT_STR___SHA__                                                \
        TF_PLAT_STR___SSE2_MATH__                                          \
        TF_PLAT_STR___SSE2__                                               \
        TF_PLAT_STR___SSE_MATH__                                           \
        TF_PLAT_STR___SSE__                                                \
        TF_PLAT_STR___SSE3__                                               \
        TF_PLAT_STR___SSE4A__                                              \
        TF_PLAT_STR___SSE4_1__                                             \
        TF_PLAT_STR___SSE4_2__                                             \
        TF_PLAT_STR___SSSE3__                                              \
        TF_PLAT_STR___TBM__                                                \
        TF_PLAT_STR___XOP__                                                \
        TF_PLAT_STR___XSAVEC__                                             \
        TF_PLAT_STR___XSAVEOPT__                                           \
        TF_PLAT_STR___XSAVES__                                             \
        TF_PLAT_STR___XSAVE__                                              \
        TF_PLAT_STR_TERMINATOR_
#define TF_PLAT_STR_LIST___powerpc64__()                                   \
        TF_PLAT_STR__SOFT_DOUBLE                                           \
        TF_PLAT_STR__SOFT_FLOAT                                            \
        TF_PLAT_STR___ALTIVEC__                                            \
        TF_PLAT_STR___APPLE_ALTIVEC__                                      \
        TF_PLAT_STR___CRYPTO__                                             \
        TF_PLAT_STR___FLOAT128_HARDWARE__                                  \
        TF_PLAT_STR___FLOAT128_TYPE__                                      \
        TF_PLAT_STR___FP_FAST_FMA                                          \
        TF_PLAT_STR___FP_FAST_FMAF                                         \
        TF_PLAT_STR___HTM__                                                \
        TF_PLAT_STR___NO_FPRS__                                            \
        TF_PLAT_STR___NO_LWSYNC__                                          \
        TF_PLAT_STR___POWER8_VECTOR__                                      \
        TF_PLAT_STR___POWER9_VECTOR__                                      \
        TF_PLAT_STR___PPC405__                                             \
        TF_PLAT_STR___QUAD_MEMORY_ATOMIC__                                 \
        TF_PLAT_STR___RECIPF__                                             \
        TF_PLAT_STR___RECIP_PRECISION__                                    \
        TF_PLAT_STR___RECIP__                                              \
        TF_PLAT_STR___RSQRTEF__                                            \
        TF_PLAT_STR___RSQRTE__                                             \
        TF_PLAT_STR___TM_FENCE__                                           \
        TF_PLAT_STR___UPPER_REGS_DF__                                      \
        TF_PLAT_STR___UPPER_REGS_SF__                                      \
        TF_PLAT_STR___VEC__                                                \
        TF_PLAT_STR___VSX__                                                \
        TF_PLAT_STR_TERMINATOR_
#define TF_PLAT_STR_LIST___aarch64__()                                     \
        TF_PLAT_STR___ARM_ARCH                                             \
        TF_PLAT_STR___ARM_FEATURE_CLZ                                      \
        TF_PLAT_STR___ARM_FEATURE_CRC32                                    \
        TF_PLAT_STR___ARM_FEATURE_CRC32                                    \
        TF_PLAT_STR___ARM_FEATURE_CRYPTO                                   \
        TF_PLAT_STR___ARM_FEATURE_DIRECTED_ROUNDING                        \
        TF_PLAT_STR___ARM_FEATURE_DSP                                      \
        TF_PLAT_STR___ARM_FEATURE_FMA                                      \
        TF_PLAT_STR___ARM_FEATURE_IDIV                                     \
        TF_PLAT_STR___ARM_FEATURE_LDREX                                    \
        TF_PLAT_STR___ARM_FEATURE_NUMERIC_MAXMIN                           \
        TF_PLAT_STR___ARM_FEATURE_QBIT                                     \
        TF_PLAT_STR___ARM_FEATURE_QRDMX                                    \
        TF_PLAT_STR___ARM_FEATURE_SAT                                      \
        TF_PLAT_STR___ARM_FEATURE_SIMD32                                   \
        TF_PLAT_STR___ARM_FEATURE_UNALIGNED                                \
        TF_PLAT_STR___ARM_FP                                               \
        TF_PLAT_STR___ARM_NEON_FP                                          \
        TF_PLAT_STR___ARM_NEON__                                           \
        TF_PLAT_STR___ARM_WMMX                                             \
        TF_PLAT_STR___IWMMXT2__                                            \
        TF_PLAT_STR___IWMMXT__                                             \
        TF_PLAT_STR___VFP_FP__                                             \
        TF_PLAT_STR_TERMINATOR_
#define TF_PLAT_STR_LIST___generic__()                                     \
        TF_PLAT_STR_TARGET_IPHONE_SIMULATOR                                \
        TF_PLAT_STR_TARGET_OS_IOS                                          \
        TF_PLAT_STR_TARGET_OS_IPHONE                                       \
        TF_PLAT_STR__MSC_VER                                               \
        TF_PLAT_STR__M_ARM                                                 \
        TF_PLAT_STR__M_ARM64                                               \
        TF_PLAT_STR__M_ARM_ARMV7VE                                         \
        TF_PLAT_STR__M_ARM_FP                                              \
        TF_PLAT_STR__M_IX86                                                \
        TF_PLAT_STR__M_X64                                                 \
        TF_PLAT_STR__WIN32                                                 \
        TF_PLAT_STR__WIN64                                                 \
        TF_PLAT_STR___ANDROID__                                            \
        TF_PLAT_STR___APPLE__                                              \
        TF_PLAT_STR___BYTE_ORDER__                                         \
        TF_PLAT_STR___CYGWIN__                                             \
        TF_PLAT_STR___FreeBSD__                                            \
        TF_PLAT_STR___LITTLE_ENDIAN__                                      \
        TF_PLAT_STR___NetBSD__                                             \
        TF_PLAT_STR___OpenBSD__                                            \
        TF_PLAT_STR_____MSYS__                                             \
        TF_PLAT_STR___aarch64__                                            \
        TF_PLAT_STR___alpha__                                              \
        TF_PLAT_STR___arm__                                                \
        TF_PLAT_STR___i386__                                               \
        TF_PLAT_STR___i686__                                               \
        TF_PLAT_STR___ia64__                                               \
        TF_PLAT_STR___linux__                                              \
        TF_PLAT_STR___mips32__                                             \
        TF_PLAT_STR___mips64__                                             \
        TF_PLAT_STR___powerpc64__                                          \
        TF_PLAT_STR___powerpc__                                            \
        TF_PLAT_STR___riscv___                                             \
        TF_PLAT_STR___s390x__                                              \
        TF_PLAT_STR___sparc64__                                            \
        TF_PLAT_STR___sparc__                                              \
        TF_PLAT_STR___x86_64__                                             \
        TF_PLAT_STR_TERMINATOR_
#if !defined(__x86_64__) && !defined(_M_X64) && \
    !defined(__i386__) && !defined(_M_IX86)
#undef TF_PLAT_STR_LIST___x86_64__
#define TF_PLAT_STR_LIST___x86_64__()
#endif
#if !defined(__powerpc64__) && !defined(__powerpc__)
#undef TF_PLAT_STR_LIST___powerpc64__
#define TF_PLAT_STR_LIST___powerpc64__()
#endif
#if !defined(__aarch64__) && !defined(_M_ARM64) && \
    !defined(__arm__) && !defined(_M_ARM)
#undef TF_PLAT_STR_LIST___aarch64__
#define TF_PLAT_STR_LIST___aarch64__()
#endif
#define TF_PLATFORM_STRINGS()                                                  \
    static const char tf_cpu_option[] =                                        \
        TF_PLAT_STR_MAGIC_PREFIX_ "TF_PLAT_STR_VERSION=" TF_PLAT_STR_VERSION_  \
        TF_PLAT_STR_LIST___x86_64__()                                          \
        TF_PLAT_STR_LIST___powerpc64__()                                       \
        TF_PLAT_STR_LIST___aarch64__()                                         \
        TF_PLAT_STR_LIST___generic__()                                         \
    ;                                                                          \
    const char *tf_cpu_option_global;                                          \
    namespace {                                                                \
    class TFCPUOptionHelper {                                                  \
     public:                                                                   \
      TFCPUOptionHelper() {                                                    \
             \
                     \
        tf_cpu_option_global = tf_cpu_option;                                  \
                 \
        printf("%s%s", tf_cpu_option, "");                                     \
      }                                                                        \
    } tf_cpu_option_avoid_omit_class;                                          \
    }  
namespace tensorflow {
int GetPlatformStrings(const std::string& path,
                       std::vector<std::string>* found);
}  
#endif  
#include "tensorflow/core/platform/platform_strings.h"
#include <cerrno>
#include <cstdio>
#include <cstring>
#include <string>
#include <vector>
namespace tensorflow {
int GetPlatformStrings(const std::string& path,
                       std::vector<std::string>* found) {
  int result;
  FILE* ifp = fopen(path.c_str(), "rb");
  if (ifp != nullptr) {
    static const char prefix[] = TF_PLAT_STR_MAGIC_PREFIX_;
    int first_char = prefix[1];
    int last_char = -1;
    int c;
    while ((c = getc(ifp)) != EOF) {
      if (c == first_char && last_char == 0) {
        int i = 2;
        while (prefix[i] != 0 && (c = getc(ifp)) == prefix[i]) {
          i++;
        }
        if (prefix[i] == 0) {
          std::string str;
          while ((c = getc(ifp)) != EOF && c != 0) {
            str.push_back(c);
          }
          if (!str.empty()) {
            found->push_back(str);
          }
        }
      }
      last_char = c;
    }
    result = (ferror(ifp) == 0) ? 0 : errno;
    if (fclose(ifp) != 0) {
      result = errno;
    }
  } else {
    result = errno;
  }
  return result;
}
}  