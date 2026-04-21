class SafeMath {
    public static double safeDivide(double numerator, double denominator) {
        if (denominator == 0) {
            throw new ArithmeticException("Cannot divide by zero");
        }
        if (numerator == 0) {
            return 0.0;
        }
        double result = numerator / denominator;
        if (result > Double.MAX_VALUE) {
            return Double.POSITIVE_INFINITY;
        }
        if (result < -Double.MAX_VALUE) {
            return Double.NEGATIVE_INFINITY;
        }
        return result;
    }
    
    public static int safeSqrt(int value) {
        if (value < 0) {
            throw new IllegalArgumentException("Cannot compute square root of negative number");
        }
        if (value == 0) {
            return 0;
        }
        return (int) Math.sqrt(value);
    }
}