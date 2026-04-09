#include <string>
#include <vector>
#include <map>
#include <memory>
#include <stdexcept>

using namespace std;

// Enhanced Media Player Interface with additional functionality
class MediaPlayer {
public:
    virtual void play(const string& fileName) = 0;
    virtual void stop() = 0;
    virtual void pause() = 0;
    virtual void setVolume(int volume) = 0;
    virtual string getCurrentStatus() const = 0;
    virtual ~MediaPlayer() = default;
};

// Advanced Media Player Interface with more formats
class AdvancedMediaPlayer {
public:
    virtual void playVlc(const string& fileName) = 0;
    virtual void playMp4(const string& fileName) = 0;
    virtual void playFlac(const string& fileName) = 0;
    virtual void playWav(const string& fileName) = 0;
    virtual ~AdvancedMediaPlayer() = default;
};

// Concrete implementations of AdvancedMediaPlayer
class VlcPlayer : public AdvancedMediaPlayer {
public:
    void playVlc(const string& fileName) override {
        cout << "Playing vlc file: " << fileName << endl;
    }
    void playMp4(const string&) override {}
    void playFlac(const string&) override {}
    void playWav(const string&) override {}
};

class Mp4Player : public AdvancedMediaPlayer {
public:
    void playVlc(const string&) override {}
    void playMp4(const string& fileName) override {
        cout << "Playing mp4 file: " << fileName << endl;
    }
    void playFlac(const string&) override {}
    void playWav(const string&) override {}
};

class FlacPlayer : public AdvancedMediaPlayer {
public:
    void playVlc(const string&) override {}
    void playMp4(const string&) override {}
    void playFlac(const string& fileName) override {
        cout << "Playing flac file: " << fileName << endl;
    }
    void playWav(const string&) override {}
};

class WavPlayer : public AdvancedMediaPlayer {
public:
    void playVlc(const string&) override {}
    void playMp4(const string&) override {}
    void playFlac(const string&) override {}
    void playWav(const string& fileName) override {
        cout << "Playing wav file: " << fileName << endl;
    }
};

// Enhanced MediaAdapter with state management
class MediaAdapter : public MediaPlayer {
private:
    unique_ptr<AdvancedMediaPlayer> advancedPlayer;
    int volume = 50;
    bool isPlaying = false;
    bool isPaused = false;
    string currentFile;

public:
    explicit MediaAdapter(const string& audioType) {
        if (audioType == "vlc") {
            advancedPlayer = make_unique<VlcPlayer>();
        } else if (audioType == "mp4") {
            advancedPlayer = make_unique<Mp4Player>();
        } else if (audioType == "flac") {
            advancedPlayer = make_unique<FlacPlayer>();
        } else if (audioType == "wav") {
            advancedPlayer = make_unique<WavPlayer>();
        } else {
            throw invalid_argument("Unsupported media type: " + audioType);
        }
    }

    void play(const string& fileName) override {
        string ext = getFileExtension(fileName);
        
        if (ext == "vlc") {
            advancedPlayer->playVlc(fileName);
        } else if (ext == "mp4") {
            advancedPlayer->playMp4(fileName);
        } else if (ext == "flac") {
            advancedPlayer->playFlac(fileName);
        } else if (ext == "wav") {
            advancedPlayer->playWav(fileName);
        } else {
            throw invalid_argument("Unsupported file format: " + ext);
        }
        
        isPlaying = true;
        isPaused = false;
        currentFile = fileName;
        cout << "Volume set to: " << volume << "%" << endl;
    }

    void stop() override {
        if (isPlaying || isPaused) {
            cout << "Stopped playing: " << currentFile << endl;
            isPlaying = false;
            isPaused = false;
            currentFile = "";
        }
    }

    void pause() override {
        if (isPlaying) {
            cout << "Paused: " << currentFile << endl;
            isPaused = true;
            isPlaying = false;
        }
    }

