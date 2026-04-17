#include <string>
#include <bitset>
#include <vector>
#include <algorithm>

using namespace std;

class AdvancedTokenizer {
private:
    string m_text;
    bitset<256> m_delimiters;
    size_t m_currentPos;
    bool m_keepDelimiters;
    bool m_caseSensitive;

public:
    AdvancedTokenizer(const string& text, const string& delimiters, 
                     bool keepDelimiters = false, bool caseSensitive = true)
        : m_text(text), m_currentPos(0), 
          m_keepDelimiters(keepDelimiters), m_caseSensitive(caseSensitive) {
        setDelimiters(delimiters);
    }

    void setDelimiters(const string& delimiters) {
        m_delimiters.reset();
        for (char c : delimiters) {
            if (!m_caseSensitive) {
                m_delimiters.set(tolower(c));
                m_delimiters.set(toupper(c));
            } else {
                m_delimiters.set(c);
            }
        }
    }

    void reset() {
        m_currentPos = 0;
    }

    bool nextToken(string& token) {
        token.clear();
        
        // Skip leading delimiters
        while (m_currentPos < m_text.size() && m_delimiters.test(m_text[m_currentPos])) {
            if (m_keepDelimiters) {
                token = m_text.substr(m_currentPos, 1);
                m_currentPos++;
                return true;
            }
            m_currentPos++;
        }

        if (m_currentPos >= m_text.size()) {
            return false;
        }

        size_t startPos = m_currentPos;
        while (m_currentPos < m_text.size() && !m_delimiters.test(m_text[m_currentPos])) {
            m_currentPos++;
        }

        token = m_text.substr(startPos, m_currentPos - startPos);
        return true;
    }

    vector<string> getAllTokens() {
        vector<string> tokens;
        string token;
        reset();
        while (nextToken(token)) {
            tokens.push_back(token);
        }
        return tokens;
    }

    string getRemainingText() {
        return m_currentPos < m_text.size() ? m_text.substr(m_currentPos) : "";
    }

    size_t getCurrentPosition() const {
        return m_currentPos;
    }
};
