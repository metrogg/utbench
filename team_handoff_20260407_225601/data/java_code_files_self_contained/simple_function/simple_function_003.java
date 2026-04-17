import java.util.*;
import java.lang.*;

class Solution {
    /**
    Given two lists operator, and operand. The first list has basic algebra operations, and
    the second list is a list of integers. Use the two given lists to build the algebric
    expression and return the evaluation of this expression.

    The basic algebra operations:
    Addition ( + )
    Subtraction ( - )
    Multiplication ( * )
    Floor division ( / )
    Exponentiation ( ** )

    Example:
    operator["+", "*", "-"]
    array = [2, 3, 4, 5]
    result = 2 + 3 * 4 - 5
    => result = 9

    Note:
        The length of operator list is equal to the length of operand list minus one.
        Operand is a list of of non-negative integers.
        Operator list has at least one operator, and operand list has at least two operands.
     */
    public int doAlgebra(List<String> operator, List<Integer> operand) {
        List<String> ops = new ArrayList<>(operator);
        List<Integer> nums = new ArrayList<>(operand);
        for (int i = ops.size() - 1; i >= 0; i--) {
            if (ops.get(i).equals("**")) {
                nums.set(i, (int) Math.round(Math.pow(nums.get(i), nums.get(i + 1))));
                nums.remove(i + 1);
                ops.remove(i);
            }
        }
        for (int i = 0; i < ops.size(); i++) {
            if (ops.get(i).equals("*")) {
                nums.set(i, nums.get(i) * nums.get(i + 1));
                nums.remove(i + 1);
                ops.remove(i);
                i--;
            } else if (ops.get(i).equals("/")) {
                nums.set(i, nums.get(i) / nums.get(i + 1));
                nums.remove(i + 1);
                ops.remove(i);
                i--;
            }
        }
        for (int i = 0; i < ops.size(); i++) {
            if (ops.get(i).equals("+")) {
                nums.set(i, nums.get(i) + nums.get(i + 1));
                nums.remove(i + 1);
                ops.remove(i);
                i--;
            } else if (ops.get(i).equals("-")) {
                nums.set(i, nums.get(i) - nums.get(i + 1));
                nums.remove(i + 1);
                ops.remove(i);
                i--;
            }
        }
        return nums.get(0);
    }
}
