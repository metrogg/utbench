#include <vector>
#include <string>
#include <map>
#include <sstream>
#include <algorithm>

using namespace std;

class ErrorMessageGenerator {
private:
    static string joinStrings(const vector<string>& strings, const string& delimiter) {
        if (strings.empty()) return "";
        string result;
        for (size_t i = 0; i < strings.size() - 1; ++i) {
            result += strings[i] + delimiter;
        }
        return result + strings.back();
    }

    static string formatPosition(int row, int col) {
        return "[" + to_string(row) + ":" + to_string(col) + "]";
    }

public:
    // Basic error messages
    static const map<string, string> BASIC_ERRORS;

    // Error builders
    static string buildSortError(const string& sortName, const string& errorType, 
                               int rowLeft = -1, int colLeft = -1, 
                               int rowRight = -1, int colRight = -1) {
        stringstream ss;
        ss << "Sort error: " << sortName << " - " << errorType;
        if (rowLeft != -1) {
            ss << " at position " << formatPosition(rowLeft, colLeft);
            if (rowRight != -1) {
                ss << "-" << formatPosition(rowRight, colRight);
            }
        }
        return ss.str();
    }

    static string buildFunctionError(const string& funcName, const string& errorType,
                                   const vector<string>& params = {},
                                   const string& retType = "",
                                   int rowLeft = -1, int colLeft = -1,
                                   int rowRight = -1, int colRight = -1) {
        stringstream ss;
        ss << "Function error: " << funcName << " - " << errorType;
        
        if (!params.empty()) {
            ss << " (params: " << joinStrings(params, ", ") << ")";
        }
        
        if (!retType.empty()) {
            ss << " (return type: " << retType << ")";
        }
        
        if (rowLeft != -1) {
            ss << " at position " << formatPosition(rowLeft, colLeft);
            if (rowRight != -1) {
                ss << "-" << formatPosition(rowRight, colRight);
            }
        }
        
        return ss.str();
    }

    static string buildQuantifierError(const string& quantType, const string& errorDesc,
                                     const string& expectedSort, const string& actualSort,
                                     int rowLeft, int colLeft, int rowRight, int colRight) {
        stringstream ss;
        ss << "Quantifier error: " << quantType << " - " << errorDesc
           << " (expected: " << expectedSort << ", actual: " << actualSort << ")"
           << " at position " << formatPosition(rowLeft, colLeft)
           << "-" << formatPosition(rowRight, colRight);
        return ss.str();
    }

    static string buildContextError(const string& context, const string& errorDesc,
                                  const vector<string>& additionalInfo = {}) {
        stringstream ss;
        ss << "Context error in " << context << ": " << errorDesc;
        if (!additionalInfo.empty()) {
            ss << " (" << joinStrings(additionalInfo, "; ") << ")";
        }
        return ss.str();
    }

    static string getBasicError(const string& errorKey) {
        auto it = BASIC_ERRORS.find(errorKey);
        if (it != BASIC_ERRORS.end()) {
            return it->second;
        }
        return "Unknown error: " + errorKey;
    }
};

// Initialize static const map
const map<string, string> ErrorMessageGenerator::BASIC_ERRORS = {
    {"ERR_NULL_NODE_VISIT", "Attempted to visit null AST node"},
    {"ERR_SYMBOL_MALFORMED", "Malformed symbol encountered"},
    {"ERR_ASSERT_MISSING_TERM", "Assert command missing term"},
    {"ERR_DECL_CONST_MISSING_NAME", "Constant declaration missing name"},
    {"ERR_DECL_CONST_MISSING_SORT", "Constant declaration missing sort"},
    {"ERR_DECL_FUN_MISSING_NAME", "Function declaration missing name"},
    {"ERR_DECL_FUN_MISSING_RET", "Function declaration missing return sort"},
    {"ERR_SET_LOGIC_MISSING_LOGIC", "Set-logic command missing logic name"},
    {"ERR_QUAL_ID_MISSING_ID", "Qualified identifier missing identifier"},
    {"ERR_QUAL_ID_MISSING_SORT", "Qualified identifier missing sort"},
    {"ERR_SORT_MISSING_ID", "Sort missing identifier"},
    {"ERR_QUAL_TERM_MISSING_ID", "Qualified term missing identifier"}
};
