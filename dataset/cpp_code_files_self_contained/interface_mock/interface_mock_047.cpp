#include <cstdint>
#include <functional>
#include <iostream>
#include <iomanip>
#include <map>

using namespace std;

enum CommandType {
    COMMAND_PLAYER_POS = 0x01,
    COMMAND_GET_CHUNK = 0x02,
    COMMAND_PLACE_BLOCK_ABSOLUTE = 0x03,
    COMMAND_DESTROY_BLOCK_ABSOLUTE = 0x04,
    COMMAND_MULTI_ACTION = 0x05
};

class PacketBuilder {
private:
    vector<uint8_t> buffer;
    uint8_t calculateChecksum(const vector<uint8_t>& data) {
        uint8_t checksum = 0;
        for (auto byte : data) {
            checksum ^= byte;
        }
        return checksum;
    }

public:
    // Serialize signed integers of various sizes
    void serializeSigned(int64_t value, uint8_t size, function<void(uint8_t)> callback) {
        for (int i = 0; i < size; ++i) {
            callback(static_cast<uint8_t>((value >> (8 * i)) & 0xFF));
        }
    }

    // Serialize unsigned integers of various sizes
    void serializeUnsigned(uint64_t value, uint8_t size, function<void(uint8_t)> callback) {
        for (int i = 0; i < size; ++i) {
            callback(static_cast<uint8_t>((value >> (8 * i)) & 0xFF));
        }
    }

    // Build a player position packet
    vector<uint8_t> buildPlayerPosPacket(int playerX, int playerY) {
        buffer.clear();
        uint8_t playerXSize = sizeof(playerX);
        uint8_t playerYSize = sizeof(playerY);
        
        buffer.push_back(static_cast<uint8_t>(COMMAND_PLAYER_POS));
        buffer.push_back(playerXSize);
        serializeSigned(playerX, playerXSize, [this](uint8_t b) { buffer.push_back(b); });
        buffer.push_back(playerYSize);
        serializeSigned(playerY, playerYSize, [this](uint8_t b) { buffer.push_back(b); });
        
        return finalizePacket();
    }

    // Build a chunk request packet
    vector<uint8_t> buildChunkRequestPacket(int64_t chunkX, int64_t chunkY) {
        buffer.clear();
        uint8_t chunkXSize = sizeof(chunkX);
        uint8_t chunkYSize = sizeof(chunkY);
        
        buffer.push_back(static_cast<uint8_t>(COMMAND_GET_CHUNK));
        buffer.push_back(chunkXSize);
        serializeSigned(chunkX, chunkXSize, [this](uint8_t b) { buffer.push_back(b); });
        buffer.push_back(chunkYSize);
        serializeSigned(chunkY, chunkYSize, [this](uint8_t b) { buffer.push_back(b); });
        
        return finalizePacket();
    }

    // Build a multi-action packet
    vector<uint8_t> buildMultiActionPacket(const vector<map<string, int64_t>>& actions) {
        buffer.clear();
        buffer.push_back(static_cast<uint8_t>(COMMAND_MULTI_ACTION));
        buffer.push_back(static_cast<uint8_t>(actions.size()));
        
        for (const auto& action : actions) {
            uint8_t actionType = action.at("type");
            buffer.push_back(actionType);
            
            switch (actionType) {
                case COMMAND_PLAYER_POS: {
                    int x = action.at("x");
                    int y = action.at("y");
                    uint8_t xSize = sizeof(x);
                    uint8_t ySize = sizeof(y);
                    buffer.push_back(xSize);
                    serializeSigned(x, xSize, [this](uint8_t b) { buffer.push_back(b); });
                    buffer.push_back(ySize);
                    serializeSigned(y, ySize, [this](uint8_t b) { buffer.push_back(b); });
                    break;
                }
                case COMMAND_PLACE_BLOCK_ABSOLUTE: {
                    uint32_t id = action.at("id");
                    int64_t chunkX = action.at("chunkX");
                    int64_t chunkY = action.at("chunkY");
                    uint16_t x = action.at("x");
                    uint16_t y = action.at("y");
                    
                    uint8_t idSize = sizeof(id);
                    uint8_t chunkXSize = sizeof(chunkX);
                    uint8_t chunkYSize = sizeof(chunkY);
                    uint8_t xSize = sizeof(x);
                    uint8_t ySize = sizeof(y);
                    
                    buffer.push_back(idSize);
                    serializeUnsigned(id, idSize, [this](uint8_t b) { buffer.push_back(b); });
                    buffer.push_back(chunkXSize);
                    serializeSigned(chunkX, chunkXSize, [this](uint8_t b) { buffer.push_back(b); });
                    buffer.push_back(chunkYSize);
                    serializeSigned(chunkY, chunkYSize, [this](uint8_t b) { buffer.push_back(b); });
                    buffer.push_back(xSize);
                    serializeSigned(x, xSize, [this](uint8_t b) { buffer.push_back(b); });
                    buffer.push_back(ySize);
                    serializeSigned(y, ySize, [this](uint8_t b) { buffer.push_back(b); });
                    break;
                }
            }
        }
        
        return finalizePacket();
    }

private:
    vector<uint8_t> finalizePacket() {
        uint8_t totalSize = buffer.size() + 2; // +2 for size and checksum
        vector<uint8_t> finalPacket;
        finalPacket.push_back(totalSize);
        finalPacket.insert(finalPacket.end(), buffer.begin(), buffer.end());
        finalPacket.push_back(calculateChecksum(finalPacket));
        return finalPacket;
    }
};
