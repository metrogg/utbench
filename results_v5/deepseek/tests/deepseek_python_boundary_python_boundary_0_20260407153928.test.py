import pytest
from boundary_python_boundary_0 import (
    LedgerValidationError,
    LedgerSnapshot,
    _coerce_operation,
    _normalize_operations,
    running_balance,
    first_negative_index,
    below_zero,
    minimum_balance,
    overdraft_count,
)


class TestCoerceOperation:
    def test_coerce_operation_int(self):
        assert _coerce_operation(5, index=0) == 5
        assert _coerce_operation(-3, index=1) == -3
        assert _coerce_operation(0, index=2) == 0

    def test_coerce_operation_bool_raises(self):
        with pytest.raises(LedgerValidationError) as exc:
            _coerce_operation(True, index=0)
        assert "operation[0] must be int, got bool" in str(exc.value)
        with pytest.raises(LedgerValidationError) as exc:
            _coerce_operation(False, index=1)
        assert "operation[1] must be int, got bool" in str(exc.value)

    def test_coerce_operation_non_int_raises(self):
        with pytest.raises(LedgerValidationError) as exc:
            _coerce_operation("5", index=0)
        assert "operation[0] must be int, got str" in str(exc.value)
        with pytest.raises(LedgerValidationError) as exc:
            _coerce_operation(3.14, index=1)
        assert "operation[1] must be int, got float" in str(exc.value)
        with pytest.raises(LedgerValidationError) as exc:
            _coerce_operation(None, index=2)
        assert "operation[2] must be int, got NoneType" in str(exc.value)


class TestNormalizeOperations:
    def test_normalize_operations_empty(self):
        assert _normalize_operations([]) == []
        assert _normalize_operations(()) == []

    def test_normalize_operations_valid_int(self):
        assert _normalize_operations([1, -2, 3]) == [1, -2, 3]
        assert _normalize_operations((0, 5, -1)) == [0, 5, -1]

    def test_normalize_operations_mixed_valid(self):
        assert _normalize_operations([10, 20, 30]) == [10, 20, 30]

    def test_normalize_operations_with_bool_raises(self):
        with pytest.raises(LedgerValidationError) as exc:
            _normalize_operations([True, 2])
        assert "operation[0] must be int, got bool" in str(exc.value)

    def test_normalize_operations_with_non_int_raises(self):
        with pytest.raises(LedgerValidationError) as exc:
            _normalize_operations([1, "2", 3])
        assert "operation[1] must be int, got str" in str(exc.value)

    def test_normalize_operations_length_boundary(self):
        ops = list(range(10000))
        result = _normalize_operations(ops)
        assert len(result) == 10000
        assert result == ops

    def test_normalize_operations_too_many_raises(self):
        ops = list(range(10001))
        with pytest.raises(LedgerValidationError) as exc:
            _normalize_operations(ops)
        assert "too many operations" in str(exc.value)


class TestRunningBalance:
    def test_running_balance_empty_operations(self):
        result = running_balance([])
        assert result == []
        result = running_balance([], initial_balance=100)
        assert result == []

    def test_running_balance_single_operation(self):
        result = running_balance([5])
        assert len(result) == 1
        assert result[0].index == 0
        assert result[0].operation == 5
        assert result[0].balance == 5

        result = running_balance([-3], initial_balance=10)
        assert len(result) == 1
        assert result[0].index == 0
        assert result[0].operation == -3
        assert result[0].balance == 7

    def test_running_balance_multiple_operations(self):
        result = running_balance([1, -2, 3], initial_balance=10)
        assert len(result) == 3
        assert result[0].index == 0
        assert result[0].operation == 1
        assert result[0].balance == 11
        assert result[1].index == 1
        assert result[1].operation == -2
        assert result[1].balance == 9
        assert result[2].index == 2
        assert result[2].operation == 3
        assert result[2].balance == 12

    def test_running_balance_initial_balance_bool_raises(self):
        with pytest.raises(LedgerValidationError) as exc:
            running_balance([], initial_balance=True)
        assert "initial_balance must be int" in str(exc.value)
        with pytest.raises(LedgerValidationError) as exc:
            running_balance([], initial_balance=False)
        assert "initial_balance must be int" in str(exc.value)

    def test_running_balance_initial_balance_non_int_raises(self):
        with pytest.raises(LedgerValidationError) as exc:
            running_balance([], initial_balance="100")
        assert "initial_balance must be int" in str(exc.value)

    def test_running_balance_propagates_normalize_errors(self):
        with pytest.raises(LedgerValidationError) as exc:
            running_balance([1, "2"])
        assert "operation[1] must be int, got str" in str(exc.value)

        with pytest.raises(LedgerValidationError) as exc:
            running_balance(list(range(10001)))
        assert "too many operations" in str(exc.value)


