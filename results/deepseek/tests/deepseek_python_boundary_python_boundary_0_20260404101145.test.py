import pytest
from typing import List


def test_below_zero_normal_positive_operations():
    """正常路径：所有操作都是正数，余额永远不会低于0"""
    operations = [1, 2, 3, 4, 5]
    result = below_zero(operations)
    assert result is False


def test_below_zero_normal_negative_at_end():
    """正常路径：操作序列中最后一步使余额低于0"""
    operations = [10, 20, -35]
    result = below_zero(operations)
    assert result is True


def test_below_zero_normal_negative_in_middle():
    """正常路径：操作序列中间某步使余额低于0"""
    operations = [5, -10, 3, 2]
    result = below_zero(operations)
    assert result is True


def test_below_zero_boundary_empty_list():
    """边界条件：空列表"""
    operations: List[int] = []
    result = below_zero(operations)
    assert result is False


def test_below_zero_boundary_single_positive():
    """边界条件：单个正数操作"""
    operations = [100]
    result = below_zero(operations)
    assert result is False


def test_below_zero_boundary_single_zero():
    """边界条件：单个零操作"""
    operations = [0]
    result = below_zero(operations)
    assert result is False


def test_below_zero_boundary_single_negative():
    """边界条件：单个负数操作"""
    operations = [-1]
    result = below_zero(operations)
    assert result is True


def test_below_zero_boundary_exact_zero_balance():
    """边界条件：余额恰好为0，不应返回True"""
    operations = [10, -10, 5, -5]
    result = below_zero(operations)
    assert result is False


def test_below_zero_boundary_large_numbers():
    """边界条件：大数字操作"""
    operations = [1000000, -999999, -2]
    result = below_zero(operations)
    assert result is True


def test_below_zero_edge_immediate_recovery():
    """边界条件：余额低于0后立即恢复为正数"""
    operations = [5, -10, 10]
    result = below_zero(operations)
    assert result is True


def test_below_zero_edge_multiple_below_zero_points():
    """边界条件：多个低于0的点，函数应在第一次检测到时就返回"""
    operations = [5, -10, 3, -20, 10]
    result = below_zero(operations)
    assert result is True


def test_below_zero_example_from_docstring_1():
    """文档字符串中的示例1"""
    operations = [1, 2, 3]
    result = below_zero(operations)
    assert result is False


def test_below_zero_example_from_docstring_2():
    """文档字符串中的示例2"""
    operations = [1, 2, -4, 5]
    result = below_zero(operations)
    assert result is True


def test_below_zero_mixed_operations_no_below_zero():
    """混合操作：正负交替但余额从未低于0"""
    operations = [10, -5, 3, -2, 1]
    result = below_zero(operations)
    assert result is False


def test_below_zero_mixed_operations_with_below_zero():
    """混合操作：正负交替且余额低于0"""
    operations = [10, -15, 5, -3]
    result = below_zero(operations)
    assert result is True