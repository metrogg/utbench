#include <vector>
#include <string>
#include <map>
#include <cassert>

// Simulated Direct3D interfaces for our self-contained implementation
class IDirect3DDevice9 {
public:
    virtual ~IDirect3DDevice9() {}
};

class IDirect3DTexture9 {
public:
    virtual ~IDirect3DTexture9() {}
};

// Simulated D3DX functions
namespace D3DX {
    enum { OK = 0 };
}

// Texture manager class with enhanced functionality
class TextureManager {
private:
    IDirect3DTexture9* mTexture;
    std::string mFileName;
    unsigned int mWidth;
    unsigned int mHeight;
    unsigned int mMipLevels;
    bool mIsLoaded;
    
    // Static device pointer (simplified from original)
    static IDirect3DDevice9* mDevice;
    
public:
    // Constructor
    TextureManager() : mTexture(nullptr), mWidth(0), mHeight(0), 
                      mMipLevels(0), mIsLoaded(false) {}
    
    // Load texture from file with additional options
    bool load(const std::string& fileName, bool generateMipMaps = true, 
              unsigned int maxDimension = 0) {
        // Validate input
        if (fileName.empty() || !mDevice) {
            return false;
        }
        
        // Release existing texture
        release();
        
        // Simulated texture loading
        mTexture = new IDirect3DTexture9();
        mFileName = fileName;
        
        // Simulate getting texture info (in real D3D this would come from file)
        mWidth = 256;  // Default simulated values
        mHeight = 256;
        mMipLevels = generateMipMaps ? calculateMipLevels(mWidth, mHeight) : 1;
        
        mIsLoaded = true;
        return true;
    }
    
    // Calculate mip levels based on texture dimensions
    static unsigned int calculateMipLevels(unsigned int width, unsigned int height) {
        unsigned int levels = 1;
        while (width > 1 || height > 1) {
            width = std::max(1u, width / 2);
            height = std::max(1u, height / 2);
            levels++;
        }
        return levels;
    }
    
    // Bind texture to a specific slot
    void bind(unsigned int slot = 0) const {
        assert(mDevice != nullptr);
        assert(mIsLoaded && "Texture not loaded");
        
        // Simulated texture binding
        std::cout << "Binding texture '" << mFileName << "' to slot " << slot << std::endl;
    }
    
    // Get texture information
    void getInfo(unsigned int& width, unsigned int& height, 
                 unsigned int& mipLevels) const {
        width = mWidth;
        height = mHeight;
        mipLevels = mMipLevels;
    }
    
    // Release texture resources
    void release() {
        if (mTexture) {
            delete mTexture;
            mTexture = nullptr;
        }
        mIsLoaded = false;
        mWidth = mHeight = mMipLevels = 0;
        mFileName.clear();
    }
    
    // Check if texture is loaded
    bool isLoaded() const { return mIsLoaded; }
    
    // Get filename
    const std::string& getFileName() const { return mFileName; }
    
    // Static device management
    static void setDevice(IDirect3DDevice9* device) { mDevice = device; }
    static IDirect3DDevice9* getDevice() { return mDevice; }
    
    // Destructor
    ~TextureManager() {
        release();
    }
};

// Initialize static member
IDirect3DDevice9* TextureManager::mDevice = nullptr;