class TestFirstNegativeIndex:
    def test_first_negative_index_no_negative(self):
        assert first_negative_index([1, 2, 3]) is None
        assert first_negative_index([], initial_balance=100) is None
        assert first_negative_index([5, 3], initial_balance=10) is None

    def test_first_negative_index_first_operation_negative(self):
        assert first_negative_index([-1], initial_balance=0) == 0
        assert first_negative_index([-5], initial_balance=4) == 0

    def test_first_negative_index_later_negative(self):
        assert first_negative_index([1, -2, 3], initial_balance=0) == 1
        assert first_negative_index([10, -5, -6, 2], initial_balance=0) == 1

    def test_first_negative_index_boundary_zero_balance(self):
        assert first_negative_index([0, 0, 0]) is None
        assert first_negative_index([1, -1, 2], initial_balance=0) is None

    def test_first_negative_index_with_initial_balance(self):
        assert first_negative_index([], initial_balance=-1) == 0
        assert first_negative_index([5], initial_balance=-10) == 0
        assert first_negative_index([3, -5], initial_balance=1) == 1

    def test_first_negative_index_propagates_validation_errors(self):
        with pytest.raises(LedgerValidationError):
            first_negative_index([True])
        with pytest.raises(LedgerValidationError):
            first_negative_index([], initial_balance="0")


class TestBelowZero:
    def test_below_zero_false(self):
        assert below_zero([1, 2, 3]) is False
        assert below_zero([0, 0, 0]) is False
        assert below_zero([5, -3, 2]) is False

    def test_below_zero_true(self):
        assert below_zero([1, 2, -4, 5]) is True
        assert below_zero([-1]) is True
        assert below_zero([5, -10]) is True

    def test_below_zero_empty(self):
        assert below_zero([]) is False

    def test_below_zero_boundary_exact_zero(self):
        assert below_zero([1, -1, 2]) is False


class TestMinimumBalance:
    def test_minimum_balance_empty_operations(self):
        assert minimum_balance([]) == 0
        assert minimum_balance([], initial_balance=100) == 100
        assert minimum_balance([], initial_balance=-50) == -50

    def test_minimum_balance_all_positive(self):
        assert minimum_balance([1, 2, 3]) == 0
        assert minimum_balance([5, 10], initial_balance=20) == 20

    def test_minimum_balance_negative_dip(self):
        assert minimum_balance([1, -2, 3]) == -1
        assert minimum_balance([10, -15, 5], initial_balance=0) == -5
        assert minimum_balance([-5, 3, -10, 2], initial_balance=10) == 5

    def test_minimum_balance_initial_negative(self):
        assert minimum_balance([], initial_balance=-5) == -5
        assert minimum_balance([2, 3], initial_balance=-10) == -10
        assert minimum_balance([-3, 1], initial_balance=-5) == -8

    def test_minimum_balance_boundary_equal_not_update(self):
        assert minimum_balance([0, 0, 0], initial_balance=5) == 5
        assert minimum_balance([-1, 2, -1], initial_balance=0) == -1

    def test_minimum_balance_propagates_validation_errors(self):
        with pytest.raises(LedgerValidationError):
            minimum_balance([True])
        with pytest.raises(LedgerValidationError):
            minimum_balance([], initial_balance="0")


class TestOverdraftCount:
    def test_overdraft_count_empty(self):
        assert overdraft_count([]) == 0
        assert overdraft_count([], initial_balance=100) == 0
        assert overdraft_count([], initial_balance=-10) == 1

    def test_overdraft_count_no_overdraft(self):
        assert overdraft_count([1, 2, 3]) == 0
        assert overdraft_count([5, -3], initial_balance=10) == 0

    def test_overdraft_count_single_overdraft(self):
        assert overdraft_count([-1]) == 1
        assert overdraft_count([1, -2, 3], initial_balance=0) == 1

    def test_overdraft_count_multiple_overdrafts(self):
        assert overdraft_count([-1, 2, -3, 4, -5]) == 3
        assert overdraft_count([-5, 3, -2, 1], initial_balance=0) == 2

    def test_overdraft_count_initial_negative(self):
        assert overdraft_count([], initial_balance=-1) == 1
        assert overdraft_count([2], initial_balance=-5) == 1
        assert overdraft_count([-3, 1], initial_balance=-10) == 2

    def test_overdraft_count_boundary_zero_balance_not_overdraft(self):
        assert overdraft_count([0, 0, 0]) == 0
        assert overdraft_count([1, -1, 2]) == 0

    def test_overdraft_count_propagates_validation_errors(self):
        with pytest.raises(LedgerValidationError):
            overdraft_count([True])
        with pytest.raises(LedgerValidationError):
            overdraft_count([], initial_balance="0")