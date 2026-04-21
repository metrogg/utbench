class RangeValidator {
    public static int clamp(int value, int minVal, int maxVal) {
        if (minVal > maxVal) {
            throw new IllegalArgumentException("minVal must be less than or equal to maxVal");
        }
        if (value < minVal) {
            return minVal;
        }
        if (value > maxVal) {
            return maxVal;
        }
        return value;
    }
}