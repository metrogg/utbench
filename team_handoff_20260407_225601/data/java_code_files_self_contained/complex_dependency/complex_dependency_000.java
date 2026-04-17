// Converted Java method
import java.util.Objects;

class EnhancedSpotLight {
    private PointLight pointLight;
    private Vector3 direction;
    private float cutoff;
    private float outerCutoff;
    private float falloffExponent;
    private boolean castsShadows;

    /**
     * Enhanced spot light with additional physical properties
     * @param pointLight The base point light properties
     * @param direction Direction of the spotlight (will be normalized)
     * @param cutoff Inner cutoff angle in degrees
     * @param outerCutoff Outer cutoff angle in degrees
     * @param falloffExponent Light falloff exponent (1.0 for linear)
     * @param castsShadows Whether this light casts shadows
     */
    public EnhancedSpotLight(PointLight pointLight, Vector3 direction, 
                           float cutoff, float outerCutoff, 
                           float falloffExponent, boolean castsShadows) {
        setPointLight(pointLight);
        setDirection(direction);
        setCutoff(cutoff);
        setOuterCutoff(outerCutoff);
        setFalloffExponent(falloffExponent);
        this.castsShadows = castsShadows;
    }

    // Calculate light intensity at a given point
    public float calculateIntensity(Vector3 point) {
        Vector3 lightToPoint = point.sub(pointLight.getPosition()).normalized();
        float theta = direction.dot(lightToPoint);
        float epsilon = (float)Math.cos(Math.toRadians(cutoff)) - 
                       (float)Math.cos(Math.toRadians(outerCutoff));
        float intensity = (theta - (float)Math.cos(Math.toRadians(outerCutoff))) / epsilon;
        intensity = (float)Math.pow(Math.max(0, intensity), falloffExponent);
        
        // Apply distance attenuation
        float distance = pointLight.getPosition().sub(point).length();
        float attenuation = 1.0f / (pointLight.getConstant() + 
                                   pointLight.getLinear() * distance + 
                                   pointLight.getExponent() * distance * distance);
        
        return intensity * attenuation * pointLight.getIntensity();
    }

    // Getters and setters with validation
    public PointLight getPointLight() {
        return pointLight;
    }

    public void setPointLight(PointLight pointLight) {
        this.pointLight = Objects.requireNonNull(pointLight, "PointLight cannot be null");
    }

    public Vector3 getDirection() {
        return direction;
    }

    public void setDirection(Vector3 direction) {
        this.direction = Objects.requireNonNull(direction, "Direction cannot be null").normalized();
    }

    public float getCutoff() {
        return cutoff;
    }

    public void setCutoff(float cutoff) {
        if (cutoff <= 0 || cutoff >= 90) {
            throw new IllegalArgumentException("Cutoff must be between 0 and 90 degrees");
        }
        this.cutoff = cutoff;
    }

    public float getOuterCutoff() {
        return outerCutoff;
    }

    public void setOuterCutoff(float outerCutoff) {
        if (outerCutoff <= 0 || outerCutoff >= 90) {
            throw new IllegalArgumentException("Outer cutoff must be between 0 and 90 degrees");
        }
        if (outerCutoff <= cutoff) {
            throw new IllegalArgumentException("Outer cutoff must be greater than inner cutoff");
        }
        this.outerCutoff = outerCutoff;
    }

    public float getFalloffExponent() {
        return falloffExponent;
    }

    public void setFalloffExponent(float falloffExponent) {
        if (falloffExponent <= 0) {
            throw new IllegalArgumentException("Falloff exponent must be positive");
        }
        this.falloffExponent = falloffExponent;
    }

    public boolean isCastsShadows() {
        return castsShadows;
    }

    public void setCastsShadows(boolean castsShadows) {
        this.castsShadows = castsShadows;
    }
}

// Supporting classes
class PointLight {
    private Vector3 position;
    private float intensity;
    private float constant;
    private float linear;
    private float exponent;

    public PointLight(Vector3 position, float intensity, float constant, float linear, float exponent) {
        this.position = position;
        this.intensity = intensity;
        this.constant = constant;
        this.linear = linear;
        this.exponent = exponent;
    }

    // Getters and setters
    public Vector3 getPosition() { return position; }
    public float getIntensity() { return intensity; }
    public float getConstant() { return constant; }
    public float getLinear() { return linear; }
    public float getExponent() { return exponent; }
}

class Vector3 {
    private float x, y, z;
    
    public Vector3(float x, float y, float z) {
        this.x = x;
        this.y = y;
        this.z = z;
    }
    
    public Vector3 normalized() {
        float length = length();
        return new Vector3(x/length, y/length, z/length);
    }
    
    public float length() {
        return (float)Math.sqrt(x*x + y*y + z*z);
    }
    
    public float dot(Vector3 other) {
        return x*other.x + y*other.y + z*other.z;
    }
    
    public Vector3 sub(Vector3 other) {
        return new Vector3(x-other.x, y-other.y, z-other.z);
    }
    
    // Getters
    public float getX() { return x; }
    public float getY() { return y; }
    public float getZ() { return z; }
}
