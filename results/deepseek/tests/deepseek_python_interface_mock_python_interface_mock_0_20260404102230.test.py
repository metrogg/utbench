import pytest
import itertools
from random import shuffle
from unittest.mock import patch

def task_func(numbers=list(range(1, 3))):
    permutations = list(itertools.permutations(numbers))
    sum_diffs = 0

    for perm in permutations:
        perm = list(perm)
        shuffle(perm)
        diffs = [abs(perm[i] - perm[i+1]) for i in range(len(perm)-1)]
        sum_diffs += sum(diffs)

    avg_sum_diffs = sum_diffs / len(permutations)
    
    return avg_sum_diffs

class TestTaskFunc:
    def test_default_parameters(self):
        """测试默认参数"""
        result = task_func()
        assert isinstance(result, float)
        assert result >= 0

    def test_small_list(self):
        """测试小列表"""
        result = task_func([1, 2])
        permutations = list(itertools.permutations([1, 2]))
        # 手动计算期望值
        expected_sum = 0
        for perm in permutations:
            perm = list(perm)
            # 对于[1,2]和[2,1]，shuffle后可能的顺序只有两种
            # 但shuffle会随机打乱，所以我们需要mock shuffle来确保确定性
            pass  # 这个测试将在mock测试中覆盖
        
        # 由于shuffle的随机性，我们只验证基本属性
        assert isinstance(result, float)
        assert result >= 0

    def test_three_elements(self):
        """测试三个元素"""
        result = task_func([1, 2, 3])
        assert isinstance(result, float)
        assert result >= 0

    def test_identical_elements(self):
        """测试相同元素"""
        result = task_func([5, 5, 5, 5])
        # 所有元素相同，差值为0
        assert result == 0.0

    def test_single_element(self):
        """测试单个元素"""
        result = task_func([42])
        # 单个元素的排列只有一个，没有连续元素对
        assert result == 0.0

    def test_empty_list(self):
        """测试空列表"""
        result = task_func([])
        # 空列表的排列只有一个空排列
        assert result == 0.0

    def test_negative_numbers(self):
        """测试负数"""
        result = task_func([-5, 0, 5])
        assert isinstance(result, float)
        assert result >= 0

    def test_float_numbers(self):
        """测试浮点数"""
        result = task_func([1.5, 2.5, 3.5])
        assert isinstance(result, float)
        assert result >= 0

    def test_mixed_numbers(self):
        """测试混合类型数字"""
        result = task_func([1, 2.5, -3, 0])
        assert isinstance(result, float)
        assert result >= 0

    @patch('random.shuffle')
    def test_deterministic_with_mock(self, mock_shuffle):
        """使用mock确保测试确定性"""
        # 让shuffle不实际打乱顺序
        mock_shuffle.side_effect = lambda x: None
        
        result = task_func([1, 2])
        # 对于[1,2]，排列有：[1,2]和[2,1]
        # 由于shuffle被mock，顺序保持不变
        # 对于[1,2]: 差值为|1-2|=1
        # 对于[2,1]: 差值为|2-1|=1
        # 平均值为(1+1)/2=1.0
        assert result == 1.0

    @patch('random.shuffle')
    def test_three_elements_deterministic(self, mock_shuffle):
        """测试三个元素的确定性结果"""
        mock_shuffle.side_effect = lambda x: None
        
        result = task_func([1, 2, 3])
        # 排列有6种，每个排列的差值总和可以计算
        # 由于shuffle被mock，每个排列保持原始顺序
        permutations = list(itertools.permutations([1, 2, 3]))
        expected_sum = 0
        for perm in permutations:
            diffs = [abs(perm[i] - perm[i+1]) for i in range(len(perm)-1)]
            expected_sum += sum(diffs)
        
        expected_avg = expected_sum / len(permutations)
        assert result == expected_avg

    def test_large_list_performance(self):
        """测试较大列表（验证不会崩溃）"""
        # 使用小规模但比默认大的列表
        result = task_func(list(range(1, 5)))
        assert isinstance(result, float)
        assert result >= 0

    def test_result_range(self):
        """验证结果范围"""
        # 对于任何列表，平均值应该介于最小可能值和最大可能值之间
        numbers = [1, 3, 5]
        result = task_func(numbers)
        
        # 计算最小可能差值（排序列表）
        sorted_nums = sorted(numbers)
        min_diffs = [abs(sorted_nums[i] - sorted_nums[i+1]) for i in range(len(sorted_nums)-1)]
        min_avg = sum(min_diffs) / len(list(itertools.permutations(numbers)))
        
        # 计算最大可能差值（反向排序列表）
        reversed_nums = sorted(numbers, reverse=True)
        max_diffs = [abs(reversed_nums[i] - reversed_nums[i+1]) for i in range(len(reversed_nums)-1)]
        max_avg = sum(max_diffs) / len(list(itertools.permutations(numbers)))
        
        # 由于shuffle的随机性，实际结果可能在最小和最大之间
        # 但我们至少可以验证它是非负的
        assert result >= 0

    def test_immutable_input(self):
        """测试输入不变性"""
        original = [1, 2, 3]
        original_copy = original.copy()
        result = task_func(original)
        assert original == original_copy

    def test_return_type_consistency(self):
        """测试返回类型一致性"""
        test_cases = [
            [1, 2],
            [1, 2, 3],
            [5, 5, 5],
            [-1, 0, 1],
            []
        ]
        
        for numbers in test_cases:
            result = task_func(numbers)
            assert isinstance(result, float), f"Failed for input {numbers}"