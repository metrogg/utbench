#include <vector>
#include <string>
#include <stdexcept>
#include <memory>

using namespace std;

// Simulated GDI structures and functions for self-contained testing
struct HDC {};
struct BITMAP {
    int width;
    int height;
};

class ImageInfo {
public:
    int width;
    int height;
    int frameWidth;
    int frameHeight;
    int currentFrameX;
    int currentFrameY;
    HDC hMemDC;
    BITMAP bitmap;

    ImageInfo(int w, int h, int fw = 0, int fh = 0) 
        : width(w), height(h), frameWidth(fw), frameHeight(fh), 
          currentFrameX(0), currentFrameY(0) {
        if (fw == 0) frameWidth = w;
        if (fh == 0) frameHeight = h;
        bitmap.width = w;
        bitmap.height = h;
    }
};

class ImageRenderer {
private:
    shared_ptr<ImageInfo> _imageInfo;
    bool _trans;
    int _transColor;

public:
    ImageRenderer(int width, int height, bool transparent = false, int transColor = 0)
        : _imageInfo(make_shared<ImageInfo>(width, height)),
          _trans(transparent), _transColor(transColor) {}

    ImageRenderer(int width, int height, int frameWidth, int frameHeight, 
                 bool transparent = false, int transColor = 0)
        : _imageInfo(make_shared<ImageInfo>(width, height, frameWidth, frameHeight)),
          _trans(transparent), _transColor(transColor) {}

    // Simulated rendering functions that return rendering information instead of actual rendering
    string render(int destX, int destY, int adjWidth, int adjHeight) {
        string method = _trans ? "TransparentBlt" : "StretchBlt";
        
        return "Rendered with " + method + 
               " at (" + to_string(destX) + "," + to_string(destY) + ")" +
               " size (" + to_string(adjWidth) + "x" + to_string(adjHeight) + ")" +
               " from source (0,0)" +
               " size (" + to_string(_imageInfo->width) + "x" + to_string(_imageInfo->height) + ")";
    }

    string frameRender(int destX, int destY, int currentFrameX, int currentFrameY, int adjWidth, int adjHeight) {
        _imageInfo->currentFrameX = currentFrameX;
        _imageInfo->currentFrameY = currentFrameY;
        
        string method = _trans ? "TransparentBlt" : "StretchBlt";
        int srcX = currentFrameX * _imageInfo->frameWidth;
        int srcY = currentFrameY * _imageInfo->frameHeight;
        
        return "Rendered frame with " + method + 
               " at (" + to_string(destX) + "," + to_string(destY) + ")" +
               " size (" + to_string(adjWidth) + "x" + to_string(adjHeight) + ")" +
               " from source (" + to_string(srcX) + "," + to_string(srcY) + ")" +
               " size (" + to_string(_imageInfo->frameWidth) + "x" + to_string(_imageInfo->frameHeight) + ")";
    }

    // Additional functionality: image transformation
    string transformAndRender(int destX, int destY, int adjWidth, int adjHeight, 
                            float rotation = 0.0f, float scaleX = 1.0f, float scaleY = 1.0f) {
        if (rotation != 0.0f || scaleX != 1.0f || scaleY != 1.0f) {
            return "Transformed image with rotation=" + to_string(rotation) + 
                   " scale=(" + to_string(scaleX) + "," + to_string(scaleY) + ")" +
                   " then " + render(destX, destY, adjWidth, adjHeight);
        }
        return render(destX, destY, adjWidth, adjHeight);
    }
};
