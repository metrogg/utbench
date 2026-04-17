#ifndef XLA_PARSE_FLAGS_FROM_ENV_H_
#define XLA_PARSE_FLAGS_FROM_ENV_H_
#include <vector>
#include "absl/strings/string_view.h"
#include "xla/tsl/util/command_line_flags.h"
#include "xla/types.h"
namespace xla {
void ParseFlagsFromEnvAndDieIfUnknown(absl::string_view envvar,
                                      const std::vector<tsl::Flag>& flag_list);
void ParseFlagsFromEnvAndIgnoreUnknown(absl::string_view envvar,
                                       const std::vector<tsl::Flag>& flag_list);
void ResetFlagsFromEnvForTesting(absl::string_view envvar, int** pargc,
                                 std::vector<char*>** pargv);
}  
#endif  
#include "xla/parse_flags_from_env.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <memory>
#include <string>
#include <vector>
#include "absl/container/flat_hash_map.h"
#include "absl/strings/ascii.h"
#include "absl/strings/str_format.h"
#include "absl/strings/str_join.h"
#include "absl/types/span.h"
#include "xla/tsl/util/command_line_flags.h"
#include "tsl/platform/logging.h"
namespace xla {
static const char kWS[] = " \t\r\n";  
namespace {
struct FreeDeleter {
  void operator()(char* ptr) { free(ptr); }
};
struct EnvArgv {
  EnvArgv() : initialized(false), argc(0) {}
  bool initialized;         
  int argc;                 
  std::vector<char*> argv;  
  std::vector<std::unique_ptr<char, FreeDeleter>> argv_save;
};
}  
static void AppendToEnvArgv(const char* s0, size_t s0len, const char* s1,
                            size_t s1len, EnvArgv* a) {
  if (s0 == nullptr) {
    a->argv.push_back(nullptr);
    a->argv_save.push_back(nullptr);
  } else {
    std::string s = std::string(s0, s0len) + std::string(s1, s1len);
    char* str = strdup(s.c_str());
    a->argv.push_back(str);
    a->argv_save.emplace_back(str);
    a->argc++;
  }
}
static size_t FindFirstOf(const std::string& s, const char* x, size_t pos) {
  size_t result = s.find_first_of(x, pos);
  return result == std::string::npos ? s.size() : result;
}
static size_t FindFirstNotOf(const std::string& s, const char* x, size_t pos) {
  size_t result = s.find_first_not_of(x, pos);
  return result == std::string::npos ? s.size() : result;
}
static void ParseArgvFromString(const std::string& flag_str, EnvArgv* a) {
  size_t b = FindFirstNotOf(flag_str, kWS, 0);
  while (b != flag_str.size() && flag_str[b] == '-') {
    size_t e = b;
    while (e != flag_str.size() && isascii(flag_str[e]) &&
           (strchr("-_", flag_str[e]) != nullptr ||
            absl::ascii_isalnum(flag_str[e]))) {
      e++;
    }
    if (e != flag_str.size() && flag_str[e] == '=' &&
        e + 1 != flag_str.size() && strchr("'\"", flag_str[e + 1]) != nullptr) {
      int c;
      e++;  
      size_t eflag = e;
      char quote = flag_str[e];
      e++;  
      std::string value;
      for (; e != flag_str.size() && (c = flag_str[e]) != quote; e++) {
        if (quote == '"' && c == '\\' && e + 1 != flag_str.size()) {
          e++;
          c = flag_str[e];
        }
        value += c;
      }
      if (e != flag_str.size()) {  
        e++;
      }
      AppendToEnvArgv(flag_str.data() + b, eflag - b, value.data(),
                      value.size(), a);
    } else {  
      e = FindFirstOf(flag_str, kWS, e);
      AppendToEnvArgv(flag_str.data() + b, e - b, "", 0, a);
    }
    b = FindFirstNotOf(flag_str, kWS, e);
  }
}
static void SetArgvFromEnv(absl::string_view envvar, EnvArgv* a) {
  if (!a->initialized) {
    static const char kDummyArgv[] = "<argv[0]>";
    AppendToEnvArgv(kDummyArgv, strlen(kDummyArgv), nullptr, 0,
                    a);  
    const char* env = getenv(std::string(envvar).c_str());
    if (env == nullptr || env[0] == '\0') {
    } else if (env[strspn(env, kWS)] == '-') {  
      ParseArgvFromString(env, a);
    } else {  
      FILE* fp = fopen(env, "r");
      if (fp != nullptr) {
        std::string str;
        char buf[512];
        int n;
        while ((n = fread(buf, 1, sizeof(buf), fp)) > 0) {
          str.append(buf, n);
        }
        fclose(fp);
        ParseArgvFromString(str, a);
      } else {
        LOG(QFATAL)
            << "Could not open file \"" << env
            << "\" to read flags for environment variable \"" << envvar
            << "\". (We assumed \"" << env
            << "\" was a file name because it did not start with a \"--\".)";
      }
    }
    AppendToEnvArgv(nullptr, 0, nullptr, 0, a);  
    a->initialized = true;
  }
}
static absl::flat_hash_map<std::string, EnvArgv>& EnvArgvs() {
  static auto* env_argvs = new absl::flat_hash_map<std::string, EnvArgv>();
  return *env_argvs;
}
static absl::Mutex env_argv_mu(absl::kConstInit);
static void DieIfEnvHasUnknownFlagsLeft(absl::string_view envvar);
void ParseFlagsFromEnvAndDieIfUnknown(absl::string_view envvar,
                                      const std::vector<tsl::Flag>& flag_list) {
  ParseFlagsFromEnvAndIgnoreUnknown(envvar, flag_list);
  DieIfEnvHasUnknownFlagsLeft(envvar);
}
void ParseFlagsFromEnvAndIgnoreUnknown(
    absl::string_view envvar, const std::vector<tsl::Flag>& flag_list) {
  absl::MutexLock lock(&env_argv_mu);
  auto* env_argv = &EnvArgvs()[envvar];
  SetArgvFromEnv(envvar, env_argv);  
  if (VLOG_IS_ON(1)) {
    VLOG(1) << "For env var " << envvar << " found arguments:";
    for (int i = 0; i < env_argv->argc; i++) {
      VLOG(1) << "  argv[" << i << "] = " << env_argv->argv[i];
    }
  }
  QCHECK(tsl::Flags::Parse(&env_argv->argc, env_argv->argv.data(), flag_list))
      << "Flag parsing failed.\n"
      << tsl::Flags::Usage(getenv(std::string(envvar).c_str()), flag_list);
}
static void DieIfEnvHasUnknownFlagsLeft(absl::string_view envvar) {
  absl::MutexLock lock(&env_argv_mu);
  auto* env_argv = &EnvArgvs()[envvar];
  SetArgvFromEnv(envvar, env_argv);
  if (env_argv->argc != 1) {
    auto unknown_flags = absl::MakeSpan(env_argv->argv);
    unknown_flags.remove_prefix(1);
    LOG(QFATAL) << "Unknown flag" << (unknown_flags.size() > 1 ? "s" : "")
                << " in " << envvar << ": "
                << absl::StrJoin(unknown_flags, " ");
  }
}
void ResetFlagsFromEnvForTesting(absl::string_view envvar, int** pargc,
                                 std::vector<char*>** pargv) {
  absl::MutexLock lock(&env_argv_mu);
  EnvArgvs().erase(envvar);
  auto& env_argv = EnvArgvs()[envvar];
  *pargc = &env_argv.argc;
  *pargv = &env_argv.argv;
}
}  