    void setVolume(int vol) override {
        if (vol < 0) vol = 0;
        if (vol > 100) vol = 100;
        volume = vol;
        cout << "Volume changed to: " << volume << "%" << endl;
    }

    string getCurrentStatus() const override {
        if (isPlaying) return "Playing: " + currentFile;
        if (isPaused) return "Paused: " + currentFile;
        return "Stopped";
    }

private:
    string getFileExtension(const string& fileName) {
        size_t dotPos = fileName.find_last_of('.');
        if (dotPos == string::npos) return "";
        return fileName.substr(dotPos + 1);
    }
};

// AudioPlayer with playlist functionality
class AudioPlayer : public MediaPlayer {
private:
    unique_ptr<MediaAdapter> mediaAdapter;
    vector<string> playlist;
    size_t currentTrack = 0;
    int volume = 50;
    bool isPlaying = false;
    bool isPaused = false;

public:
    void addToPlaylist(const string& fileName) {
        playlist.push_back(fileName);
    }

    void clearPlaylist() {
        playlist.clear();
        currentTrack = 0;
        isPlaying = false;
        isPaused = false;
    }

    void play(const string& fileName) override {
        string ext = getFileExtension(fileName);
        
        if (ext == "mp3") {
            cout << "Playing mp3 file: " << fileName << endl;
            cout << "Volume set to: " << volume << "%" << endl;
            isPlaying = true;
            isPaused = false;
        } else if (ext == "vlc" || ext == "mp4" || ext == "flac" || ext == "wav") {
            try {
                mediaAdapter = make_unique<MediaAdapter>(ext);
                mediaAdapter->setVolume(volume);
                mediaAdapter->play(fileName);
                isPlaying = true;
                isPaused = false;
            } catch (const invalid_argument& e) {
                cerr << "Error: " << e.what() << endl;
            }
        } else {
            cerr << "Invalid media. " << ext << " format not supported" << endl;
        }
    }

    void playNext() {
        if (playlist.empty()) return;
        
        currentTrack = (currentTrack + 1) % playlist.size();
        play(playlist[currentTrack]);
    }

    void playPrevious() {
        if (playlist.empty()) return;
        
        currentTrack = (currentTrack == 0) ? playlist.size() - 1 : currentTrack - 1;
        play(playlist[currentTrack]);
    }

    void stop() override {
        if (mediaAdapter) {
            mediaAdapter->stop();
        } else if (isPlaying || isPaused) {
            cout << "Stopped playing" << endl;
        }
        isPlaying = false;
        isPaused = false;
    }

    void pause() override {
        if (mediaAdapter) {
            mediaAdapter->pause();
        } else if (isPlaying) {
            cout << "Paused" << endl;
            isPaused = true;
            isPlaying = false;
        }
    }

    void setVolume(int vol) override {
        if (vol < 0) vol = 0;
        if (vol > 100) vol = 100;
        volume = vol;
        if (mediaAdapter) {
            mediaAdapter->setVolume(vol);
        } else {
            cout << "Volume changed to: " << volume << "%" << endl;
        }
    }

    string getCurrentStatus() const override {
        if (mediaAdapter) {
            return mediaAdapter->getCurrentStatus();
        }
        if (isPlaying) return "Playing";
        if (isPaused) return "Paused";
        return "Stopped";
    }

    void showPlaylist() const {
        if (playlist.empty()) {
            cout << "Playlist is empty" << endl;
            return;
        }
        
        cout << "Current playlist:" << endl;
        for (size_t i = 0; i < playlist.size(); ++i) {
            cout << (i == currentTrack ? "> " : "  ") << i+1 << ". " << playlist[i] << endl;
        }
    }

private:
    string getFileExtension(const string& fileName) {
        size_t dotPos = fileName.find_last_of('.');
        if (dotPos == string::npos) return "";
        return fileName.substr(dotPos + 1);
    }
};
