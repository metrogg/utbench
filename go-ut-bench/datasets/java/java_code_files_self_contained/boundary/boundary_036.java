import java.util.*;
import java.lang.*;

class Solution {
    /**
    Create a function that takes integers, floats, or strings representing
    real numbers, and returns the larger variable in its given variable type.
    Return None if the values are equal.
    Note: If a real number is represented as a string, the floating point might be . or ,

    compareOne(1, 2.5) -> Optional.of(2.5)
    compareOne(1, "2,3") -> Optional.of("2,3")
    compareOne("5,1", "6") -> Optional.of("6")
    compareOne("1", 1) -> Optional.empty()
     */
    public Optional<Object> compareOne(Object a, Object b) {
        double temp_a = 0, temp_b = 0;
        if (a instanceof Integer) {
            temp_a = (Integer) a * 1.0;
        } else if (a instanceof Double) {
            temp_a = (double) a;
        } else if (a instanceof String) {
            temp_a = Double.parseDouble(((String) a).replace(',', '.'));
        }
        if (b instanceof Integer) {
            temp_b = (Integer) b * 1.0;
        } else if (b instanceof Double) {
            temp_b = (double) b;
        } else if (b instanceof String) {
            temp_b = Double.parseDouble(((String) b).replace(',', '.'));
        }
        if (temp_a == temp_b) {
            return Optional.empty();
        } else if (temp_a > temp_b) {
            return Optional.of(a);
        } else {
            return Optional.of(b);
        }
    }
}
