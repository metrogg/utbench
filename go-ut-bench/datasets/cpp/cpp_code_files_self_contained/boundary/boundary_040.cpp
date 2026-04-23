#include <vector>
#include <map>
#include <string>
#include <stdexcept>

// Simulated CUDA functions and types for self-contained example
enum cudaError {
    cudaSuccess,
    cudaErrorNoDevice,
    cudaErrorInsufficientDriver,
    cudaErrorInvalidDevice
};

struct cudaDeviceProp {
    std::string name;
    int major;
    int minor;
    size_t totalGlobalMem;
    size_t sharedMemPerBlock;
    int regsPerBlock;
    int warpSize;
    size_t memPitch;
    int maxThreadsPerBlock;
    int maxThreadsDim[3];
    int maxGridSize[3];
    int clockRate;
    size_t totalConstMem;
    int computeMode;
    int deviceOverlap;
    int multiProcessorCount;
    int kernelExecTimeoutEnabled;
    int integrated;
    int canMapHostMemory;
    int concurrentKernels;
    int eccEnabled;
    int pciBusID;
    int pciDeviceID;
    int tccDriver;
    int unifiedAddressing;
};

// Mock CUDA functions for demonstration
cudaError cudaGetDeviceCount(int* count) {
    // Simulate having 1 device for testing
    *count = 1;
    return cudaSuccess;
}

cudaError cudaGetDeviceProperties(cudaDeviceProp* prop, int device) {
    if (device != 0) return cudaErrorInvalidDevice;
    
    prop->name = "Simulated CUDA Device";
    prop->major = 7;
    prop->minor = 5;
    prop->totalGlobalMem = 8ULL * 1024 * 1024 * 1024; // 8GB
    prop->sharedMemPerBlock = 48 * 1024;
    prop->regsPerBlock = 65536;
    prop->warpSize = 32;
    prop->maxThreadsPerBlock = 1024;
    prop->maxThreadsDim[0] = 1024;
    prop->maxThreadsDim[1] = 1024;
    prop->maxThreadsDim[2] = 64;
    prop->maxGridSize[0] = 2147483647;
    prop->maxGridSize[1] = 65535;
    prop->maxGridSize[2] = 65535;
    prop->clockRate = 1545000;
    prop->multiProcessorCount = 40;
    
    return cudaSuccess;
}

cudaError cudaMemGetInfo(size_t* free, size_t* total) {
    *total = 8ULL * 1024 * 1024 * 1024; // 8GB
    *free = 6ULL * 1024 * 1024 * 1024;  // 6GB free
    return cudaSuccess;
}

int findCudaDevice(int argc, char** argv) {
    // Simple mock that always returns device 0
    return 0;
}

// Enhanced system capability checker
std::map<std::string, std::string> checkSystemCapabilities(bool requireAtomic = true, int minComputeMajor = 2, int minComputeMinor = 0) {
    std::map<std::string, std::string> results;
    
    int device_count = 0;
    cudaError t = cudaGetDeviceCount(&device_count);
    
    if (t != cudaSuccess) {
        results["status"] = "error";
        results["message"] = "First call to CUDA Runtime API failed. Are the drivers installed?";
        return results;
    }

    if (device_count < 1) {
        results["status"] = "error";
        results["message"] = "No CUDA devices found. Make sure CUDA device is powered, connected and available.";
        return results;
    }

    results["status"] = "success";
    results["device_count"] = std::to_string(device_count);
    
    int device = findCudaDevice(0, 0);
    cudaDeviceProp properties;
    cudaGetDeviceProperties(&properties, device);
    
    results["device_name"] = properties.name;
    results["compute_capability"] = std::to_string(properties.major) + "." + std::to_string(properties.minor);
    results["multi_processor_count"] = std::to_string(properties.multiProcessorCount);
    results["max_threads_per_block"] = std::to_string(properties.maxThreadsPerBlock);
    
    size_t free, total;
    cudaMemGetInfo(&free, &total);
    results["total_memory_mb"] = std::to_string(total >> 20);
    results["free_memory_mb"] = std::to_string(free >> 20);
    
    if (requireAtomic && (properties.major < minComputeMajor || 
                         (properties.major == minComputeMajor && properties.minor < minComputeMinor))) {
        results["status"] = "warning";
        results["compatibility"] = "Your CUDA device has compute capability " + 
                                  results["compute_capability"] + 
                                  ". The minimum required is " + 
                                  std::to_string(minComputeMajor) + "." + 
                                  std::to_string(minComputeMinor) + " for atomic operations.";
    } else {
        results["compatibility"] = "Device meets minimum requirements";
    }
    
    return results;
}
