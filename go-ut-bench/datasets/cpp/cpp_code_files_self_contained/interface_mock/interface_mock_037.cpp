#include <fstream>
#include <algorithm>
#include <vector>
#include <map>
#include <chrono>
#include <memory>
#include <stdexcept>
#include <iomanip>
#include <sstream>

using namespace std;
using namespace std::chrono;

class Symbol {
public:
    string name;
    string module;
    size_t start_address;
    size_t end_address;

    Symbol(const string& n, const string& m, size_t start, size_t end)
        : name(n), module(m), start_address(start), end_address(end) {}
};

class SymbolResolver {
private:
    string perf_map_path;
    vector<Symbol> symbols;
    system_clock::time_point last_update;
    milliseconds update_interval;
    size_t last_read_offset;

    void parse_perf_map_line(const string& line, string& module, string& name, size_t& start, size_t& size) {
        istringstream iss(line);
        string address_str, size_str;
        
        iss >> address_str >> size_str;
        
        // Parse address and size (hexadecimal)
        start = stoull(address_str, nullptr, 16);
        size = stoull(size_str, nullptr, 16);
        
        // The rest is the symbol name which may contain spaces
        getline(iss, name);
        
        // Trim whitespace
        name.erase(0, name.find_first_not_of(" \t"));
        name.erase(name.find_last_not_of(" \t") + 1);
        
        // Extract module name if present (format: [module] function)
        size_t module_start = name.find('[');
        size_t module_end = name.find(']');
        
        if (module_start != string::npos && module_end != string::npos && module_end > module_start) {
            module = name.substr(module_start + 1, module_end - module_start - 1);
            name = name.substr(module_end + 1);
            
            // Trim whitespace from function name
            name.erase(0, name.find_first_not_of(" \t"));
            name.erase(name.find_last_not_of(" \t") + 1);
        } else {
            module = "unknown";
        }
    }

public:
    SymbolResolver(const string& path, milliseconds interval = 100ms)
        : perf_map_path(path), update_interval(interval), last_read_offset(0) {
        update_symbols();
    }

    void update_symbols() {
        ifstream file(perf_map_path);
        if (!file.is_open()) {
            throw runtime_error("Could not open perf map file: " + perf_map_path);
        }

        file.seekg(last_read_offset);
        string line;
        vector<Symbol> new_symbols;

        while (getline(file, line)) {
            last_read_offset = file.tellg();
            
            string module, name;
            size_t start, size;
            
            try {
                parse_perf_map_line(line, module, name, start, size);
                new_symbols.emplace_back(name, module, start, start + size);
            } catch (const exception& e) {
                cerr << "Error parsing line: " << line << " - " << e.what() << endl;
            }
        }

        // Merge with existing symbols and sort
        symbols.insert(symbols.end(), new_symbols.begin(), new_symbols.end());
        sort(symbols.begin(), symbols.end(), 
            [](const Symbol& a, const Symbol& b) { return a.start_address < b.start_address; });

        last_update = system_clock::now();
    }

    shared_ptr<Symbol> resolve(size_t address, bool force_update = false) {
        auto now = system_clock::now();
        if (force_update || now - last_update > update_interval) {
            update_symbols();
        }

        // Binary search for the symbol
        auto it = upper_bound(symbols.begin(), symbols.end(), address,
            [](size_t addr, const Symbol& sym) { return addr < sym.end_address; });

        if (it != symbols.end() && address >= it->start_address) {
            return make_shared<Symbol>(*it);
        }

        return nullptr;
    }

    vector<shared_ptr<Symbol>> get_symbols_in_range(size_t start, size_t end) {
        vector<shared_ptr<Symbol>> result;
        
        auto lower = lower_bound(symbols.begin(), symbols.end(), start,
            [](const Symbol& sym, size_t addr) { return sym.end_address < addr; });
            
        auto upper = upper_bound(symbols.begin(), symbols.end(), end,
            [](size_t addr, const Symbol& sym) { return addr < sym.start_address; });
            
        for (auto it = lower; it != upper; ++it) {
            result.push_back(make_shared<Symbol>(*it));
        }
        
        return result;
    }

    void print_symbol_table() {
        cout << "Symbol Table (" << symbols.size() << " symbols):" << endl;
        cout << "+------------------+------------------+------------------+------------------------+" << endl;
        cout << "| Start Address    | End Address      | Size             | Module::Function       |" << endl;
        cout << "+------------------+------------------+------------------+------------------------+" << endl;
        
        for (const auto& sym : symbols) {
            cout << "| 0x" << hex << setw(14) << setfill('0') << sym.start_address << " "
                 << "| 0x" << hex << setw(14) << setfill('0') << sym.end_address << " "
                 << "| 0x" << hex << setw(14) << setfill('0') << (sym.end_address - sym.start_address) << " "
                 << "| " << left << setw(20) << setfill(' ') << (sym.module + "::" + sym.name) << " |" << endl;
        }
        
        cout << "+------------------+------------------+------------------+------------------------+" << endl;
    }
};
